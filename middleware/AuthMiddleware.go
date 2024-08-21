package middleware

import (
	"crypto/tls"
	"fmt"
	"github.com/beego/beego/v2/server/web/context"
	"net/http"
	"strings"
	"time"
)

func AuthMiddleware(ctx *context.Context) {
	if ctx.Input.Method() == "OPTIONS" {
		// Разрешаем OPTIONS запросы без проверки токена
		ctx.Output.SetStatus(http.StatusOK)
		return
	}

	// Ваш существующий код для обработки авторизации...
	token := ctx.Input.Header("Authorization")
	if token == "" {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]string{"error": "Authorization header is missing"}, true, true)
		return
	}

	// Add "Bearer " prefix if it's not present
	if !strings.HasPrefix(token, "Bearer ") {
		token = "Bearer " + token
	}

	authURL := "https://api-dev.testhub.kz/accounts/api/v1/me"
	req, err := http.NewRequest("GET", authURL, nil)
	if err != nil {
		ctx.Output.SetStatus(http.StatusInternalServerError)
		ctx.Output.JSON(map[string]string{"error": "Failed to create request"}, true, true)
		return
	}

	req.Header.Set("Authorization", token)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 10, // Set a 10-second timeout
	}

	resp, err := client.Do(req)
	if err != nil {
		ctx.Output.SetStatus(http.StatusInternalServerError)
		ctx.Output.JSON(map[string]string{"error": fmt.Sprintf("Failed to perform request: %v", err)}, true, true)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]string{"error": "Unauthorized"}, true, true)
		return
	}
}
