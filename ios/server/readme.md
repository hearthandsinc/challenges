# iOS - Simple Server

This server is used as a companion for the iOS challenge. Persistence is
in-memory, restarting the server will reset its state.

The server will sometime be a bit capricious, you might observe:
- failure to receive a response
- server timing out when trying to send a response
- duplicate events being received in the messages stream

You should retry and react accordingly.

## How to start the server?

1. You need to install Go on your system ([official documentation](https://go.dev/doc/install))
2. You can then run the server with: `go run .`

## API Documentation

### Notes

- All the `POST` methods expect an `Idempotency-Key` header, it must be a string
  unique for each request (and their respective retries). For more details, see
  [here](https://stripe.com/docs/api/idempotent_requests). A valid header could
  be: `Idempotency-Key: 459cfe7e-5952-43a0-a0ff-b2d8f1f4cfad`.

### Methods

- `GET /events?stream=messages`: a [Server-Sent Events](https://en.wikipedia.org/wiki/Server-sent_events) stream that sends you new `Message`s
- `GET /chats`: list all `Chat`s.
- `POST /chats/{chat_id}/messages`: send a new message in a chat. Expects a `{
  "text": "..." }` payload.
- `GET /chats/{chat_id}/messages`: list all `Message`s in a chat

### Entities

- `Chat`
  * `id` (`number`): the chat id
  * `name` (`string`): the chat name

- `Message`
  * `id` (`number`): the message id
  * `author` (`string`): either the constant string `"user"` or `"bot"`
  * `text` (`string`): the actual text of the message
