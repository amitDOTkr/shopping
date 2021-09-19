package categoryservice

import (
	"strings"
	"time"

	"github.com/amitdotkr/go-shopping/src/entities"
	"github.com/amitdotkr/go-shopping/src/global"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddCategory(c *fiber.Ctx) error {

	var category entities.Category

	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"error": entities.Error{Type: "JSON Error", Detail: err.Error()},
		})
	}

	categoryName := strings.TrimSpace(category.Name)
	isCategoryNameExist, err := IsCategoryNameAlreadyExist(categoryName)
	if err != nil {
		return c.Status(fiber.StatusFound).JSON(fiber.Map{
			"error": entities.Error{Type: "Database/JSON Error", Detail: err.Error()},
		})
	}
	if isCategoryNameExist {
		return c.Status(fiber.StatusFound).JSON(fiber.Map{
			"error": entities.Error{Type: "Category Name Is Already Exist"},
		})
	}

	categorySlug := strings.TrimSpace(category.Slug)
	isSlugExist, err := IsSlugAlreadyExist(categorySlug)
	if err != nil {
		return c.Status(fiber.StatusFound).JSON(fiber.Map{
			"error": entities.Error{Type: "Database/JSON Error", Detail: err.Error()},
		})
	}
	if isSlugExist {
		return c.Status(fiber.StatusFound).JSON(fiber.Map{
			"error": entities.Error{Type: "Slug Already Exist"},
		})
	}

	category.Name = categoryName
	category.Slug = categorySlug
	category.CreatedAt = time.Now()
	category.ModifiedAt = time.Now()

	res, err := CategoryCollection.InsertOne(global.Ctx, category)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": entities.Error{Type: "Database Error", Detail: err.Error()},
		})
	}
	oid, added := res.InsertedID.(primitive.ObjectID)
	if !added {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"error": entities.Error{Type: "Database Error", Detail: "Object Id not valid"},
		})
	}
	category.ID = oid
	return c.Status(fiber.StatusCreated).JSON(category)
}

func IsSlugAlreadyExist(slug string) (bool, error) {
	count, err := CategoryCollection.CountDocuments(global.Ctx, bson.M{"slug": slug}, options.Count())
	if err != nil {
		return true, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func IsCategoryNameAlreadyExist(name string) (bool, error) {
	count, err := CategoryCollection.CountDocuments(global.Ctx, bson.M{"name": name}, options.Count())
	if err != nil {
		return true, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
