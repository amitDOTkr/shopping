package routes

import (
	"github.com/amitdotkr/go/shopping/src/services/productservice"
	"github.com/amitdotkr/go/shopping/src/services/sellerservice"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"api_working": true,
		})
	})

	// api := app.Group("/api")
	// v1 := api.Group("/v1")
	// ProductApi := v1.Group("/products")

	// Products Routes
	products := app.Group("products")
	products.Post("/", productservice.AddProduct)
	products.Get("/:id", productservice.GetProductByID)
	products.Get("/slug/:slug", productservice.GetProductBySlug)

	seller := app.Group("seller")
	seller.Post("/signup", sellerservice.SellerSignup)
	seller.Post("/signin", sellerservice.SellerSignin)
}

// app.Post("/products", productroutes.AddProduct)
// app.Get("/products/:id", productroutes.GetProductByID)
// app.Get("/products/slug/:slug", productroutes.GetProductBySlug)
