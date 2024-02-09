package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vhladko/books/models"
)

var secretKey = []byte("secret-key")

func CreateToken(user models.User, expirationTime *jwt.NumericDate) *jwt.Token {
	claims := models.UserClaims{
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

func ParseToken(token string) *models.UserClaims {
	parsedAccessToken, _ := jwt.ParseWithClaims(token, &models.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	return parsedAccessToken.Claims.(*models.UserClaims)
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
