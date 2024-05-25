package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
	Role     Role      `json:"role"`
	Exp      time.Time `json:"exp"`
	jwt.StandardClaims
}
