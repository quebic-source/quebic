package common

import (
	jwt "github.com/dgrijalva/jwt-go"
)

//CreateJWTToken create jwt token
func CreateJWTToken(claims jwt.MapClaims, secret string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil

}
