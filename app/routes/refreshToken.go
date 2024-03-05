package routes

import (
	"context"


	"github.com/gofiber/fiber/v2"
)

type requestRefreshDTO struct {
	refreshToken string `json:"refreshtoken"`
}

type response struct {
	Description  string `json:"description,omitempty"`
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}


func RefreshToken(c *fiber.Ctx) error {
	println("Handle CreateRoom")

	queryInfo := new(requestRefreshDTO)
	if err := c.BodyParser(&queryInfo); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

	answer, err :=  storage.UpdateToken(context.Background(), queryInfo.refreshToken)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "error on refresh token",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"asnwer": answer,
	})

}