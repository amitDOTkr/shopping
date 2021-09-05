package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID               primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name             string               `json:"name,omitempty" bson:"name,omitempty"`
	ShortDescription string               `json:"shortDescription,omitempty" bson:"shortDescription,omitempty"`
	Description      string               `json:"description,omitempty" bson:"description,omitempty"`
	Slug             string               `json:"slug,omitempty" bson:"slug,omitempty"`
	Categories       []primitive.ObjectID `json:"categories,omitempty" bson:"categories,omitempty"`
	Price            float64              `json:"price,omitempty" bson:"price,omitempty"`
	IsActive         bool                 `json:"isActive" bson:"isActive"`
	FeaturedImage    string               `json:"featuredImage,omitempty" bson:"featuredImage,omitempty"`
	Images           []string             `json:"images,omitempty" bson:"images,omitempty"`
	Tags             []string             `json:"tags,omitempty" bson:"tags,omitempty"`
	SellerId         primitive.ObjectID   `json:"sellerId,omitempty" bson:"sellerId,omitempty"`
	Reviews          []string             `json:"reviews,omitempty" bson:"reviews,omitempty"`
	CreatedAt        time.Time            `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	ModifiedAt       time.Time            `json:"modifiedAt,omitempty" bson:"modifiedAt,omitempty"`
}
