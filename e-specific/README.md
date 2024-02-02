# SSE Example: Sending Specific Messages to Users

This example demonstrates how to send specific messages to users based on their email address. Users can register and receive personalized messages from the server.

# How to Run?

- Install Go on your computer.
- Run the main.go file with Go: `go run main.go`
- Open the index.html file in a web browser.
- Register with your username and email address.
- Click the "Try" button.

# What it Does?

Users can register with the server using their email address.
The server sends an SSE event containing a random message every 5 seconds.
Events are only sent to registered users and are personalized based on each user's email address.

Note: In this example, authentication is implemented without tokens for simplicity. In a real-world application, it is important to use a more secure authentication mechanism.