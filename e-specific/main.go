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

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Event struct {
	Content string `json:"content"`
}

var clients = make(map[*bufio.Writer]User)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Post("/register", func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return err
		}
		for _, u := range clients {
			if u.Email == user.Email {
				c.Status(fiber.StatusConflict)
				return c.SendString("User already exists")
			}
		}
		c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			clients[w] = *user
		}))
		c.Status(fiber.StatusOK)
		return nil
	})

	app.Post("/api", authMiddleware, func(c *fiber.Ctx) error {
		user := c.Locals("user").(User)

		for w, u := range clients {
			if u.Email == user.Email {
				sendEvent(w, &Event{Content: "Hello from API!"})
			}
		}

		c.Status(fiber.StatusOK)
		return nil
	})

	app.Get("/events", authMiddleware, func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")
		user := c.Locals("user").(User)

		done := c.Context().Done()

		c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			sendEvent(w, &Event{Content: fmt.Sprintf("Hello %s!", user.Name)})

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

func sendEvent(w *bufio.Writer, v interface{}) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return
	}
	fmt.Fprintf(w, "data: %s\n\n", bytes)
	id, _ := nanoid.New() // last event id
	fmt.Fprintf(w, "id: %s\n\n", id)
	w.Flush()
}

func authMiddleware(c *fiber.Ctx) error {
	bearerEmail := c.Query("email")
	if bearerEmail == "" {
		c.Status(fiber.StatusUnauthorized)
		return c.SendString("Unauthorized")
	}

	for _, user := range clients {
		if user.Email == bearerEmail {
			c.Locals("user", user)
			return c.Next()
		}
	}

	c.Status(fiber.StatusUnauthorized)
	return c.SendString("Unauthorized")
}
