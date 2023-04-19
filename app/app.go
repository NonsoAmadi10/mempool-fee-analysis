package app

import (
	"github.com/NonsoAmadi10/mempool-fee/handler"
	"github.com/gofiber/fiber/v2"
)

func App() *fiber.App {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/best_fee", handler.EstmateBestFee)

	app.Get("/recommended-fees", handler.EstimateFees)
	return app
}
