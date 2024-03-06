package routes

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"service/app/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)




// var jwtSecretKey = []byte("medods_Task1")
const privateKeyPath = "./key.pem"


func AuthHandler(c *fiber.Ctx) error{
	guid := c.Params("guid")

	_, err := storage.SearchTokenByGuid(context.TODO(), guid) 
	if err == nil {
		p := new(models.UserCookie)
		if err := c.CookieParser(p); err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
		}

		fmt.Printf(p.AccessToken)
		
		claims, errParse := ParseToken(p.AccessToken)
		
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
	

	refresh, err :=  storage.CreateToken(context.Background(), guid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "error on created token",
		})
	}

	t, err := CreateToken(guid)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	
	setCookie(c, "accesstoken", t)
	setCookie(c, "refreshtoken", refresh)

	answer := models.AccessResponse{
		Access:  t,
		Refresh: refresh,
	}

	return c.Status(200).JSON(fiber.Map{
		"answer": answer,
	})
}

func setCookie(c *fiber.Ctx, name, value string) {
	cookie := new(fiber.Cookie)
	cookie.Name = name
	cookie.Value = value
	c.Cookie(cookie)
}

func CreateKeys() (*rsa.PrivateKey, *rsa.PublicKey) {
	if _, err := os.Stat(privateKeyPath); err == nil {
		return nil, nil
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil
	}
	publicKey := &privateKey.PublicKey


	privKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	err = os.WriteFile(privateKeyPath, privKeyPEM, 0600)
	if err != nil {
		return nil, nil
	}

	return privateKey, publicKey
}

func GetPrivateKey() *rsa.PrivateKey {
	privKeyPEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil
	}

	block, _ := pem.Decode(privKeyPEM)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil
	}

	return privateKey
}

func GetPublicKeyByPrivateKey(privateKeyFilePath string) *rsa.PublicKey {
	privateKey := GetPrivateKey()

	return &privateKey.PublicKey
}
func CreateToken(guid string) (string, error){
	payload := jwt.MapClaims{
        "sub":  guid,
        // "exp":  timeConnect,
    }

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, payload)
	CreateKeys()
	privateKey := GetPrivateKey()
	t, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return t, nil
}


func ParseToken(tokenString string) (*jwt.MapClaims, error) {
	claims := &jwt.MapClaims{} 
	publicKey := GetPublicKeyByPrivateKey(privateKeyPath) 

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return nil, err 
	}

	if token.Valid {
		return claims, nil 
	} else {
		return nil, err
	}
}