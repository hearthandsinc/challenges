# iOS - Simple Server

This is a simple server that is used as a companion for the iOS challenge.
Persistence is in-memory, restarting the server will reset its state.

You need to install Go ([documentation](https://go.dev/doc/install)), you can
then run it with: `go run .`.

The server will sometime fail to send a response, take a very long time to
respond, or send duplicate events in the stream. You should react accordingly.

## Documentation

### Methods

- `POST /stream`
- `POST /messages`: post a message
- `GET /messages`: get all messages
- `GET /messages/{id}`: get a messages
- `PUT /messages/{id}`: update a message
- `DELETE /messages/{id}`: delete a message

### Entities

- `Message`
  * `idempotency_key`
  * `id`
  * `text`
  * `created_at`
