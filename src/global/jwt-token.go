package global

import (
	"fmt"
	"io/ioutil"
	"log"

	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ExtractToken(c *fiber.Ctx) string {
	bearToken := c.Get("Authorization")
	if strings.HasPrefix(bearToken, "Bearer ") {
		tokenStringFinal := strings.TrimPrefix(bearToken, "Bearer ")
		return tokenStringFinal
	}
	return ""
}

func TokenClaims(c *fiber.Ctx) (jwt.MapClaims, error) {
	token, err := ParsingToken(c)
	if err != nil {
		return nil, err
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims, nil
}

func ParsingToken(c *fiber.Ctx) (*jwt.Token, error) {
	extractedtoken := ExtractToken(c)
	pubKey, err := ioutil.ReadFile(PUBKEY_LOC)
	if err != nil {
		log.Fatalln(err)
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(pubKey)
	if err != nil {
		log.Printf("error: %v", err.Error())
	}

	token, err := jwt.Parse(extractedtoken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			alg := token.Header["alg"]
			log.Printf("unexpected signing method: %v", alg)
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func CreateAccessToken(userid string, role string) (string, error) {

	prvKey, err := ioutil.ReadFile(PRVKEY_LOC)
	if err != nil {
		log.Printf("err: %v", err.Error())
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	if err != nil {
		log.Printf("error: %v", err.Error())
	}

	atClaims := jwt.MapClaims{}
	atClaims["id"] = userid
	atClaims["role"] = role
	atClaims["type"] = "access"
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	at, err := jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims).SignedString(key)
	if err != nil {
		return "", err
	}

	return at, nil
}

func CreateRefreshToken(userid string, role string) (string, error) {

	kid := primitive.NewObjectID()
	keyId := kid.Hex()

	prvKey, err := ioutil.ReadFile(PRVKEY_LOC)
	if err != nil {
		log.Printf("err: %v", err.Error())
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	if err != nil {
		log.Printf("error: %v", err.Error())
	}

	tokenExpire := time.Now().Add(time.Hour * 600).Unix()

	rtClaims := jwt.MapClaims{}
	rtClaims["kid"] = keyId
	rtClaims["user"] = userid
	rtClaims["role"] = role
	rtClaims["type"] = "refresh"
	rtClaims["exp"] = tokenExpire

	rt, err := jwt.NewWithClaims(jwt.SigningMethodRS256, rtClaims).SignedString(key)
	if err != nil {
		return "", err
	}

	if err := RefreshTokenInDatabase(userid, tokenExpire, kid, rt); err != nil {
		return "", nil
	}

	return rt, nil
}

// Regenerate Refresh Token Using Refresh Token

func RegenerateRefreshToken(userid, role, keyId string) (string, error) {

	kid, err := primitive.ObjectIDFromHex(keyId)
	if err != nil {
		return "", err
	}

	prvKey, err := ioutil.ReadFile(PRVKEY_LOC)
	if err != nil {
		log.Printf("err: %v", err.Error())
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	if err != nil {
		log.Printf("error: %v", err.Error())
	}

	tokenExpire := time.Now().Add(time.Hour * 600).Unix()

	rtClaims := jwt.MapClaims{}
	rtClaims["kid"] = keyId
	rtClaims["user"] = userid
	rtClaims["role"] = role
	rtClaims["type"] = "refresh"
	rtClaims["exp"] = tokenExpire

	rt, err := jwt.NewWithClaims(jwt.SigningMethodRS256, rtClaims).SignedString(key)
	if err != nil {
		return "", err
	}

	if err := RefreshTokenInDatabase(userid, tokenExpire, kid, rt); err != nil {
		return "", nil
	}

	return rt, nil
}

// Saving Refresh Token in Database

func RefreshTokenInDatabase(userid string, tokenExpire int64, kid primitive.ObjectID, rt string) error {
	uid, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return err
	}

	t := time.Unix(tokenExpire, 0)

	filter := bson.M{"_id": kid}

	update := bson.M{
		"$set": bson.M{
			"token":    rt,
			"user":     uid,
			"expireAt": t,
		},
	}
	opts := options.Update().SetUpsert(true)
	RefreshCollection.UpdateOne(Ctx, filter, update, opts)

	return nil
}

func CreateAccessTokenGo(userid string, role string, atch chan string) {

	defer close(atch)

	prvKey, err := ioutil.ReadFile(PRVKEY_LOC)
	if err != nil {
		log.Printf("err: %v", err.Error())
		atch <- ""
		close(atch)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	if err != nil {
		log.Printf("error: %v", err.Error())
		atch <- ""
		close(atch)
	}

	atClaims := jwt.MapClaims{}
	atClaims["id"] = userid
	atClaims["role"] = role
	atClaims["type"] = "access"
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	at, err := jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims).SignedString(key)
	if err != nil {
		log.Printf("error: %v", err.Error())
		atch <- ""
		close(atch)
	}

	atch <- at
}

func CreateRefreshTokenGo(userid string, role string, rtch chan string) {

	defer close(rtch)

	kid := primitive.NewObjectID()
	keyId := kid.Hex()

	prvKey, err := ioutil.ReadFile(PRVKEY_LOC)
	if err != nil {
		log.Printf("err: %v", err.Error())
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	if err != nil {
		log.Printf("error: %v", err.Error())
	}

	tokenExpire := time.Now().Add(time.Hour * 600).Unix()

	rtClaims := jwt.MapClaims{}
	rtClaims["kid"] = keyId
	rtClaims["user"] = userid
	rtClaims["role"] = role
	rtClaims["type"] = "refresh"
	rtClaims["exp"] = tokenExpire

	rt, err := jwt.NewWithClaims(jwt.SigningMethodRS256, rtClaims).SignedString(key)
	if err != nil {
		log.Printf("error: %v", err.Error())
		rtch <- ""
		close(rtch)
	}

	if err := RefreshTokenInDatabase(userid, tokenExpire, kid, rt); err != nil {
		log.Printf("error: %v", err.Error())

		rtch <- ""
		close(rtch)
	}
	rtch <- rt

}
