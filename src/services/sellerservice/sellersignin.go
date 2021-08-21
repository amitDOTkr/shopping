package sellerservice

import (
	"github.com/amitdotkr/go/shopping/src/entities"
	"github.com/amitdotkr/go/shopping/src/global"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func SellerSignin(c *fiber.Ctx) error {

	var seller entities.Seller

	validate := validator.New()

	if err := c.BodyParser(&seller); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": entities.Error{Type: "JSON Error", Detail: err.Error()},
		})
	}

	if err := validate.Struct(seller); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": entities.Error{Type: "Validation Error", Detail: err.Error()},
		})
	}

	isExist, err := IsEmailAlreadyExist(seller.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": entities.Error{
				Type:   "Database Error",
				Detail: err.Error()},
		})
	}
	if !isExist {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": entities.Error{
				Type:   "Authentication Error",
				Detail: "The Credentials you provided cannot be authenticated."},
		})
	}

	data := &entities.Seller{}
	res := SellerCollection.FindOne(global.Ctx, bson.M{"email": seller.Email})
	if err := res.Decode(data); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": entities.Error{
				Type:   "Database/Json Error",
				Detail: err.Error()},
		})
	}

	validUser := global.CheckPasswordHash(seller.Password, data.Password)
	if !validUser {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": entities.Error{
				Type:   "Authentication Error",
				Detail: "The Credentials you provided cannot be Authenticated."},
		})
	}

	token, _ := global.CreateToken(data.ID.Hex(), "seller")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":    data.ID,
		"email": seller.Email,
		"token": token,
	})
}
