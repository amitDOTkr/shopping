package main

import (
	"github.com/amitdotkr/go-shopping/src/global"
	"github.com/amitdotkr/go-shopping/src/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"

	// "github.com/gofiber/fiber/v2/middleware/favicon"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	if global.Debugger_Val {
		app.Use(logger.New())
	}
	app.Use(cors.New())
	// app.Use(favicon.New(favicon.Config{
	// 	File: "./favicon.ico",
	// }))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	routes.Register(app)

	app.Listen(":3000")
}
