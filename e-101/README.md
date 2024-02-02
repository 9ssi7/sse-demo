# SSE Example: Sending a "Hello World!" Message

This example demonstrates the basic functionality of SSE by sending a "Hello World!" message to all connected clients.

## How to Run?

- Install Go on your computer.
- Run the main.go file with Go: `go run main.go`
- Open the index.html file in a web browser.

## What it Does?

The server sends an SSE event containing a "Hello World!" message every 5 seconds.
JavaScript uses the EventSource API to listen for events sent by the server.
Each time a new event is received, the message is displayed in the browser.