package sellerservice

import (
	"log"

	"github.com/amitdotkr/go-shopping/src/entities"
	"github.com/amitdotkr/go-shopping/src/global"
	"github.com/amitdotkr/go-shopping/src/services/productservice"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MyProducts(c *fiber.Ctx) error {
	sellerId, err := global.ValidatingUser(c)
	if err != nil {
		log.Printf("error in myproduct func")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": entities.Error{Type: "Authentication Error", Detail: err.Error()},
		})
	}
	// token := c.Cookies("access_token", "")

	// if token == "" {
	// 	log.Println("access token cookie is empty")
	// 	RegenerateToken(c)
	// }

	// claims, err := global.TokenClaims(token)
	// if err != nil {
	// 	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
	// 		"error": entities.Error{Type: "Authentication Error", Detail: err.Error()},
	// 	})
	// }

	// claimType := claims["type"].(string)
	// if claimType != "access" {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 		"error": entities.Error{Type: "Authentication Error:", Detail: "Token is not access type"},
	// 	})
	// }

	// sellerId := claims["id"].(string)
	sid, err := primitive.ObjectIDFromHex(sellerId)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": entities.Error{Type: "Seller Id Error", Detail: err.Error()},
		})
	}

	filter := bson.M{"sellerId": sid}
	productservice.ListProducts(c, filter)
	return nil
}
