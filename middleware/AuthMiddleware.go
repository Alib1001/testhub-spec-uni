package middleware

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/beego/beego/v2/server/web/context"
	"log"
	"net/http"
	"strings"
	"time"
)

type UserInfo struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	UUID      string `json:"uuid"`
	Role      string `json:"role"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Balance   int    `json:"balance"`
}

func AuthMiddleware(ctx *context.Context) {

	path := ctx.Input.URL()

	if path == "/api/cities" || path == "/api/subjects" {
		return
	}

	if ctx.Input.Method() == "OPTIONS" {
		// Разрешаем OPTIONS запросы без проверки токена
		ctx.Output.SetStatus(http.StatusOK)
		return
	}

	token := ctx.Input.Header("Authorization")
	if token == "" {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]string{"error": "Authorization header is missing"}, true, true)
		return
	}

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

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		ctx.Output.SetStatus(http.StatusInternalServerError)
		ctx.Output.JSON(map[string]string{"error": "Failed to decode response"}, true, true)
		return
	}

	// Проверка, является ли пользователь суперпользователем
	isSuperUser, err := IsSuperUser(userInfo.ID)
	if err != nil {
		ctx.Output.SetStatus(http.StatusInternalServerError)
		ctx.Output.JSON(map[string]string{"error": "Failed to check superuser status"}, true, true)
		return
	}

	if !isSuperUser {
		ctx.Output.SetStatus(http.StatusForbidden)
		ctx.Output.JSON(map[string]string{"error": "Access forbidden: only superusers allowed"}, true, true)
		return
	}
}

func IsSuperUser(userId int) (bool, error) {
	o := orm.NewOrm()
	var isSuperUser bool
	err := o.Raw("SELECT is_superuser FROM accounts.user WHERE id = ?", userId).QueryRow(&isSuperUser)
	if err != nil {
		log.Printf("Error querying user superuser status: %v", err)
		return false, err
	}
	return isSuperUser, nil
}
