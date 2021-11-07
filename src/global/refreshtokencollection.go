package global

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var RefreshCollection mongo.Collection

func Refresh() {
	RefreshCollection = *DB.Collection("refresh")
}
