package routes

import (
	"service/app/models"

	"github.com/gofiber/fiber/v2"
)



type responseDTO struct {
	Description  string `json:"description,omitempty"`
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`}



func GetToken(c *fiber.Ctx) error {
	p := new(models.UserCookie)

	if err := c.CookieParser(p); err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"AccessToken": p.AccessToken,
		"RefreshToken": p.RefreshToken,
	})
}