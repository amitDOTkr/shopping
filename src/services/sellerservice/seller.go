package sellerservice

import (
	"github.com/amitdotkr/go-shopping/src/global"
	"go.mongodb.org/mongo-driver/mongo"
)

var SellerCollection mongo.Collection

func init() {
	SellerCollection = *global.DB.Collection("sellers")
}
