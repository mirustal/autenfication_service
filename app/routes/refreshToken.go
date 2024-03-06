package routes

import (
	"context"
	"service/app/models"

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
	guid := c.Params("guid")
	p := new(models.UserCookie)

	if err := c.CookieParser(p); err != nil {
		return err
	}

	refreshToken, err := storage.SearchTokenByGuid(context.Background(), guid)
	if err != nil {
		return err
	}

	if refreshToken != p.RefreshToken {
		return c.Status(400).JSON(fiber.Map{
			"answer": "not access",
		})
	}

	refreshToken, err = storage.UpdateToken(context.Background(), guid)
	if err != nil {
		return err
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "refreshtoken"
	cookie.Value = refreshToken
	c.Cookie(cookie)

	return c.Status(200).JSON(fiber.Map{
		"AccessToken": p.AccessToken,
		"RefreshToken": refreshToken,
	})
}