package global

import (
	"log"
	"os"
	"strconv"
	// "github.com/joho/godotenv"
)

var Debugger_Val bool
var DB_URL string
var PRVKEY_LOC string
var PUBKEY_LOC string
var AccessTokenTime int
var RefreshTokenTime int

func settingVariable() {

	// Using .env files during development, should be commented in production
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Println("Error loading .env file")
	// }

	// Setting Enviroment Variables
	debugger, err := strconv.ParseBool(os.Getenv("Debugger_Val"))
	if err != nil {
		log.Println("Enviroment Variable is not set OR Unable to Parse Debuuger_Value.")
	}

	access_token_time, err := strconv.Atoi(os.Getenv("AccessTokenTime"))
	if err != nil {
		log.Println("Enviroment Variable is not set OR Unable to Parse Access Token Time.")
	}

	refresh_token_time, err := strconv.Atoi(os.Getenv("RefreshTokenTime"))
	if err != nil {
		log.Println("Enviroment Variable is not set OR Unable to Parse Refresh Token Time.")
	}

	Debugger_Val = debugger
	AccessTokenTime = access_token_time
	RefreshTokenTime = refresh_token_time
	DB_URL = os.Getenv("MongoDB_URL")

	// For Production
	PRVKEY_LOC = "/etc/cert/id_rsa"
	PUBKEY_LOC = "/etc/cert/id_rsa.pub"

	// for Development
	// PRVKEY_LOC = "/home/amit/projects/go/shopping/src/cert/id_rsa"
	// PUBKEY_LOC = "/home/amit/projects/go/shopping/src/cert/id_rsa.pub"

	// Initializing Database
	ConnectToMongo()
	Refresh()
}

func init() {
	settingVariable()
}
