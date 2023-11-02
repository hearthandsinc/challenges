# iOS - Simple Server

This server is used as a companion for the iOS challenge. Persistence is
in-memory, restarting the server will reset its state.

The server will sometime fail to send a response, take a very long time to
respond, or send duplicate events in the stream. You should react accordingly.

## How to run the server?

1. You need to install Go on your system: [official documentation](https://go.dev/doc/install)
2. You can then run the server with: `go run .`

## API Documentation

### Notes

- All the `POST` methods expect an `Idempotency-Key` header, it must be a
  string unique for each request (and their respective retries). [Read
  more](https://stripe.com/docs/api/idempotent_requests). For example, a valid
  header could be: `Idempotency-Key: 459cfe7e-5952-43a0-a0ff-b2d8f1f4cfad`.

### Methods

**Chats:**
- `GET /chats`: list all chats

**Messages:**
- `POST /chats/{chat_id}/messages`: send a new message in a chat
  * the message should be passed as a JSON payload
- `GET /chats/{chat_id}/messages`: list all messages in a chat

### Entities

- `Chat`
  * `id` (`number`): the chat id
  * `name` (`string`): the chat name

- `Message`
  * `id` (`number`): the message id
  * `author` (`string`): either the constant string `"user"` or `"bot"`
  * `text` (`string`): the actual text of the message
