package main

import (
	"bufio"
	"encoding/json"
	"fmt"

	"github.com/9ssi7/nanoid"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/valyala/fasthttp"
)

type Event struct {
	Content string `json:"content"`
}

var clients = make(map[*bufio.Writer]bool)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Post("/api", func(c *fiber.Ctx) error {
		for w := range clients {
			sendEvent(w, &Event{Content: "Hello World too!"})
		}

		c.Status(fiber.StatusOK)
		return nil
	})

	app.Get("/events", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		done := c.Context().Done()

		c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			clients[w] = true
			sendEvent(w, &Event{Content: "Hello World!"})

			go func() {
				<-done
				delete(clients, w)
			}()

			select {}
		}))
		return nil
	})

	err := app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}

func sendEvent(w *bufio.Writer, event *Event) {
	bytes, err := json.Marshal(event)
	if err != nil {
		return
	}
	fmt.Fprintf(w, "data: %s\n\n", bytes)
	id, _ := nanoid.New() // last event id
	fmt.Fprintf(w, "id: %s\n\n", id)
	w.Flush()
}
