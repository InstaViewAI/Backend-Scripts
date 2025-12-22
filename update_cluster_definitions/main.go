package main

import (
	"context"
	"fmt"
	"log"

	"github.com/update_cluster_group/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	updates := getUpdateMap()

	// fmt.Printf("map %+v\n", updates)

	atlas_connection_string := "mongodb://admin:password@rs0.db.staging.instaview.ai:27017,rs1.db.staging.instaview.ai:27017,rs2.db.staging.instaview.ai:27017/?replicaSet=stg-mongo&tls=true&tlsInsecure=true"

	db.NewPartnerDBs([]db.Database{
		{Partner: "instaview", ConnectionString: atlas_connection_string},
	})

	// update Cluster Definitions
	fmt.Println("Updating cluster definitions...")
	updateClusterDefinitions(updates)
	fmt.Println("Cluster definitions updated successfully.")
}

func updateClusterDefinitions(updates map[string]map[string]any) {

	ctx := context.TODO()

	partnerId := "instaview"

	cur, err := db.ClusterDefinition(partnerId).Find(ctx, bson.M{})
	if err != nil {
		println("Error finding device:", err)
		return
	}

	for cur.Next(ctx) {
		var c map[string]any

		if err := cur.Decode(&c); err != nil {
			log.Printf("Error decoding cluster document: %v", err)
			continue
		}

		clusterID, ok := c["_id"].(string)
		if !ok || clusterID == "" {
			log.Printf("Skipping document missing or invalid _id: %v", c["_id"])
			continue
		}

		updateSet, found := updates[clusterID]
		if !found {
			continue // no updates for this cluster
		}

		// --- Update cluster description and name if available ---
		updateFields := bson.M{}
		if clusterDesc, ok := updateSet["description"]; ok {
			updateFields["description"] = clusterDesc
		}
		if clusterName, ok := updateSet["name"]; ok {
			updateFields["name"] = clusterName
		}
		if len(updateFields) > 0 {
			filter := bson.M{"_id": clusterID}
			update := bson.M{"$set": updateFields}

			res, err := db.ClusterDefinition(partnerId).UpdateOne(ctx, filter, update)
			if err != nil {
				log.Printf("Failed to update cluster %s description/name: %v", clusterID, err)
			} else if res.ModifiedCount > 0 {
				log.Printf("Updated cluster %s description/name", clusterID)
			}
		}

		// --- Update attribute descriptions and names ---
		for attrID, attrDataRaw := range updateSet {

			if attrID == "description" || attrID == "name" {
				continue // skip cluster-level keys
			}

			// attrDataRaw is expected to be a map[string]string with keys "description" and "name"
			attrData, ok := attrDataRaw.(map[string]any)
			if !ok {
				continue
			}

			// update both name and description if present
			attrUpdate := bson.M{}
			if desc, ok := attrData["description"]; ok {
				attrUpdate["attributes.$[elem].description"] = desc
			}
			if name, ok := attrData["name"]; ok {
				attrUpdate["attributes.$[elem].name"] = name
			}
			if len(attrUpdate) == 0 {
				continue
			}

			filter := bson.M{"_id": clusterID}
			update := bson.M{"$set": attrUpdate}
			arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
				Filters: []interface{}{bson.M{"elem.id": attrID}},
			})

			res, err := db.ClusterDefinition(partnerId).UpdateOne(ctx, filter, update, arrayFilters)
			if err != nil {
				log.Printf("Failed to update attribute %s in cluster %s: %v", attrID, clusterID, err)
			} else if res.ModifiedCount > 0 {
				log.Printf("Updated attribute %s in cluster %s", attrID, clusterID)
			}
		}
	}
}
