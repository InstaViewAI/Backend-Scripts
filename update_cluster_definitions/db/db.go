package db

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// partnerClientsOnce is a sync.Once object that atomically makes sure there is only one instantiation of itself.
// We will be using it to create Singletons.
var (
	partnerClientsOnce sync.Once
	partnerClients     *PartnerClients
)

type Database struct {
	Partner          string `yaml:"partner"`
	ConnectionString string `yaml:"connection_string"`
}

// PartnerClients is a wrapper type on map of partners and mongo db clients.
// Each key will contain a separate partner replica set connection.
type PartnerClients struct {
	mu      sync.RWMutex
	clients map[string]*mongo.Client
}

func NewPartnerDBs(dbs []Database) *PartnerClients {
	partnerClientsOnce.Do(
		func() {
			// Iterate over the db configurations.
			// Each partner will have its own set of db configurations.
			for _, dbCfg := range dbs {
				// Use the SetServerAPIOptions() method to set the Stable API version to 1
				serverAPI := options.ServerAPI(options.ServerAPIVersion1)
				opts := options.Client().ApplyURI(dbCfg.ConnectionString).SetServerAPIOptions(serverAPI)

				// Create a new client and connect to the server
				client, err := mongo.Connect(context.TODO(), opts)
				if err != nil {
					panic(err)
				}

				// Send a ping to confirm a successful connection
				var result bson.M
				if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
					panic(err)
				} else {
					println(fmt.Sprintf("Successfully connected to MongoDB: %s", dbCfg.Partner))
				}

				// DB name will match the service name.
				// Fetching a service db handle
				if partnerClients == nil {
					partnerClients = &PartnerClients{
						clients: make(map[string]*mongo.Client),
					}
				}
				partnerClients.Set(dbCfg.Partner, client)
			}
		},
	)

	return partnerClients
}

func (pc *PartnerClients) Set(key string, c *mongo.Client) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	if pc.clients == nil {
		pc.clients = make(map[string]*mongo.Client)
	}
	pc.clients[key] = c
}

// Client is a getter that retrieves a partner db handler by key from the map.
func (pc *PartnerClients) Client(pID string) *mongo.Client {

	client, ok := pc.clients[pID]
	if !ok {
		panic(fmt.Errorf("no client found for given partner: %s", pID))
	}

	return client
}

func ClusterDefinition(pid string) *mongo.Collection {
	return partnerClients.Client(pid).Database("zeus").Collection("cluster_definition")
}
