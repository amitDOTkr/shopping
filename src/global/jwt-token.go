package global

import (
	"fmt"
	"log"

	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func CreateToken(userid string, role string) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["id"] = userid
	atClaims["role"] = role
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString(JWT_Access_Secret)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ExtractToken(c *fiber.Ctx) string {
	bearToken := c.Get("Authorization")
	if strings.HasPrefix(bearToken, "Bearer ") {
		tokenStringFinal := strings.TrimPrefix(bearToken, "Bearer ")
		// log.Printf("%v", tokenStringFinal)
		return tokenStringFinal
	}
	return ""
}

func ParsingToken(c *fiber.Ctx) (*jwt.Token, error) {
	extractedtoken := ExtractToken(c)

	token, err := jwt.Parse(extractedtoken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			alg := token.Header["alg"]
			log.Printf("unexpected signing method: %v", alg)
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JWT_Access_Secret, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenClaims(c *fiber.Ctx) (jwt.MapClaims, error) {
	token, err := ParsingToken(c)
	if err != nil {
		return nil, err
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims, nil
}
