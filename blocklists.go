package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2/log"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/pelletier/go-toml/v2"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type WorldObject struct {
	FriendlyName      string
	GameObjectMapping map[string]Gameobject
}
type WorldObjectIndex struct {
	Index map[string]WorldObject
}

type Blocklist struct {
	Title       string  `toml:"-"`
	Description string  `toml:"-"`
	Maintainer  string  `toml:"-"`
	Blocks      []Block `toml:"block"`
}

type Block struct {
	FriendlyName string       `toml:"friendly_name" json:"friendly_name"`
	WorldId      string       `toml:"world_id" json:"world_id"`
	GameObjects  []Gameobject `toml:"game_objects" json:"game_objects"`
}

type CallbackContainer struct {
	Version          int      `json:"Version"`
	WorldId          string   `json:"WorldId"`
	UnmatchedObjects []string `json:"UnmatchedObjects"`
}

type Gameobject struct {
	Name     string              `toml:"name" json:"Name"`
	Position *GameobjectPosition `toml:"position" json:"Position"`
	Parent   *Gameobject         `toml:"parent" json:"Parent"`
}
type GameobjectPosition struct {
	X float64 `toml:"x" json:"X"`
	Y float64 `toml:"y" json:"Y"`
	Z float64 `toml:"z" json:"Z"`
}

func (index WorldObjectIndex) GetWorldById(HashedWorldId string) *WorldObject {
	if val, exists := index.Index[HashedWorldId]; exists {
		return &val
	}
	return nil
}

func (index WorldObjectIndex) HandleBlocklistCallback(object CallbackContainer) {
	var world *WorldObject
	if world = index.GetWorldById(object.WorldId); world == nil { // Return immediately if not under our supervision
		return
	}

	var misses []Gameobject
	for _, b64 := range object.UnmatchedObjects {
		if val, exists := world.GameObjectMapping[b64]; exists {
			misses = append(misses, val)
		}
	}

	if len(misses) == 0 { // None of the objects reported back are relevant to us
		return
	}
	p := influxdb2.NewPointWithMeasurement("callbacks").
		AddTag("worldFriendlyName", world.FriendlyName).
		AddField("misses", misses).
		SetTime(time.Now())
	err := influxdb.WritePoint(context.Background(), p)
	if err != nil {
		log.Error(err)
	}
}

func generateObjectIndex(blocklistsLocations []string) (mapping map[string]WorldObject) {
	mapping = make(map[string]WorldObject)
	for _, blocklistUrl := range blocklistsLocations {
		blocklistObject, err := fetchBlocklist(blocklistUrl)
		if err != nil {
			panic(err)
		}

		for _, object := range blocklistObject.Blocks {
			widHashEncoded := base64.StdEncoding.EncodeToString(stringToHash(object.WorldId))

			ensureMappingInititalization(mapping, widHashEncoded, object)

			mapHashesToRealObjects(object, mapping, widHashEncoded)
		}
	}
	return mapping
}

func mapHashesToRealObjects(object Block, mapping map[string]WorldObject, widHashEncoded string) {
	for _, gameObject := range object.GameObjects {
		marshal, err := json.Marshal(gameObject)
		if err != nil {
			panic(err)
		}

		b64 := base64.StdEncoding.EncodeToString(stringToHash(marshal))
		mapping[widHashEncoded].GameObjectMapping[b64] = gameObject
	}
}

func ensureMappingInititalization(mapping map[string]WorldObject, widhashEncoded string, block Block) {
	if _, exists := mapping[widhashEncoded]; !exists {
		mapping[widhashEncoded] = WorldObject{
			FriendlyName:      block.FriendlyName,
			GameObjectMapping: make(map[string]Gameobject),
		}
	}
}

func fetchBlocklist(location string) (Blocklist, error) {
	uri, err := url.ParseRequestURI(location)
	if err != nil {
		return Blocklist{}, err
	}
	var blocklistBytes []byte

	switch uri.Scheme {
	case "http", "https":
		blocklistBytes, err = downloadBlocklistFromHTTP(location)
		if err != nil {
			return Blocklist{}, err
		}
	case "file":
		_, err = os.Stat(uri.Path)
		if err != nil {
			return Blocklist{}, err
		}
		blocklistBytes, err = os.ReadFile(uri.Path)
		if err != nil {
			return Blocklist{}, err
		}
	default:
		return Blocklist{}, errors.New("unsupported scheme: " + uri.Scheme)
	}

	var blocklistObject Blocklist
	err = toml.Unmarshal(blocklistBytes, &blocklistObject)
	if err != nil {
		return Blocklist{}, err
	}
	return blocklistObject, nil
}

func downloadBlocklistFromHTTP(location string) ([]byte, error) {
	resp, err := http.Get(location)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func stringToHash[inputs string | []byte](input inputs) (output []byte) {
	hash := sha256.New()
	hash.Write([]byte(input)) // []byte as []byte stays []byte, so we don't have to explicitly check
	output = hash.Sum(nil)
	return output
}

func getInfluxDBClient() api.WriteAPIBlocking {
	influxClient := influxdb2.NewClient(os.Getenv("INFLUXDB_LOCATION"), os.Getenv("DOCKER_INFLUXDB_INIT_ADMIN_TOKEN"))
	defer influxClient.Close()
	return influxClient.WriteAPIBlocking(os.Getenv("DOCKER_INFLUXDB_INIT_ORG"), os.Getenv("DOCKER_INFLUXDB_INIT_BUCKET"))
}
