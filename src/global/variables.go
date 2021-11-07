package global

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Debugger_Val bool
var DB_URL string
var PRVKEY_LOC string
var PUBKEY_LOC string

func settingVariable() {

	// Using .env files during production
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Setting Enviroment Variables
	debugger, err := strconv.ParseBool(os.Getenv("Debugger_Val"))
	if err != nil {
		log.Println("Enviroment Variable is not set. OR Unable to Parse Debuuger_Value.")
	}

	Debugger_Val = debugger
	DB_URL = os.Getenv("MongoDB_URL")
	PRVKEY_LOC = os.Getenv("PRVKEY_LOC")
	PUBKEY_LOC = os.Getenv("PUBKEY_LOC")

	// Initializing Database
	ConnectToMongo()
	Refresh()
}

func init() {
	settingVariable()
}
