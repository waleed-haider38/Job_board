package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// SAME claims struct jo login mein use hui
type JwtCustomClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		// 1️. Authorization header read karo
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Missing authorization header",
			})
		}

		// 2️. "Bearer token" format check
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Invalid authorization format",
			})
		}

		tokenString := parts[1]

		// 3️. Token parse + verify
		token, err := jwt.ParseWithClaims(
			tokenString,
			&JwtCustomClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte("secret"), nil
			},
		)

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Invalid or expired token",
			})
		}

		// 4️. Claims extract karo
		claims := token.Claims.(*JwtCustomClaims)

		// 5️. user_id context mein daalo
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		// 6️. next handler ko call
		return next(c)
	}
}
