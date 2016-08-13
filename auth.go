package main

import (
	"log"
	"os"
	"time"

	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/numbleroot/modulist/db"
)

// Functions

// CreateJWT produces a JSON Web Token (JWT) with
// authenticated claims based on the values in the
// supplied user object.
func (app *App) CreateJWT(User *db.User) string {

	// Retrieve the session signing key from environment.
	jwtSigningSecret := os.Getenv("APP_JWT_SIGNING_SECRET")

	// Save current timestamp.
	nowTime := time.Now()
	expTime := nowTime.Add(app.JWTValidFor)

	// Create a JWT with claims to identify user.
	sessionJWT := jwt.New(jwt.SigningMethodHS512)
	claims := sessionJWT.Claims.(jwt.MapClaims)

	// Add these claims.
	claims["iss"] = User.Mail
	claims["iat"] = nowTime.Unix()
	claims["nbf"] = nowTime.Add((-1 * time.Minute)).Unix()
	claims["exp"] = expTime.Unix()

	sessionJWTString, err := sessionJWT.SignedString([]byte(jwtSigningSecret))
	if err != nil {
		log.Fatalf("[CreateJWT] Creating JWT went wrong: %s.\nTerminating.", err)
	}

	return sessionJWTString
}

// Authorize takes a supplied request, extracts the to-
// be-included JWT out of the set cookies and validates
// it on various aspects.
func (app *App) Authorize(Request *http.Request) (bool, *db.User, string) {

	// TODO: Implement this.
	return true, nil, ""
}
