package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/go-playground/validator/v10"
)

type ValidateErrorResponse struct {
	FailedField string      `json:"failed_field"`
	Tag         string      `json:"tag"`
	Params      string      `json:"params"`
	Value       interface{} `json:"value"`
}

var validate = validator.New()

func ValidateStruct(st interface{}) []*ValidateErrorResponse {
	var errors []*ValidateErrorResponse
	err := validate.Struct(st)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidateErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Params = err.Param()
			element.Value = err.Value()
			errors = append(errors, &element)
		}
	}
	return errors
}

type Job struct {
	Type   string `json:"type" validate:"required,min=3,max=32"`
	Salary int    `json:"salary" validate:"required,number"`
}

type User struct {
	Name  string `json:"name" validate:"required,min=3,max=32"`
	Email string `json:"email" validate:"required,email,min=6,max=32"`
	Job   Job    `json:"job" validate:"dive"`
}

// database!
var users = make(map[string]User)

func main() {

	app := fiber.New(fiber.Config{Prefork: true})

	// middlewares
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
	}))

	router := app.Group("/api/v1")

	router.Post("/users", func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		if errors := ValidateStruct(*user); errors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(errors)
		}

		if _, ok := users[user.Name]; ok {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "invalid Name"})
		}

		users[user.Name] = *user
		return c.Status(fiber.StatusCreated).JSON(user)
	})

	router.Get("/users/:name<minLen(3);maxLen(32)>", func(c *fiber.Ctx) error {
		user, ok := users[c.Params("name")]
		if !ok {
			return fiber.ErrNotFound
		}

		return c.JSON(user)
	})

	router.Delete("/users/:name", func(c *fiber.Ctx) error {
		if _, ok := users[c.Params("name")]; ok {
			delete(users, c.Params("name"))
			return c.SendStatus(fiber.StatusNoContent)

		} else {
			return fiber.ErrNotFound
		}
	})

	router.Patch("/users", func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		if errors := ValidateStruct(*user); errors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(errors)
		}

		if _, ok := users[user.Name]; !ok {
			return fiber.ErrNotFound
		}

		users[user.Name] = *user
		return c.JSON(user)

	})

	log.Fatal(app.Listen(":8080"))
}
