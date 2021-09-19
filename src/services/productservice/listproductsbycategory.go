package productservice

import (
	"github.com/amitdotkr/go-shopping/src/entities"
	"github.com/amitdotkr/go-shopping/src/global"
	"github.com/amitdotkr/go-shopping/src/services/categoryservice"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func ListProductsbyCategory(c *fiber.Ctx) error {

	categorySlug := c.Params("name")

	filter := bson.M{"slug": categorySlug}
	data := &entities.Category{}
	res := categoryservice.CategoryCollection.FindOne(global.Ctx, filter)
	if err := res.Decode(data); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": entities.Error{
				Type:   "Database/Json Error",
				Detail: err.Error()},
		})
	}
	if !data.IsActive {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": entities.Error{Type: "Unlisted Category",
				Detail: "Category Is Not Listed Yet"},
		})
	}
	filter = bson.M{"categories": data.ID, "isActive": true}

	ListProducts(c, filter)

	return nil
}
