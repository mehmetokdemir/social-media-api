package jwt

import "github.com/dgrijalva/jwt-go"

// Token JWT
type Token struct {
	Username string
	UserId   uint
	jwt.StandardClaims
}
