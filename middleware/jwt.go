package middleware

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/maxik12233/blog/common"
)

func GenerateJWT(ID uint, roles []uint) (string, error) {

	// Generate jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   ID,
		"roles": roles,
		"exp":   time.Now().Add(time.Hour * time.Duration(common.JWT_TOKEN_EXP_HOURS)).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return "", nil
	}

	return tokenString, nil
}
