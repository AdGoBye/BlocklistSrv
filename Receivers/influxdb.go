package Receivers

import (
	"AGB-BlocklistSrv/Processing"
	"context"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"os"
	"strconv"
	"time"
)

var (
	client = GetInstance()
)

type Influxdb struct{}

func (influx Influxdb) SendToRemote(misses []Processing.Gameobject, world *Processing.WorldObject) {
	now := time.Now()
	callbackSetId, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	for i, miss := range misses {
		p := influxdb2.NewPointWithMeasurement("callbacks").
			AddTag("callbackSetId", callbackSetId.String()).
			AddTag("blocklists", *miss.ParentBlocklist).
			AddTag("uniq", strconv.Itoa(i)).
			AddField("objectName", miss.Name).
			AddField("world", world.FriendlyName).
			SetTime(now)
		if miss.Position != nil {
			p.AddField("position", miss.Position)
		}
		if miss.Parent != nil {
			p.AddField("parentName", miss.Parent.Name)
			if miss.Parent.Position != nil {
				p.AddField("parentPosition", miss.Parent.Position)
			}
		}

		err = client.WritePoint(context.Background(), p)
		if err != nil {
			log.Error(err)
		}
	}
}

func GetInstance() api.WriteAPIBlocking {
	influxClient := influxdb2.NewClient(os.Getenv("INFLUXDB_LOCATION"), os.Getenv("DOCKER_INFLUXDB_INIT_ADMIN_TOKEN"))
	defer influxClient.Close()
	return influxClient.WriteAPIBlocking(os.Getenv("DOCKER_INFLUXDB_INIT_ORG"), os.Getenv("DOCKER_INFLUXDB_INIT_BUCKET"))
}
