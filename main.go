package main

import (
	"AGB-BlocklistSrv/Processing"
	"AGB-BlocklistSrv/Pushers"
	"AGB-BlocklistSrv/Receivers"
	"AGB-BlocklistSrv/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"time"
)

func main() {
	app := fiber.New(fiber.Config{
		Network: fiber.NetworkTCP,
	})

	Processing.Index.Index = Processing.GenerateObjectIndex(config.Configuration.Blocklists)

	app.Use(recover.New())
	log.Infof("Loaded %d blocks, passing to Fiber", len(Processing.Index.Index))

	v1Group := app.Group("/v1")
	v1Group.Post("/BlocklistCallback", submitBlocklistHit)

	Processing.ChosenReceiver = ChooseReceiverFromConfig()
	Processing.ChosenPusher = ChoosePusherFromConfig()

	if Processing.ChosenPusher.CanPusherOperate() {
		v1Group.Post("pusher", Processing.ChosenPusher.HandlePushRequest)
	} else {
		go func() {
			for range time.Tick(time.Hour * 1) { // TODO: Make this configurable
				Processing.Index.Index = Processing.GenerateObjectIndex(config.Configuration.Blocklists)
			}
		}()
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
	Processing.Index.HandleBlocklistCallback(Callback)

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
