package main

import (
	"AGB-BlocklistSrv/Processing"
	"AGB-BlocklistSrv/Pushers"
	"AGB-BlocklistSrv/Receivers"
	"AGB-BlocklistSrv/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	index = Processing.WorldObjectIndex{Index: Processing.GenerateObjectIndex(config.Configuration.Blocklists)}
)

func main() {
	app := fiber.New(fiber.Config{
		Network: fiber.NetworkTCP,
	})

	app.Use(recover.New())
	log.Infof("Loaded %d blocks, passing to Fiber", len(index.Index))

	v1Group := app.Group("/v1")
	v1Group.Post("/BlocklistCallback", submitBlocklistHit)

	Processing.ChosenReceiver = ChooseReceiverFromConfig()
	Processing.ChosenPusher = ChoosePusherFromConfig()

	if Processing.ChosenPusher.CanPusherOperate() {
		v1Group.Post("pusher", Processing.ChosenPusher.HandlePushRequest)
	}

	err := app.Listen(":80")
	if err != nil {
		panic(err)
	}
}
func submitBlocklistHit(c *fiber.Ctx) error {
	c.Accepts("application/json")
	var Callback Processing.CallbackContainer

	if err := c.BodyParser(&Callback); err != nil {
		panic(err)
	}
	index.HandleBlocklistCallback(Callback)

	return c.SendStatus(fiber.StatusNoContent)
}
func ChooseReceiverFromConfig() Processing.Receiver {
	switch config.Configuration.Reciever {
	case "influxdb":
		return Receivers.Influxdb{}
	case "stub":
		return Receivers.Stub{}
	default:
		panic("Invalid receiver")
	}
}

func ChoosePusherFromConfig() Processing.Pusher {
	switch config.Configuration.Pusher {
	case "grafghanno":
		return Pushers.GrafanaGithubWebhookAnnotation{}
	default:
		panic("Invalid pusher")
	}
}
