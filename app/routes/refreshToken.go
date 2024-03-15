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

	dataUser := RequestAuthDTO{
		Guid:     guid,		
		Login:    login,
		Password: password,
	}

	check, err := storage.ValidateRefreshToken(context.Background(), dataUser.Guid, p.RefreshToken)
	if !check || err != nil  {
		return err
	}

	
	accessToken, err := CreateAccessToken(dataUser)
	if err != nil  {
		return err
	}

	refreshToken, err := storage.UpdateRefreshToken(context.Background(),dataUser.Guid)
	if err != nil {
		return err
	}

	SetCookie(c, "accesst", accessToken)
	SetCookie(c, "refresht", refreshToken)

	return c.Status(200).JSON(fiber.Map{
		"AccessToken": accessToken,
		"RefreshToken": refreshToken,
	})
}



