package sellerservice

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/amitdotkr/go-shopping/src/entities"
	"github.com/amitdotkr/go-shopping/src/global"
	"github.com/gofiber/fiber/v2"
)

func UploadImages(c *fiber.Ctx) error {
	userId, err := global.ValidatingUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": entities.Error{Type: "Authentication Error", Detail: err.Error()},
		})
	}

	path := filepath.Join(".", "images", userId)

	if err = os.MkdirAll(path, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": entities.Error{Type: "Uploading Error", Detail: err.Error()},
		})
	}
	if form, err := c.MultipartForm(); err == nil {

		files := form.File["Image"]

		for _, file := range files {
			// fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])

			if err := c.SaveFile(file, fmt.Sprintf(path+"/%s", file.Filename)); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": entities.Error{Type: "Uploading Error", Detail: err.Error()},
				})
			}
		}
	}
	return c.SendStatus(fiber.StatusOK)
}
