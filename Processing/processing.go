package Processing

import (
	"github.com/gofiber/fiber/v2"
)

var (
	ChosenReceiver Receiver
	ChosenPusher   Pusher
)

// A Receiver is the actual endpoint that gets the blocklist data.
type Receiver interface {
	SendToRemote(misses []Gameobject, world *WorldObject)
}

// A Pusher is something that pushes data to the BlocklistSrv.
type Pusher interface {
	HandlePushRequest(c *fiber.Ctx) error
	CanPusherOperate() bool
}
