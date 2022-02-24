package jwt

import (
	"errors"
	"flow-users/user"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtCustumClaims struct {
	Id    uint64 `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateToken(user user.UserPostResponse, issuer string, secret string) (token string, err error) {
	// Set custom claims
	claims := &JwtCustumClaims{
		user.Id,
		user.Email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
		},
	}

	// Generate token
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return newToken.SignedString([]byte(secret))
}

func CheckToken(issuer string, token *jwt.Token) (id uint64, err error) {
	claims := token.Claims.(*JwtCustumClaims)

	if !claims.VerifyIssuer(issuer, true) {
		// Invalid token
		return 0, errors.New("invalid token")
	}

	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		// Token expired
		return 0, errors.New("token expired")
	}

	return claims.Id, nil
}
