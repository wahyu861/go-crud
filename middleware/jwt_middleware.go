package middleware

import (
	"fmt"
	"go-crud/config"
	"go-crud/models"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	jwtMiddleware "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

func UseJWT() echo.MiddlewareFunc {
	return jwtMiddleware.WithConfig(jwtMiddleware.Config{
		SigningKey: JWTSecret,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return jwt.MapClaims{}
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": fmt.Sprintf("Token tidak valid: %v", err),
			})
		},
	})
}

// Middleware tambahan — mengambil user dari token dan masukkan ke context
func AttachUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userToken, ok := c.Get("user").(*jwt.Token)
			if !ok || userToken == nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"message": "Token tidak valid atau user tidak ditemukan",
				})
			}

			claims, ok := userToken.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"message": "Klaim token tidak valid",
				})
			}

			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"message": "Klaim user_id tidak ditemukan",
				})
			}

			var user models.User
			if err := config.DB.First(&user, uint(userIDFloat)).Error; err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"message": "User tidak ditemukan di database",
				})
			}

			// simpan user ke context
			c.Set("authUser", user)
			return next(c)
		}
	}
}
