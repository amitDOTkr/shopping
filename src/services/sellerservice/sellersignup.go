package sellerservice

import (
	"time"

	"github.com/amitdotkr/go-shopping/src/entities"
	"github.com/amitdotkr/go-shopping/src/global"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SellerSignup(c *fiber.Ctx) error {

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
		return c.Status(fiber.StatusMultiStatus).JSON(fiber.Map{
			"error": entities.Error{Type: "DataBase Error", Detail: err.Error()},
		})
	}
	if isExist {
		return c.Status(fiber.StatusMultiStatus).JSON(fiber.Map{
			"error": entities.Error{Type: "Email Already Exist"},
		})
	}

	hashedpassword, err := global.HashPassword(seller.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": entities.Error{Type: "Password Hashing", Detail: err.Error()},
		})
	}

	sellerData := entities.Seller{
		Name:       seller.Name,
		Email:      seller.Email,
		Password:   hashedpassword,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	sellerRes, err := SellerCollection.InsertOne(global.Ctx, sellerData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": entities.Error{Type: "Database Error", Detail: err.Error()},
		})
	}

	oid, _ := sellerRes.InsertedID.(primitive.ObjectID)

	if err := CreateTokenPairGo(c, oid.Hex(), "seller"); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": entities.Error{
				Type:   "Token Generation Error",
				Detail: err.Error()},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"seller": entities.Seller{
			ID:              oid,
			Name:            seller.Name,
			Email:           seller.Email,
			IsEmailVerified: seller.IsEmailVerified,
			ProfileImage:    seller.ProfileImage,
			IsActice:        seller.IsActice,
		},
	})

}

func IsEmailAlreadyExist(email string) (bool, error) {
	count, err := SellerCollection.CountDocuments(global.Ctx, bson.M{"email": email}, options.Count())
	if err != nil {
		return true, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
