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

func EstimateImprovedBestFee(c *fiber.Ctx) error {
	bestFee := mempoolfee.GetImprovedBestFee()

	response := &Fee{
		BestFee: bestFee,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func EstimateFees(c *fiber.Ctx) error {
	high, best, low, err := mempoolfee.GetPriorityFees()

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := &Fee{
		BestFee: best,
		Low:     low,
		High:    high,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func EsitmateFastestFee(c *fiber.Ctx) error {
	fast, _, _, err := mempoolfee.GetPriorityFees()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	response := &Fee{
		BestFee: fast,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func EstimateHalfHourFee(c *fiber.Ctx) error {
	bestFee := mempoolfee.GetHalfHourFee()

	response := &Fee{
		BestFee: bestFee,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
