package categoryservice

import (
	"github.com/amitdotkr/go/shopping/src/global"
	"go.mongodb.org/mongo-driver/mongo"
)

var CategoryCollection mongo.Collection

func init() {
	CategoryCollection = *global.DB.Collection("categories")
}
