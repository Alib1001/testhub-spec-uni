package middleware

import (
	"crypto/tls"
	"fmt"
	"github.com/beego/beego/v2/server/web/context"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func AuthMiddleware(ctx *context.Context) {
	// Извлечение заголовка Authorization
	token := ctx.Input.Header("Authorization")
	if token == "" {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]string{"error": "Authorization header is missing"}, true, true)
		return
	}

	// Если используется Bearer Token, проверяем наличие префикса "Bearer "
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

	// Добавляем заголовок Authorization в запрос
	req.Header.Set("Authorization", token)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 10, // Устанавливаем тайм-аут 10 секунд
	}

	fmt.Println("Making request to:", authURL)
	fmt.Println("Authorization Header:", token)

	resp, err := client.Do(req)
	if err != nil {
		ctx.Output.SetStatus(http.StatusInternalServerError)
		ctx.Output.JSON(map[string]string{"error": fmt.Sprintf("Failed to perform request: %v", err)}, true, true)
		fmt.Println("Request Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.Output.SetStatus(http.StatusInternalServerError)
		ctx.Output.JSON(map[string]string{"error": "Failed to read response body"}, true, true)
		return
	}

	if resp.StatusCode != http.StatusOK {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]string{"error": "Unauthorized"}, true, true)
		return
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
	fmt.Println("Authenticated user:", string(body))
}
