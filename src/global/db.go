package global

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB mongo.Database
var Ctx context.Context

func ConnectToMongo() {

	client, err := mongo.NewClient(options.Client().ApplyURI(DB_URL))
	if err != nil {
		log.Fatalln("mongodb error")
		log.Fatal(err)
	}
	Ctx = context.Background()

	ctx, cancel := context.WithTimeout(Ctx, 10*time.Second)
	defer cancel()

	client.Connect(ctx)

	DB = *client.Database("shopping")
}

// func init() {
// 	ConnectToMongo()
// }
