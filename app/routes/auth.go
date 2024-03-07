package routes

import (
	"context"
	"os"
	"service/app/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = os.Getenv("CONFIG_PATH") 
//jsonobj

var login = "miroslav"
var password = "godev"

type RequestAuthDTO struct {
	Guid     string `json:"guid"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func AuthHandler(c *fiber.Ctx) error {
	guid := c.Params("guid")
	dataUser := RequestAuthDTO{
		Guid:     guid,		
		Login:    login,
		Password: password,
	}

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

		err := validatePayloadtoken(dataUser, claims); 
		if !err {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token structure",
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

	t, err := CreateAccessToken(dataUser)
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
	// cookie.Expires = time.Now().Add(1 * time.house)
	cookie.Value = value
	c.Cookie(cookie)
}

func CreateAccessToken(dataUser RequestAuthDTO) (string, error) {
	expTime := time.Now().Add(1 * time.Hour).Unix()

	payload := jwt.MapClaims{
		"guid":     dataUser.Guid,
		"login":    dataUser.Login,
		"password": dataUser.Password,
		"exp":      expTime, 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, payload)

	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}


func parseAccessToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidType
		}
		return []byte(jwtSecretKey), nil
	})



	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	} else {
		return nil, jwt.ErrTokenInvalidClaims
	}
}



func validatePayloadtoken(dataUser RequestAuthDTO, claims jwt.MapClaims) bool {

	if claimsGuid, ok := claims["guid"].(string);
	!ok || claimsGuid != dataUser.Guid {
		return false
		}

	if claimsLogin, ok := claims["login"].(string);
	!ok || claimsLogin != dataUser.Login {
		return false
	}


	if claimsPassword, ok := claims["password"].(string); 
	!ok || claimsPassword != dataUser.Password {
		return false
	}


	return true

}