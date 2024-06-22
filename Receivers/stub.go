package Receivers

import (
	"AGB-BlocklistSrv/Processing"
	"fmt"
)

type Stub struct{}

func (stub Stub) SendToRemote(misses []Processing.Gameobject, world *Processing.WorldObject) {
	fmt.Println("Received hit for " + world.FriendlyName + ":")
	for _, miss := range misses {
		fmt.Printf("%v", miss)
	}
}
