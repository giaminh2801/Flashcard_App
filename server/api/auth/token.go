package auth

import (
	"fmt"
	"go-flashcard-api/api/models"
	"go-flashcard-api/api/responses"
	"go-flashcard-api/config"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenDetails struct
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	ATExpires    int64
	RTExpires    int64
}

// GenerateJWT creates a new token to the client
func GenerateJWT(user models.User) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.ATExpires = time.Now().Add(time.Minute * 5).Unix()
	td.RTExpires = time.Now().Add(time.Hour * 24 * 7).Unix()

	ATClaims := models.Claim{
		User: user,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "Minh Le",
			ExpiresAt: td.ATExpires,
		},
	}
	RTClaims := models.Claim{
		User: user,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "Minh Le",
			ExpiresAt: td.RTExpires,
		},
	}
	var err error
	// Creating Access Token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, ATClaims)
	td.AccessToken, err = accessToken.SignedString(config.ATSECRET)
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, RTClaims)
	td.RefreshToken, err = refreshToken.SignedString(config.RTSECRET)
	if err != nil {
		return nil, err
	}

	return td, nil
}

// VerifyAccessToken verifies access token and returns that token if valid
func VerifyAccessToken(w http.ResponseWriter, r *http.Request) *jwt.Token {
	accessTokenCookie, err := r.Cookie("access_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return nil
		}
		return nil
	}
	accessTokenStr := accessTokenCookie.Value

	accessToken, err := jwt.ParseWithClaims(accessTokenStr, &models.Claim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected access token's signing method: %v", token.Header["alg"])
		}
		return config.ATSECRET, nil
	})
	if err != nil {
		return nil
	}
	return accessToken
}

// VerifyRefreshToken verifies refresh token and returns that
func VerifyRefreshToken(w http.ResponseWriter, r *http.Request) *jwt.Token {
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		if err == http.ErrNoCookie {
			responses.ERROR(w, http.StatusUnauthorized, err)
			return nil
		}
		responses.ERROR(w, http.StatusBadRequest, err)
		return nil
	}

	refreshTokenStr := refreshTokenCookie.Value
	refreshToken, err := jwt.ParseWithClaims(refreshTokenStr, &models.Claim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected access token's signing method: %v", token.Header["alg"])
		}
		return config.RTSECRET, nil
	})
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return nil
	}

	return refreshToken
}
