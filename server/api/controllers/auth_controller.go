package controllers

import (
	"context"
	"encoding/json"
	"go-flashcard-api/api/auth"
	"go-flashcard-api/api/models"
	"go-flashcard-api/api/responses"
	"go-flashcard-api/api/utils/types"
	"go-flashcard-api/config"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Login handler
func Login(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user, tokenDetails, err := auth.SignIn(user.Email, user.Password)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenDetails.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(tokenDetails.ATExpires, 0),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenDetails.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(tokenDetails.RTExpires, 0),
	})

	responses.JSON(w, http.StatusOK, user)
}

// Logout handler
func Logout(w http.ResponseWriter, r *http.Request) {
	// Delete access_token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	// Delete refresh_token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	responses.JSON(w, http.StatusOK, "logged_out")
}

// Refresh token handler
func Refresh(ch chan<- bool, pErr *error, pCtx *context.Context, refreshToken *jwt.Token, w http.ResponseWriter, r *http.Request) {
	ATExpires := time.Now().Add(time.Minute * 5).Unix()
	RTExpires := time.Now().Add(time.Hour * 24 * 7).Unix()

	user := refreshToken.Claims.(*models.Claim).User
	ATClaims := &models.Claim{
		User: user,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "Minh Le",
			ExpiresAt: ATExpires,
		},
	}
	RTClaims := &models.Claim{
		User: user,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "Minh Le",
			ExpiresAt: RTExpires,
		},
	}
	// Creating Access Token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, ATClaims)
	accessTokenStr, err := accessToken.SignedString(config.ATSECRET)
	if err != nil {
		*pErr = err
		ch <- false
		return
	}
	//Creating Refresh Token
	refreshToken = jwt.NewWithClaims(jwt.SigningMethodHS256, RTClaims)
	refreshTokenStr, err := refreshToken.SignedString(config.RTSECRET)
	if err != nil {
		*pErr = err
		ch <- false
		return
	}

	// Delete old refresh_token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Create new cookies for access- and refresh tokens
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Path:     "/",
		Value:    accessTokenStr,
		HttpOnly: true,
		Expires:  time.Unix(ATExpires, 0),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Path:     "/",
		Value:    refreshTokenStr,
		HttpOnly: true,
		Expires:  time.Unix(RTExpires, 0),
	})
	*pCtx = context.WithValue(
		r.Context(),
		types.StringKey("user"),
		accessToken.Claims.(*models.Claim).User,
	)
	ch <- true
}
