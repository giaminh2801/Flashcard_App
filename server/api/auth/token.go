package auth

import (
	"errors"
	"go-flashcard-api/api/models"
	"go-flashcard-api/api/responses"
	"go-flashcard-api/config"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

// GenerateJWT creates a new token to the client
func GenerateJWT(user models.User) (string, error) {
	claims := models.Claim{
		User: user,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "Minh Le",
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.SECRETKEY)
}

// ExtractToken retrieves the token from headers and Query
func ExtractToken(w http.ResponseWriter, r *http.Request) *jwt.Token {
	token, err := request.ParseFromRequestWithClaims(
		r,
		request.OAuth2Extractor,
		&models.Claim{},
		func(tkn *jwt.Token) (interface{}, error) {
			return config.SECRETKEY, nil
		},
	)
	if err != nil {
		errCode := http.StatusUnauthorized
		switch err.(type) {
		case *jwt.ValidationError:
			validationError := err.(*jwt.ValidationError)
			switch validationError.Errors {
			case jwt.ValidationErrorExpired:
				err = errors.New("Token has expired")
				responses.ERROR(w, errCode, err)
				return nil
			case jwt.ValidationErrorSignatureInvalid:
				err = errors.New("The signature is invalid")
				responses.ERROR(w, errCode, err)
				return nil
			default:
				responses.ERROR(w, errCode, err)
				return nil
			}
		}
	}

	return token
}
