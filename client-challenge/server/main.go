package main

import (
	"cmp"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/r3labs/sse/v2"
	"golang.org/x/exp/maps"
)

const (
	// maximum number of entities returned in a single request
	limit = 100
)

var (
	// hostname used by the server
	hostname = cmp.Or(os.Getenv("HTTP_HOST"), "localhost")

	// port used by the server
	port = cmp.Or(os.Getenv("PORT"), "3000")
)

//go:embed data.json
var dataJSON []byte

func main() {
	events := sse.New()
	app := NewApp(events)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Handle("/events", http.HandlerFunc(events.ServeHTTP))

	r.Group(func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Use(chaosMiddleware)
		r.Get("/chats", app.GetChats)
		r.Get("/chats/{chatID}/messages", app.GetChatMessages)
		r.Post("/chats/{chatID}/messages", app.PostChatMessages)
	})

	addr := net.JoinHostPort(hostname, port)
	fmt.Printf("Starting server on http://%s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		panic(fmt.Errorf("failed to start server on %s: %w", addr, err))
	}
}

//
// Entities
//

type Chat struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	messages []*Message
}

type Message struct {
	ID             uuid.UUID `json:"id"`
	ChatID         uuid.UUID `json:"chat_id"`
	Author         string    `json:"author"`
	Text           string    `json:"text"`
	SentAt         time.Time `json:"sent_at"`
	IdempotencyKey string    `json:"idempotency_key"`
}

// App is both our controller and our data store. This coupling allows to keep
// the implementation simple.
type App struct {
	events *sse.Server // server-sent events server
	jokes  []string    // jokes to send

	mu              sync.RWMutex
	idempotencyKeys map[string]struct{} // store idempotency keys
	store           map[uuid.UUID]*Chat // in-memory store of chats and messages
}

func NewApp(events *sse.Server) *App {
	var jokes []string
	if err := json.Unmarshal(dataJSON, &jokes); err != nil {
		panic(fmt.Errorf("failed to unmarshal data: %w", err))
	}

	events.CreateStream("messages")

	now := time.Now()

	// store contains all the chats, indexed by their ID
	store := map[uuid.UUID]*Chat{}
	for _, chat := range []*Chat{
		{ID: uuid.New(), Name: "John", messages: []*Message{
			{ID: uuid.New(), Author: "bot", Text: "Sounds good 👍", SentAt: now, IdempotencyKey: uuid.New().String()},
		}},
		{ID: uuid.New(), Name: "Jessica", messages: []*Message{
			{ID: uuid.New(), Author: "bot", Text: "How are you!?", SentAt: now.Add(-1 * time.Minute), IdempotencyKey: uuid.New().String()},
		}},
		{ID: uuid.New(), Name: "Matt", messages: []*Message{
			{ID: uuid.New(), Author: "bot", Text: "ok chat soon :)", SentAt: now.Add(-32 * time.Minute), IdempotencyKey: uuid.New().String()},
		}},
		{ID: uuid.New(), Name: "Sarah", messages: []*Message{
			{ID: uuid.New(), Author: "bot", Text: "ok talk later!", SentAt: now.Add(-24 * time.Hour), IdempotencyKey: uuid.New().String()},
		}},
	} {
		for _, message := range chat.messages {
			message.ChatID = chat.ID
		}
		store[chat.ID] = chat
	}

	return &App{
		events:          events,
		jokes:           jokes,
		idempotencyKeys: map[string]struct{}{},
		store:           store,
	}
}

func (app *App) GetChats(w http.ResponseWriter, r *http.Request) {
	app.mu.RLock()
	defer app.mu.RUnlock()

	chats := maps.Values(app.store)

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, chats)
}

func (app *App) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	app.mu.RLock()
	defer app.mu.RUnlock()

	chat, err := app.findChatByID(chi.URLParam(r, "chatID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	start := max(len(chat.messages)-limit, 0)
	end := len(chat.messages)

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, chat.messages[start:end])
}

func (app *App) PostChatMessages(w http.ResponseWriter, r *http.Request) {
	app.mu.Lock()
	defer app.mu.Unlock()

	idempotencyKey := r.Header.Get("Idempotency-Key")
	if len(idempotencyKey) == 0 {
		http.Error(w, "missing idempotency key", http.StatusBadRequest)
		return
	}

	if _, ok := app.idempotencyKeys[idempotencyKey]; ok {
		w.WriteHeader(http.StatusOK)
		return
	}

	var message Message
	if err := render.DecodeJSON(r.Body, &message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chat, err := app.findChatByID(chi.URLParam(r, "chatID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newMessage := &Message{
		ID:             uuid.New(),
		ChatID:         chat.ID,
		Author:         "user",
		Text:           message.Text,
		SentAt:         time.Now(),
		IdempotencyKey: idempotencyKey,
	}

	chat.messages = append(chat.messages, newMessage)
	app.publishEvent("messages", newMessage)

	go app.sendDelayedAnswer(chat)

	app.idempotencyKeys[idempotencyKey] = struct{}{}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, newMessage)
}

func (app *App) sendDelayedAnswer(chat *Chat) {
	<-time.After(3*time.Second + time.Duration(rand.Intn(3))*time.Second)

	app.mu.Lock()
	defer app.mu.Unlock()

	text := app.jokes[rand.Intn(len(app.jokes))]

	newMessage := &Message{
		ID:             uuid.New(),
		ChatID:         chat.ID,
		Author:         "bot",
		Text:           text,
		SentAt:         time.Now(),
		IdempotencyKey: uuid.New().String(),
	}

	chat.messages = append(chat.messages, newMessage)
	app.publishEvent("messages", newMessage)
}

func (app *App) findChatByID(rawID string) (*Chat, error) {
	parsedID, err := uuid.Parse(rawID)
	if err != nil {
		return nil, fmt.Errorf("invalid chat ID %s: %w", rawID, err)
	}

	chat, ok := app.store[parsedID]
	if !ok {
		return nil, fmt.Errorf("chat not found with ID %s", rawID)
	}

	return chat, nil
}

func (app *App) publishEvent(stream string, payload any) {
	serialized, err := json.Marshal(payload)
	if err != nil {
		panic(fmt.Errorf("failed to serialize payload: %w", err))
	}

	for i := 0; i <= rand.Intn(3); i++ {
		app.events.Publish(stream, &sse.Event{Data: serialized})
	}
}

// chaosMiddleware makes sure to randomly return errors and adds artificial
// latency to simulate a real world system.
func chaosMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// latency
		<-time.After(time.Duration(rand.Intn(1500)) * time.Millisecond)

		// 5% chance of unavailable
		if rand.Intn(100) < 5 {
			http.Error(w, "unavailable", http.StatusServiceUnavailable)
			return
		}

		// 5% chance of timeout after 5 seconds
		if rand.Intn(100) < 5 {
			<-time.After(5 * time.Second)
			http.Error(w, "timeout", http.StatusGatewayTimeout)
			return
		}

		next.ServeHTTP(w, r)

		// latency
		<-time.After(time.Duration(rand.Intn(1500)) * time.Millisecond)
	})
}
