package productservice

import (
	"github.com/amitdotkr/go/shopping/src/global"
	"go.mongodb.org/mongo-driver/mongo"
)

var ProductCollection mongo.Collection

func init() {
	ProductCollection = *global.DB.Collection("products")
}
