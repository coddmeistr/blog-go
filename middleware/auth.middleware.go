package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/maxik12233/blog/common"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var tokenString string
		cookies := r.Cookies()
		for _, val := range cookies {
			if val.Name == common.JWT_TOKEN_NAME {
				tokenString = val.Value
			}
		}

		if tokenString == "0" || tokenString == "" {
			common.ReturnAnauthorized(w)
			return
		}

		// Decode/validate it
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			common.ReturnAnauthorized(w)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check the exp
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				common.ReturnAnauthorized(w)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				"UserID",
				claims["sub"],
			)
			ctx = context.WithValue(
				ctx,
				"Roles",
				claims["roles"],
			)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		} else {
			common.ReturnAnauthorized(w)
		}
	})
}

func ValidateRolesMiddleware(validRoles []uint) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			rolesVal := r.Context().Value("Roles").([]interface{})
			if rolesVal == nil {
				common.ReturnAnauthorized(w)
				return
			}

			roles := make([]uint, len(rolesVal))
			for i, val := range rolesVal {
				roles[i] = uint(val.(float64))
			}
			if len(roles) == 0 {
				common.ReturnAnauthorized(w)
				return
			}

			validation := make(map[uint]uint)
			for _, val := range roles {
				validation[val] = 0
			}

			for _, val := range validRoles {
				if _, ok := validation[val]; !ok {
					common.ReturnAnauthorized(w)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
