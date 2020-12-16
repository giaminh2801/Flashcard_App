package types

import "github.com/dgrijalva/jwt-go"

// StringKey is the context key for user model
type StringKey string

// RedirectValues for context
type RedirectValues struct {
	Path         string
	RefreshToken *jwt.Token
}
