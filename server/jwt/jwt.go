package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var secretKey = "secret-key"

type UserClaims struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func CreateToken(user User, expirationTime *jwt.NumericDate) *jwt.Token {
	claims := UserClaims{
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token
}

func SignToken(token *jwt.Token) (string, error) {
	signedString, err := token.SignedString(secretKey)

	return signedString, err
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil

}

func GetExparationTime() *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
}
