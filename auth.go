package main

import (
	"errors"
	"log"
	"os"
	"time"

	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/freitagsrunde/modulist/db"
	"github.com/gin-gonic/gin"
)

// Functions

// CreateSession produces a JSON Web Token (JWT) with
// authenticated claims based on the values in the
// supplied user object and saves it inside the
// cookie storage.
func (app *App) CreateSession(c *gin.Context, User db.User) {

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

	// TODO: Set 'secure' to true.
	c.SetCookie("Token", sessionJWTString, int(app.JWTValidFor.Seconds()), "", "", false, true)
}

// Authorize takes a supplied request, extracts the to-
// be-included JWT out of the set cookies and validates
// it on various aspects. Additionally, it is checked that
// the user has at least the required minimum privilege
// specified by the calling handler.
func (app *App) Authorize(Request *http.Request, MinimumPrivilege int) (*db.User, error) {

	// Extract cookie with token from request.
	cookie, err := Request.Cookie("Token")
	if err != nil {
		return nil, errors.New("Authorization not present or correct. Please log in.")
	}

	// Parse authorization token.
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {

		// Verify that JWT was signed with correct algorithm.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("[Authorize] Unexpected signing method: %v.", token.Header["alg"])

			return nil, errors.New("Authorization not present or correct. Please log in.")
		}

		// Return our JWT signing secret as the key to verify integrity of JWT.
		return []byte(os.Getenv("APP_JWT_SIGNING_SECRET")), nil
	})

	// Check for parsing errors.
	if err != nil {
		return nil, errors.New("Authorization not present or correct. Please log in.")
	}

	// Check if JWT is valid.
	if token.Valid != true {
		return nil, errors.New("Authorization not present or correct. Please log in.")
	}

	claims := token.Claims.(jwt.MapClaims)
	if err := claims.Valid(); err != nil {

		// Claims in JWT were not valid.
		log.Printf("[Authorize] Claims in JWT were not correct: %s\n", err.Error())

		return nil, errors.New("Authorization not present or correct. Please log in.")
	}

	// Extract user's mail out of claims in JWT.
	userMail, ok := claims["iss"].(string)
	if !ok {
		return nil, errors.New("Authorization not present or correct. Please log in.")
	}

	var User db.User
	app.DB.First(&User, "\"mail\" = ?", userMail)

	// Check if logged-in user is allowed to view page.
	if User.Privileges > MinimumPrivilege {
		return nil, errors.New("You do not have sufficient privileges.")
	}

	// We found the logged-in user.
	return &User, nil
}
