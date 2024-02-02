package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/9ssi7/nanoid"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/valyala/fasthttp"
)

type ExchangeRates struct {
	XLira float64 `json:"xLira"`
	YLira float64 `json:"yLira"`
	ZLira float64 `json:"zLira"`
}

type Event struct {
	Old *ExchangeRates `json:"old"`
	New *ExchangeRates `json:"new"`
}

var clients = make(map[*bufio.Writer]bool)
var rates = &ExchangeRates{
	XLira: 1,
	YLira: 1.05,
	ZLira: 1.23,
}

func randomExchangeRate() float64 {
	rate := 1 + (rand.Float64() * 2)
	return math.Round(rate*math.Pow(10, 2)) / math.Pow(10, 2)
}

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	go func() {
		for {
			oldRates := *rates
			rates.XLira = randomExchangeRate()
			rates.YLira = randomExchangeRate()
			rates.ZLira = randomExchangeRate()

			for w := range clients {
				sendEvent(w, &Event{
					Old: &oldRates,
					New: rates,
				})
			}
			time.Sleep(5 * time.Second)
		}
	}()

	app.Get("/events", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		done := c.Context().Done()

		c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			clients[w] = true
			sendEvent(w, map[string]interface{}{"Message": "connected!"})

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
