package main

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTData struct {
	jwt.StandardClaims
	ShieldooClaims map[string]string `json:"shieldoo"`
}

func GenerateJWTAccessToken(instance string) (string, error) {
	// prepare claims for token
	claims := JWTData{
		StandardClaims: jwt.StandardClaims{
			// set token lifetime in timestamp
			ExpiresAt: time.Now().Add(time.Duration(time.Minute * 10)).Unix(),
		},
		ShieldooClaims: map[string]string{
			"instance": instance,
		},
	}

	// generate a string using claims and HS256 algorithm
	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// sign the generated key using secretKey
	token, err := tokenString.SignedString([]byte(shieldooApiKey))

	return token, err
}
