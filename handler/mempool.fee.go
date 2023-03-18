package handler

import (
	mempoolfee "github.com/NonsoAmadi10/mempool-fee/mempool-fee"
	"github.com/gofiber/fiber/v2"
)

type Fee struct {
	High    float64 `json:"high,omitempty"`
	Low     float64 `json:"low,omitempty"`
	Normal  float64 `json:"normal,omitempty"`
	BestFee float64 `json:"best_fee,omitempty"`
}

func EstmateBestFee(c *fiber.Ctx) error {
	bestFee := mempoolfee.GetBestFee()

	response := &Fee{
		BestFee: bestFee,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
