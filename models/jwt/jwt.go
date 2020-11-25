package jwt

import (
	"fmt"
	"time"

	"github.com/DLTcollab/vehicle-data-explorer/models/config"
	"github.com/dgrijalva/jwt-go"
)

// custom claims
type Claims struct {
	Hash string `json:"hash"`
	jwt.StandardClaims
}

func CreateJwtToken(hash string) (string, error) {

	claims := Claims{
		Hash: hash,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(config.ConfigHandler.GetString("SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func VerifyJwtToken(tokenString string) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.ConfigHandler.GetString("SECRET")), nil
	})
	return token, err
}
