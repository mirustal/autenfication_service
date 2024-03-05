package routes

import (
	"context"


	"github.com/gofiber/fiber/v2"
)

type requestDTO struct {
	Description string `json:"description,omitempty"`
}

type responseDTO struct {
	Description  string `json:"description,omitempty"`
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}


func GetToken(c *fiber.Ctx) error {
	println("Handle CreateRoom")

	queryInfo := new(requestDTO)
	if err := c.BodyParser(&queryInfo); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

	room, err :=  storage.CreateToken(context.Background())
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "error on created token",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"room": room,
	})

}