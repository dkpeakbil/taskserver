package domain

import "github.com/dgrijalva/jwt-go"

var TokenSecret = []byte("my_token_secret")

type TokenClaims struct {
	jwt.StandardClaims
	UserID int
}
