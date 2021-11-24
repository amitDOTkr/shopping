package productservice

import (
	"time"

	"github.com/amitdotkr/go-shopping/src/entities"
	"github.com/amitdotkr/go-shopping/src/global"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddProduct(c *fiber.Ctx) error {

	// var sellerId string

	userId, err := global.ValidatingUser(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": entities.Error{Type: "Authentication Error", Detail: err.Error()},
		})
	}

	// at := c.Cookies("access_token", "")
	// if at != "" {
	// 	claims, err := global.TokenClaims(at)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
	// 			"error": entities.Error{Type: "Authentication Error", Detail: err.Error()},
	// 		})
	// 	}
	// 	sellerId = claims["id"].(string)
	// }

	// if at == "" {
	// 	global.RegenerateTokenUsingRefreshToken(c)
	// 	rt := c.Cookies("refresh_token")
	// 	claims, err := global.TokenClaims(rt)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
	// 			"error": entities.Error{Type: "Authentication Error", Detail: err.Error()},
	// 		})
	// 	}
	// 	sellerId = claims["user"].(string)
	// }

	var product entities.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"error": entities.Error{Type: "JSON Error", Detail: err.Error()},
		})
	}

	isSlugExist, err := IsSlugAlreadyExist(product.Slug)
	if err != nil {
		return c.Status(fiber.StatusFound).JSON(fiber.Map{
			"error": entities.Error{Type: "Database/JSON Error", Detail: err.Error()},
		})
	}
	if isSlugExist {
		return c.Status(fiber.StatusFound).JSON(fiber.Map{
			"error": entities.Error{Type: "Slug Already Exist"},
		})
	}

	sid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": entities.Error{Type: "Seller Id Error", Detail: err.Error()},
		})
	}

	// productCategories := product.Categories

	product.SellerId = sid
	product.CreatedAt = time.Now()
	product.ModifiedAt = time.Now()

	res, err := ProductCollection.InsertOne(global.Ctx, product)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": entities.Error{Type: "Database Error", Detail: err.Error()},
		})
	}
	oid, added := res.InsertedID.(primitive.ObjectID)
	if !added {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"error": entities.Error{Type: "Database Error", Detail: "Object Id not valid"},
		})
	}
	product.ID = oid
	return c.Status(fiber.StatusCreated).JSON(product)
}

func IsSlugAlreadyExist(slug string) (bool, error) {
	count, err := ProductCollection.CountDocuments(global.Ctx, bson.M{"slug": slug}, options.Count())
	if err != nil {
		return true, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
