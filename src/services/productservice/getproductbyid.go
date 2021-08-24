package productservice

import (
	"github.com/amitdotkr/go/shopping/src/entities"
	"github.com/amitdotkr/go/shopping/src/global"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetProductByID(c *fiber.Ctx) error {
	prodId := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(prodId)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": entities.Error{Type: "Object Id Error", Detail: err.Error()},
		})
	}
	data := &entities.Product{}
	filter := bson.M{"_id": oid}
	if err := ProductCollection.FindOne(global.Ctx, filter).Decode(data); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": entities.Error{Type: "Not Found", Detail: err.Error()},
		})
	}
	if !data.IsActive {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": entities.Error{Type: "Unlisted Product",
				Detail: "Product Is Not Listed Yet"},
		})
	}
	return c.Status(fiber.StatusOK).JSON(data)
}
