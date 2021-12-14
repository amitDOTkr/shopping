package routes

import (
	"github.com/amitdotkr/go-shopping/src/services/categoryservice"
	"github.com/amitdotkr/go-shopping/src/services/productservice"
	"github.com/amitdotkr/go-shopping/src/services/sellerservice"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"api_working": true,
		})
	})

	// api := app.Group("api")

	// Products Routes
	products := app.Group("products")
	// products := api.Group("products")
	products.Post("/", productservice.AddProduct)
	products.Get("/:id", productservice.GetProductByID)
	products.Get("/slug/:slug", productservice.GetProductBySlug)
	products.Get("/category/:name", productservice.ListProductsbyCategory)
	products.Get("/seller/:sellerId", productservice.ListProductsOfSeller)
	products.Post("/featuredimage", productservice.AddFeaturedImageToProduct)
	products.Post("/images", productservice.AddImagesToProduct)

	// Categories Routes
	categories := app.Group("categories")
	// categories := api.Group("categories")
	categories.Post("/", categoryservice.AddCategory)
	categories.Get("/id/:id", categoryservice.GetCategoryByID)
	categories.Get("/:slug", categoryservice.GetCategoryBySlug)

	seller := app.Group("seller")
	// seller := api.Group("seller")
	seller.Post("/signup", sellerservice.SellerSignup)
	seller.Post("/signin", sellerservice.SellerSigninGo)
	seller.Get("/myproducts", sellerservice.MyProducts)
	seller.Post("/regeneratetoken", sellerservice.RegenerateToken)
	seller.Post("/uploadimages", sellerservice.UploadImages)

}
