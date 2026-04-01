package main

import (
	"encoding/json"
	"fmt"
	"os"

	i "github.com/instaview/nvr-shadow-update/iot"
)

type Device struct {
	ID      string `bson:"_id"`
	SpaceID string `bson:"space_id"`
}

type ThingEntry struct {
	ThingName string `json:"thing_name"`
	DeviceID  string `json:"device_id"`
	SpaceID   string `json:"space_id"`
}

func main() {
	// atlas_connection_string := "mongodb://admin:password@rs0.db.staging.instaview.ai:27017,rs1.db.staging.instaview.ai:27017,rs2.db.staging.instaview.ai:27017/?replicaSet=stg-mongo&tls=true&tlsInsecure=true"

	// db.NewPartnerDBs([]db.Database{
	// 	{Partner: "alder_nvr", ConnectionString: atlas_connection_string}})

	// ctx := context.TODO()

	// filter := bson.M{
	// 	"thing_name": bson.M{
	// 		"$exists": true,
	// 	},
	// }

	// cur, err := db.DeviceCollection("alder_nvr").Find(ctx,filter)
	// if err != nil {
	// 	println("Error finding device:", err)
	// }

	// var thing_name []string

	// for cur.Next(ctx) {
	// 	var device Device
	// 	if err := cur.Decode(&device); err != nil {
	// 		println("Error decoding device:", err.Error())
	// 		continue
	// 	}
	// 	thing_name = append(thing_name, fmt.Sprintf("%s_%s", device.SpaceID, device.ID))
	// }

	// iotCore := i.NewIoTCore()

	// for _, thing := range thing_name {
	// 	if err := iotCore.UpdateNVRShadow(thing); err != nil {
  	//  		fmt.Printf("Failed to update shadow: %v\n", err)
	// 	}
	// 	//println(thing)
	// }
	// fmt.Print("Updated the All NVR Shadows\n")

	// Read file
	data, err := os.ReadFile("alder_nvr_things_with_subs.json")
	if err != nil {
		panic(fmt.Errorf("failed to read json file: %w", err))
	}

	// Unmarshal JSON
	var things []ThingEntry
	if err := json.Unmarshal(data, &things); err != nil {
		panic(fmt.Errorf("failed to unmarshal json: %w", err))
	}

	iotCore := i.NewIoTCore()

	// Iterate through all the things read from the JSON File and update shadow
	for _, t := range things {
		// first check the things currently exist or not 
		if !iotCore.ThingExists(t.ThingName) {
			fmt.Printf("Thing does NOT exist, skipping: %s\n", t.ThingName)
			continue
		}

		err := iotCore.UpdateNVRSubShadow(t.ThingName)
		if err != nil {
			fmt.Printf("Failed to Update the Subscription Shadow for %s: %v\n", t.ThingName, err)
			continue
		}
		fmt.Printf("Updated the Subscription Shadow %s successfully\n", t.ThingName)
	}
	fmt.Print("Updated the All NVR Shadows\n")

}
