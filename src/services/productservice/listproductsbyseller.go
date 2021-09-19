package productservice

import (
	"github.com/amitdotkr/go-shopping/src/entities"
	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ListProductsOfSeller(c *fiber.Ctx) error {
	sellerId := c.Params("sellerId")
	sid, err := primitive.ObjectIDFromHex(sellerId)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": entities.Error{Type: "Database Error", Detail: err.Error()},
		})
	}
	filter := bson.M{"sellerId": sid, "isActive": true}
	ListProducts(c, filter)
	return nil
}
