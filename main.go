package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	configuration = loadConfiguration()
	index         = WorldObjectIndex{
		Index: generateObjectIndex(configuration.Blocklists),
	}
	influxdb = getInfluxDBClient()
)

func main() {
	app := fiber.New(fiber.Config{
		Network: fiber.NetworkTCP,
	})

	app.Use(recover.New())
	log.Infof("Loaded %d blocks, passing to Fiber", len(index.Index))

	v1Group := app.Group("/v1")
	v1Group.Post("/BlocklistCallback", submitBlocklistHit)

	app.Post("webhook", handleWebhookRequest)

	err := app.Listen(":80")
	if err != nil {
		panic(err)
	}
}
func submitBlocklistHit(c *fiber.Ctx) error {
	c.Accepts("application/json")
	var Callback CallbackContainer

	if err := c.BodyParser(&Callback); err != nil {
		panic(err)
	}
	index.HandleBlocklistCallback(Callback)

	return c.SendStatus(fiber.StatusNoContent)
}
