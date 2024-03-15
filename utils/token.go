package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// create RSA512 TOKEN

const privateKeyPath = "./key.pem"


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

