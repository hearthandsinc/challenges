package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/r3labs/sse/v2"
	"golang.org/x/exp/maps"
)

var (
	hostname = envOrDefault("HOSTNAME", "localhost")
	port     = envOrDefault("PORT", "3000")
)

func main() {
	events := sse.New()
	events.CreateStream("message")

	app := NewApp(events)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Handle("/events", http.HandlerFunc(events.ServeHTTP))

	r.Group(func(r chi.Router) {
		// chaos middleware
		r.Use(func(next http.Handler) http.Handler {
			fn := func(w http.ResponseWriter, r *http.Request) {
				// 5% chance of unavailable
				if rand.Intn(100) < 5 {
					http.Error(w, "unavailable", http.StatusServiceUnavailable)
					return
				}

				// 5% chance of timeout
				if rand.Intn(100) < 5 {
					<-time.After(10 * time.Second)
					http.Error(w, "timeout", http.StatusGatewayTimeout)
					return
				}

				next.ServeHTTP(w, r)
			}
			return http.HandlerFunc(fn)
		})

		r.Get("/chats", app.GetChats)
		r.Post("/chats/{chatID}/messages", app.PostMessages)
		r.Get("/chats/{chatID}/messages", app.GetMessages)
	})

	addr := fmt.Sprintf("%s:%s", hostname, port)
	fmt.Printf("Starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		panic(fmt.Errorf("failed to start server on %s: %w", addr, err))
	}
}

//
// Entities
//

type Chat struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`

	messages []*Message // not sent to the clients
}

type Message struct {
	ID     uint32    `json:"id"`
	ChatID uint32    `json:"chat_id"`
	Author string    `json:"author"`
	Text   string    `json:"text"`
	SentAt time.Time `json:"sent_at"`
}

// App is both our controller and our data store. This coupling allows to keep
// the implementation simple.
type App struct {
	events *sse.Server // server-sent events server

	mu              sync.RWMutex
	idempotencyKeys map[string]struct{} // store idempotency keys
	store           map[uint32]*Chat    // in-memory store of chats and messages
}

func NewApp(events *sse.Server) *App {
	now := time.Now()

	// store contains all the chats, manually indexed by their ID
	store := map[uint32]*Chat{}
	for _, chat := range []*Chat{
		{Name: "John", messages: []*Message{
			{ID: newID(), Author: "bot", Text: "Sounds good 👍", SentAt: now},
		}},
		{Name: "Jessica", messages: []*Message{
			{ID: newID(), Author: "bot", Text: "How are you!?", SentAt: now.Add(-1 * time.Minute)},
		}},
		{Name: "Matt", messages: []*Message{
			{ID: newID(), Author: "bot", Text: "ok chat soon :)", SentAt: now.Add(-32 * time.Minute)},
		}},
		{Name: "Sarah", messages: []*Message{
			{ID: newID(), Author: "bot", Text: "ok talk later!", SentAt: now.Add(-32 * time.Hour)},
		}},
	} {
		chatID := newID()
		chat.ID = chatID
		for _, message := range chat.messages {
			message.ChatID = chatID
		}
		store[chatID] = chat
	}

	return &App{
		events:          events,
		idempotencyKeys: map[string]struct{}{},
		store:           store,
	}
}

func (a *App) GetChats(w http.ResponseWriter, r *http.Request) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	chats := maps.Values(a.store)

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, chats)
}

func (a *App) PostMessages(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()

	idempotencyKey := r.Header.Get("Idempotency-Key")
	if len(idempotencyKey) == 0 {
		http.Error(w, "missing idempotency key", http.StatusBadRequest)
		return
	}

	if _, ok := a.idempotencyKeys[idempotencyKey]; ok {
		w.WriteHeader(http.StatusOK)
		return
	}

	var message Message
	if err := render.DecodeJSON(r.Body, &message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chat, err := a.findChatByID(chi.URLParam(r, "chatID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newMessage := &Message{
		ID:     newID(),
		ChatID: chat.ID,
		Author: "user",
		Text:   message.Text,
		SentAt: time.Now(),
	}

	chat.messages = append(chat.messages, newMessage)
	a.idempotencyKeys[idempotencyKey] = struct{}{}

	if err := a.publishEvent("message", newMessage); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *App) GetMessages(w http.ResponseWriter, r *http.Request) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	chat, err := a.findChatByID(chi.URLParam(r, "chatID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, chat.messages)
}

func (a *App) findChatByID(rawID string) (*Chat, error) {
	chatID, err := strconv.ParseUint(rawID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid chat ID %s: %w", rawID, err)
	}

	chat, ok := a.store[uint32(chatID)]
	if !ok {
		return nil, fmt.Errorf("chat %s not found", rawID)
	}

	return chat, nil
}

func (a *App) publishEvent(stream string, payload any) error {
	serialized, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to serialize payload: %w", err)
	}

	for i := 0; i < rand.Intn(3); i++ {
		a.events.Publish(stream, &sse.Event{Data: serialized})
	}

	return nil
}

//
// Helpers
//

// lastID is the last ID generated by newID.
var lastID = uint32(time.Now().Unix())

// newID generates a new ID, unique for the lifetime of this server.
func newID() uint32 {
	return atomic.AddUint32(&lastID, 1)
}

// envOrDefault returns the value of the environment variable at the given key.
// Fallbacks to the given default if the value found is missing or empty.
func envOrDefault(key string, defaultValue string) string {
	if v := strings.TrimSpace(os.Getenv(key)); len(v) > 0 {
		return v
	}
	return defaultValue
}
