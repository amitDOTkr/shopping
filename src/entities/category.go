package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name             string             `json:"name,omitempty" bson:"name,omitempty"`
	ShortDescription string             `json:"short_description,omitempty" bson:"short_description,omitempty"`
	Description      string             `json:"description,omitempty" bson:"description,omitempty"`
	Slug             string             `json:"slug,omitempty" bson:"slug,omitempty"`
	IsParentCategory bool               `json:"is_parent_category" bson:"is_parent_category"`
	ParentCategory   primitive.ObjectID `json:"parent_category,omitempty" bson:"parent_category,omitempty"`
	SubCategories    []string           `json:"sub_categories,omitempty" bson:"sub_categories,omitempty"`
	IsActive         bool               `json:"is_active" bson:"is_active"`
	FeaturedImage    string             `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
	CreatedAt        time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	ModifiedAt       time.Time          `json:"modifiedAt,omitempty" bson:"modifiedAt,omitempty"`
}
