package productservice

import (
	"log"

	"github.com/amitdotkr/go-shopping/src/entities"
	"github.com/amitdotkr/go-shopping/src/global"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ListProducts(c *fiber.Ctx, filter primitive.M) error {
	// log.Print(sid)
	products := []entities.Product{}

	data := &entities.Product{}
	cur, err := ProductCollection.Find(global.Ctx, filter)
	if err != nil {
		log.Printf("cur error: %v", err)
	}

	defer cur.Close(global.Ctx)

	for cur.Next(global.Ctx) {
		err := cur.Decode(data)
		if err != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": entities.Error{Type: "Database Error", Detail: err.Error()},
			})
		}

		products = append(products, *data)

	}
	if err := cur.Err(); err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": entities.Error{Type: "Database Error", Detail: err.Error()},
		})
	}
	c.Status(fiber.StatusOK).JSON(products)
	return nil
}
