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

	check, err := storage.ValidateRefreshToken(context.Background(), guid, p.RefreshToken)
	if !check || err != nil  {
		return err
	}

	accessToken, err := CreateAccessToken(guid)
	if err != nil  {
		return err
	}

	refreshToken, err := storage.UpdateRefreshToken(context.Background(), guid)
	if err != nil {
		return err
	}

	SetCookie(c, "accesst", accessToken)
	SetCookie(c, "refresht", refreshToken)

	return c.Status(200).JSON(fiber.Map{
		"AccessToken": p.AccessToken,
		"RefreshToken": p.RefreshToken,
	})
}



