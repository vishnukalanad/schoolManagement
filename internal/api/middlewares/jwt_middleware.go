package middlewares

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"schoolManagement/pkg/utils"
)

func JWTMiddleware(next http.Handler) http.Handler {
	fmt.Println("-------------( JWT MIDDLEWARE STARTED )-------------")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-------------( INSIDE JWT MIDDLEWARE )-------------")
		token, err := r.Cookie("Bearer")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		jwtSecret := os.Getenv("JWT_SECRET")

		parsedToken, err := jwt.Parse(token.Value, func(token *jwt.Token) (interface{}, error) {
			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				http.Error(w, "Token expired!", http.StatusUnauthorized)
				return
			}

			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if parsedToken.Valid {
			log.Println("Valid JWT token")
		} else {
			http.Error(w, "Token invalid!", http.StatusUnauthorized)
			log.Println("Invalid JWT token")
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Token invalid!", http.StatusUnauthorized)
			return
		}

		// ContextKey is nothing but a custom type which actually is a string; this is to prevent compiler warning saying not to use string for keys in context;
		ctx := context.WithValue(r.Context(), utils.ContextKey("role"), claims["role"])
		ctx = context.WithValue(ctx, utils.ContextKey("expiry"), claims["exp"])
		ctx = context.WithValue(ctx, utils.ContextKey("username"), claims["user"])
		ctx = context.WithValue(ctx, utils.ContextKey("userid"), claims["uid"])

		next.ServeHTTP(w, r.WithContext(ctx))
		fmt.Println("-------------( SENT RESPONSE FROM JWT MIDDLEWARE )-------------")
	})
}
