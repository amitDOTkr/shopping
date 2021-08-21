package global

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func CreateToken(userid string, role string) (string, error) {
	var err error
	//Creating Access Token
	// os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	// atClaims["authorized"] = true
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

func Verify(c *fiber.Ctx) bool {
	isValid := false
	var errorMessage string
	errorMessage = ""
	tokenString := "Bearer knkjnkjnkknkjnknkjknkjkjknjk"
	// c.Request().Header.Header("Authorization")
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// return rsakeys[token.Header["kid"].(string)], nil
			return token, nil
		})
		if err != nil {
			errorMessage = err.Error()
		} else if !token.Valid {
			errorMessage = "Invalid token"
		} else if token.Header["alg"] == nil {
			errorMessage = "alg must be defined"
		} else if token.Claims.(jwt.MapClaims)["aud"] != "api://default" {
			errorMessage = "Invalid aud"
		} else if !strings.Contains(token.Claims.(jwt.MapClaims)["iss"].(string), os.Getenv("OKTA_DOMAIN")) {
			errorMessage = "Invalid iss"
		} else {
			isValid = true
		}
		if !isValid {
			c.SendStatus(fiber.StatusUnauthorized)
			log.Printf(errorMessage)
		}
	} else {
		c.SendStatus(fiber.StatusOK)
	}
	return isValid
}

// func ExtractToken(r *http.Request) string {
// 	bearToken := r.Header.Get("Authorization")
// 	//normally Authorization the_token_xxx
// 	strArr := strings.Split(bearToken, " ")
// 	if len(strArr) == 2 {
// 	   return strArr[1]
// 	}
// 	return ""
//   }

//   func VerifyToken(r *http.Request) (*jwt.Token, error) {
// 	tokenString := ExtractToken(r)
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 	   //Make sure that the token method conform to "SigningMethodHMAC"
// 	   if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 		  return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 	   }
// 	   return []byte(os.Getenv("ACCESS_SECRET")), nil
// 	})
// 	if err != nil {
// 	   return nil, err
// 	}
// 	return token, nil
//   }
