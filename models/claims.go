package models

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID    int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`

	jwt.StandardClaims
}
