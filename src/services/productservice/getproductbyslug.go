package productservice

import (
	"github.com/amitdotkr/go/shopping/src/entities"
	"github.com/amitdotkr/go/shopping/src/global"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetProductBySlug(c *fiber.Ctx) error {

	slug := c.Params("slug")

	filter := bson.M{"slug": slug}
	data := &entities.Product{}
	res := ProductCollection.FindOne(global.Ctx, filter)
	if err := res.Decode(data); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": entities.Error{
				Type:   "Database/Json Error",
				Detail: err.Error()},
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
