package models

import (
	"github.com/dgrijalva/jwt-go"
)

// ATClaim struct
type Claim struct {
	User User `json:"user"`
	jwt.StandardClaims
}
