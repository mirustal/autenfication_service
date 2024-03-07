package routes

import (
	"context"
	"service/app/models"



	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = ("medods_Task1")

func AuthHandler(c *fiber.Ctx) error{
	guid := c.Params("guid")

	_, err := storage.SearchTokenByGuid(context.Background(), guid) 
	if err == nil {
		p := new(models.UserCookie)
		if err := c.CookieParser(p); err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
		}

		claims, errParse := parseAccessToken(p.AccessToken)
		
		if errParse != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "not access",
			})
		}
		
		sub, ok := (*claims)["sub"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token structure",
			})
		}

		if sub != guid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "not access",
			})
		}
		return c.Next()
	}
	
	refresh, err :=  storage.CreateRefreshToken(context.Background(), guid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "error on created token",
		})
	}

	t, err := CreateAccessToken(guid)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	SetCookie(c, "accesst", t)
	SetCookie(c, "refresht", refresh)

	return c.Next()

	
}

func SetCookie(c *fiber.Ctx, name, value string) {
	cookie := new(fiber.Cookie)
	cookie.Name = name
	// cookie.Expires = time.Now().Add(1 * time.Second)
	cookie.Value = value
	c.Cookie(cookie)
}

func CreateAccessToken(guid string) (string, error) {
	payload := jwt.MapClaims{
		"sub": guid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, payload)

	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}


func parseAccessToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidType
		}
		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, jwt.ErrTokenInvalidClaims
	}
}
