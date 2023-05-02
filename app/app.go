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

	v1 := app.Group("/api/v1")

	v2 := app.Group("/api/v2")
	v1.Get("/best-fee", handler.EstmateBestFee)

	v1.Get("/recommended-fees", handler.EstimateFees)
	v2.Get("/best-fee", handler.EstimateImprovedBestFee)
	v2.Get("/half-hour-fee", handler.EstimateHalfHourFee)
	return app
}
