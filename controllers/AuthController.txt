package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testhub-spec-uni/models"
	"time"

	"github.com/beego/beego/v2/server/web"
	"github.com/dgrijalva/jwt-go"
)

type AuthController struct {
	web.Controller
}

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// PostLogin performs user login and issues a JWT token.
func (c *AuthController) PostLogin() {
	username := c.GetString("username")
	password := c.GetString("password")

	if username == "" || password == "" {
		requestBody := string(c.Ctx.Input.CopyBody(1024))
		var user models.User
		err := json.Unmarshal([]byte(requestBody), &user)
		if err != nil || user.Username == "" || user.Password == "" {
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			c.Ctx.Output.Body([]byte("Invalid request"))
			fmt.Println("Error decoding JSON:", err)
			return
		}
		username = user.Username
		password = user.Password
	}

	dbUser, err := models.GetUserByUsername(username)
	if err != nil || !models.CheckPasswordHash(password, dbUser.Password) {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.Body([]byte("Unauthorized"))
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: dbUser.Username,
		Role:     dbUser.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.Body([]byte("Internal Server Error"))
		return
	}

	// Check if the request is from Postman or a similar tool
	contentType := c.Ctx.Input.Header("Content-Type")
	isPostman := strings.Contains(contentType, "application/json")

	if isPostman {
		// If it's a Postman request, return the token in the response body
		c.Ctx.Output.Body([]byte("Token: " + tokenString))
	} else {
		// If it's a browser request, set the cookie and redirect
		c.Ctx.SetCookie("token", tokenString, 24*60*60, "/")
		next := c.GetString("next")
		if next != "" {
			c.Redirect(next, http.StatusFound)
		} else {
			c.Redirect("/", http.StatusFound)
		}
	}
}

// GetLogin renders the login page.
func (c *AuthController) GetLogin() {
	next := c.GetString("next")
	c.Data["next"] = next
	c.TplName = "auth/login.tpl"
	c.Render()
}

// VerifyToken verifies the authenticity of a JWT token.
func (c *AuthController) VerifyToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, fmt.Errorf("invalid token signature")
		}
		return nil, fmt.Errorf("invalid token")
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

// AuthMiddleware is a middleware function to authenticate requests using JWT token.
func (c *AuthController) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		claims, err := c.VerifyToken(tokenStr.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "username", claims.Username)
		ctx = context.WithValue(ctx, "role", claims.Role)
		next(w, r.WithContext(ctx))
	}
}
