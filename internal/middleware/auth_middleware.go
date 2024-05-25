package middleware

import (
	"fmt"
	"smart_electricity_tracker_backend/internal/config"
	"smart_electricity_tracker_backend/internal/helpers"
	"smart_electricity_tracker_backend/internal/models"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type AuthMiddlewareService struct {
	cfg *config.Config
}

func NewAuthMiddleware(cfg *config.Config) *AuthMiddlewareService {
	return &AuthMiddlewareService{cfg: cfg}
}

func (s *AuthMiddlewareService) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := &models.Claims{}
		tokenString := c.Get("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.cfg.JWTSecret), nil
		})
		if err != nil {
			log.Infof("Error: %v", err)
			return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
		}
		if !token.Valid {
			log.Infof("Error: %v", err)
			return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
		}

		if claims.Exp.Before(time.Now()) {
			return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "Token expired")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func (s *AuthMiddlewareService) Permission(roleApprover []models.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)

		for _, roleApp := range roleApprover {
			if role == fmt.Sprint(roleApp) {
				return c.Next()
			}
		}

		return helpers.ErrorResponse(c, fiber.StatusForbidden, "Forbidden")
	}
}
