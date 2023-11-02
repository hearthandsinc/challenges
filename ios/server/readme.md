# iOS Challenge - Server Companion

This server is used as a companion for the iOS challenge. Persistence is
in-memory, restarting the server will reset its state.

The server will sometime be a bit capricious, you might observe:
- failure to receive a response
- server timing out when trying to send a message
- duplicate events being received in the messages stream

When that happens, you should retry and react accordingly.

## How to start the server?

1. You need to install Go on your system: follow the [official documentation](https://go.dev/doc/install)
2. Then execute the command: `go run .`

## API Documentation

This API exposes a few HTTP endpoints, and an SSE stream. It allows to list
chats, and read & write messages.

### Notes

- All the `POST` endpoints expect an `Idempotency-Key` header, it must be a string
  unique to each request (and their respective retries). For more details, see
  [here](https://stripe.com/docs/api/idempotent_requests). A valid header could
  be: `Idempotency-Key: 459cfe7e-5952-43a0-a0ff-b2d8f1f4cfad`.

### Endpoints

- `GET /events?stream=messages`: an [SSE](https://en.wikipedia.org/wiki/Server-sent_events) stream that sends you `Message` entities as they are received by the server
- `GET /chats`: returns a list of all `Chat` entities
- `GET /chats/{chat_id}/messages`: returns a list of the 100 most recents `Message` entities in a chat
- `POST /chats/{chat_id}/messages`: send a new message in a chat and returns
  the newly created `Message` entity. It expects a JSON payload of the form: `{ "text":
  "..." }`

### Entities

#### `Chat`
- `id` (`number`): the chat id
- `name` (`string`): the chat name

#### `Message`
- `id` (`number`): the message id
- `chat_id` (`number`): the id of the chat this message belongs to
- `author` (`string`): the message author (either the constant string `"user"` or `"bot"`)
- `text` (`string`): the actual content of the message
- `sent_at` (`string`): the date at which the message was sent
