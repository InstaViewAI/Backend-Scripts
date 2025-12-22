package main

import (
	"context"
	"fix-shadow-script/db"
	i "fix-shadow-script/iot"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type Device struct {
	ID      string `bson:"_id"`
	SpaceID string `bson:"space_id"`
}

func main() {
	atlas_connection_string := "mongodb://admin:Qq8wG62nm.W2@rs0.db.dev.instaview.ai:27017,rs1.db.dev.instaview.ai:27017/?tls=true&tlsInsecure=true"

	db.NewPartnerDBs([]db.Database{
		{Partner: "instaview", ConnectionString: atlas_connection_string},
		{Partner: "luna", ConnectionString: atlas_connection_string}})

	ctx := context.TODO()

	cur, err := db.DeviceCollection("instaview").Find(ctx, bson.M{
		"cluster_group_version": bson.M{"$ne": "v0"},
	})
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

	for _, thing := range thing_name {
		iotCore.UpdateBitRateToEnum(thing)
	}

	fmt.Print("Updated all devices to enum bitrate\n")

}
