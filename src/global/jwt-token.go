package global

import (
	"fmt"
	"log"

	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func CreateToken(userid string, role string) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = userid
	atClaims["role"] = role
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(JWT_Access_Secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ExtractToken(c *fiber.Ctx) string {
	bearToken := c.Get("Authorization")
	if strings.HasPrefix(bearToken, "Bearer ") {
		tokenStringFinal := strings.TrimPrefix(bearToken, "Bearer ")
		log.Printf("%v", tokenStringFinal)
		return tokenStringFinal
	}
	return ""
}

func VerifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			alg := token.Header["alg"]
			log.Printf("unexpected signing method: %v", alg)
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_Access_Secret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(c *fiber.Ctx) error {
	token, err := VerifyToken(c)
	log.Printf("token: %v", token)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}
