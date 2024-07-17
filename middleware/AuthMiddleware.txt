package middleware

import (
	"net/http"
	"net/url"
	"testhub-spec-uni/controllers"

	"github.com/beego/beego/v2/server/web/context"
	"github.com/dgrijalva/jwt-go"
)

func AuthMiddleware(ctx *context.Context) {
	tokenString := ctx.Input.Header("Authorization")
	if tokenString == "" {
		cookie, err := ctx.Request.Cookie("token")
		if err == nil {
			tokenString = cookie.Value
		}
	}

	if tokenString == "" {
		redirectToLoginWithNext(ctx)
		return
	}

	claims := &controllers.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your_secret_key"), nil
	})

	if err != nil || !token.Valid {
		redirectToLoginWithNext(ctx)
		return
	}

	if claims.Role != "Admin" {
		ctx.Output.SetStatus(http.StatusForbidden)
		ctx.Output.Body([]byte("Forbidden"))
		return
	}

	ctx.Input.SetData("username", claims.Username)
	ctx.Input.SetData("role", claims.Role)
}

func redirectToLoginWithNext(ctx *context.Context) {
	next := url.QueryEscape(ctx.Input.URI())
	ctx.Redirect(http.StatusFound, "/login?next="+next)
}
