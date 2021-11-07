package sellerservice

import (
	"log"

	"github.com/amitdotkr/go-shopping/src/entities"
	"github.com/amitdotkr/go-shopping/src/global"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	data := &entities.Seller{}
	res := SellerCollection.FindOne(global.Ctx, bson.M{"email": seller.Email})
	if err := res.Decode(data); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
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

	at, rt, err := CreateTokenPair(data.ID.Hex(), "seller")
	if err != nil {
		log.Printf("err: %v", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":            data.ID,
		"email":         seller.Email,
		"access_token":  at,
		"refresh_token": rt,
	})
}

// Regerate Token Using Refresh Tokens

func RegenerateToken(c *fiber.Ctx) error {

	claims, err := global.TokenClaims(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": entities.Error{Type: "Authentication Error", Detail: err.Error()},
		})
	}

	claimType := claims["type"].(string)
	if claimType != "refresh" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": entities.Error{Type: "Authentication Error:", Detail: "Token is not access type"},
		})
	}

	userId := claims["user"].(string)
	keyId := claims["kid"].(string)

	kid, err := primitive.ObjectIDFromHex(keyId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": entities.Error{Type: "Object Id Error", Detail: err.Error()},
		})
	}

	token := global.ExtractToken(c)

	data := &entities.Refreshtoken{}
	filter := bson.M{"_id": kid}
	if err := global.RefreshCollection.FindOne(global.Ctx, filter).Decode(data); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": entities.Error{Type: "Not Found", Detail: err.Error()},
		})
	}

	if token != data.Token {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": entities.Error{Type: "Token Not Found", Detail: "Refresh Token is not available in DB"},
		})
	}

	at, rt, err := CreateTokenPairUsingRefreshToken(userId, "seller", keyId)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error: ": entities.Error{
				Type:   "Token Parsing Error",
				Detail: err.Error()},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  at,
		"refresh_token": rt,
	})

}

func CreateTokenPair(userid, role string) (string, string, error) {
	at, err := global.CreateAccessToken(userid, role)
	if err != nil {
		return "", "", err
	}
	rt, err := global.CreateRefreshToken(userid, role)
	if err != nil {
		return "", "", err
	}

	return at, rt, nil
}

func CreateTokenPairUsingRefreshToken(userid, role, keyId string) (string, string, error) {
	at, err := global.CreateAccessToken(userid, role)
	if err != nil {
		return "", "", err
	}
	rt, err := global.RegenerateRefreshToken(userid, role, keyId)
	if err != nil {
		return "", "", err
	}

	return at, rt, nil
}

func SellerSigninGo(c *fiber.Ctx) error {

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

	data := &entities.Seller{}
	res := SellerCollection.FindOne(global.Ctx, bson.M{"email": seller.Email})
	if err := res.Decode(data); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
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

	at, rt := CreateTokenPairGo(data.ID.Hex(), "seller")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":            data.ID,
		"email":         seller.Email,
		"access_token":  at,
		"refresh_token": rt,
	})
}

func CreateTokenPairGo(userid, role string) (string, string) {

	atch := make(chan string)
	rtch := make(chan string)
	go global.CreateAccessTokenGo(userid, role, atch)

	go global.CreateRefreshTokenGo(userid, role, rtch)

	return <-atch, <-rtch

}
