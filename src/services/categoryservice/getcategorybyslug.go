package categoryservice

import (
	"github.com/amitdotkr/go/shopping/src/entities"
	"github.com/amitdotkr/go/shopping/src/global"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetCategoryBySlug(c *fiber.Ctx) error {

	slug := c.Params("slug")

	filter := bson.M{"slug": slug}
	data := &entities.Category{}
	res := CategoryCollection.FindOne(global.Ctx, filter)
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
	return c.Status(fiber.StatusOK).JSON(data)
}
