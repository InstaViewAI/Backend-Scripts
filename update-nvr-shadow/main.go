package main

import (
	"context"
	"fmt"

	i "github.com/instaview/nvr-shadow-update/iot"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/instaview/nvr-shadow-update/db"
)

type Device struct {
	ID      string `bson:"_id"`
	SpaceID string `bson:"space_id"`
}

func main() {
	atlas_connection_string := "mongodb://admin:password@rs0.db.staging.instaview.ai:27017,rs1.db.staging.instaview.ai:27017,rs2.db.staging.instaview.ai:27017/?replicaSet=stg-mongo&tls=true&tlsInsecure=true"

	db.NewPartnerDBs([]db.Database{
		{Partner: "alder_nvr", ConnectionString: atlas_connection_string}})

	ctx := context.TODO()

	filter := bson.M{
		"thing_name": bson.M{
			"$exists": true,
		},
	}

	cur, err := db.DeviceCollection("alder_nvr").Find(ctx,filter)
	if err != nil {
		println("Error finding device:", err)
	}

	var thing_name []string

	for cur.Next(ctx) {
		var device Device
		if err := cur.Decode(&device); err != nil {
			println("Error decoding device:", err.Error())
			continue
		}
		thing_name = append(thing_name, fmt.Sprintf("%s_%s", device.SpaceID, device.ID))
	}

	iotCore := i.NewIoTCore()

	// for _, thing := range thing_name {
	// 	fmt.Println(thing)
	// }
	thing_my_name := "36sMjQdQgZkRqQWTz5F8ASMY4Ta_cbx3fjdgipbp3z29mj"
	fmt.Println(thing_my_name)

	if err := iotCore.UpdateNVRShadow(thing_my_name); err != nil {
  	  fmt.Printf("Failed to update shadow: %v\n", err)
	}
	fmt.Print("all devices things names done\n")

}
