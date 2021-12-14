package productservice

import (
	"time"

	"github.com/amitdotkr/go-shopping/src/entities"
	"github.com/amitdotkr/go-shopping/src/global"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddFeaturedImageToProduct(c *fiber.Ctx) error {

	userId, err := global.ValidatingUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": entities.Error{Type: "Authentication Error", Detail: err.Error()},
		})
	}

	var product entities.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"error": entities.Error{Type: "JSON Error", Detail: err.Error()},
		})
	}

	data := &entities.Product{}
	filter := bson.M{"_id": product.ID}
	if err := ProductCollection.FindOne(global.Ctx, filter).Decode(data); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": entities.Error{Type: "Not Found", Detail: err.Error()},
		})
	}

	if data.SellerId.Hex() != userId {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": entities.Error{Type: "Unauthorized Access", Detail: "SellerIDs does Not match"},
		})
	}

	modifiedTime := time.Now()

	update := bson.M{
		"$set": bson.M{
			"featuredImage": product.FeaturedImage,
			"modifiedAt":    modifiedTime,
		},
	}
	opts := options.Update().SetUpsert(true)
	ProductCollection.UpdateOne(global.Ctx, filter, update, opts)

	return c.SendStatus(fiber.StatusOK)
}

func AddImagesToProduct(c *fiber.Ctx) error {

	userId, err := global.ValidatingUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": entities.Error{Type: "Authentication Error", Detail: err.Error()},
		})
	}

	var product entities.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"error": entities.Error{Type: "JSON Error", Detail: err.Error()},
		})
	}

	data := &entities.Product{}
	filter := bson.M{"_id": product.ID}
	if err := ProductCollection.FindOne(global.Ctx, filter).Decode(data); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": entities.Error{Type: "Not Found", Detail: err.Error()},
		})
	}

	if data.SellerId.Hex() != userId {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": entities.Error{Type: "Unauthorized Access", Detail: "SellerIDs does Not match"},
		})
	}

	modifiedTime := time.Now()

	update := bson.M{
		"$set": bson.M{
			"images":     product.Images,
			"modifiedAt": modifiedTime,
		},
	}
	opts := options.Update().SetUpsert(true)
	ProductCollection.UpdateOne(global.Ctx, filter, update, opts)

	return c.SendStatus(fiber.StatusOK)
}
