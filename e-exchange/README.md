# SSE Example: Real-time Exchange Rate Updates

This example demonstrates how to use SSE to provide real-time updates of exchange rates. The server randomly generates new exchange rates and sends them to connected clients.

# How to Run?

- Install Go on your computer.
- Run the main.go file with Go: `go run main.go`
- Open the index.html file in a web browser.

# What it Does?

The server sends an SSE event containing randomly updated exchange rates every 5 seconds.
JavaScript uses the EventSource API to listen for events sent by the server.
Each time a new event is received, the exchange rates are updated in the browser.