package customerservice

import (
	"github.com/amitdotkr/go/shopping/src/global"
	"go.mongodb.org/mongo-driver/mongo"
)

var CustomerCollection mongo.Collection

func init() {
	CustomerCollection = *global.DB.Collection("customers")
}
