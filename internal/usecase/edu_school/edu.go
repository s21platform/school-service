package edu_school

import (
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type TokenResponse struct {
	Error            string `json:"error"`
	AccessToken      string `json:"access_token"`
	ExpiresIn        int64  `json:"expires_in"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	IDToken          string `json:"id_token"`
	NotBeforePolicy  int64  `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

func LoginToPlatform(email, password string) (*TokenResponse, error) {
	// URL для запроса токена
	tokenURL := "https://auth.sberclass.ru/auth/realms/EduPowerKeycloak/protocol/openid-connect/token"

	// Параметры запроса
	data := url.Values{}
	data.Set("username", email)
	data.Set("password", password)
	data.Set("grant_type", "password")
	data.Set("client_id", "s21-open-api")

	// Создание HTTP-клиента и отправка POST-запроса
	client := &http.Client{}
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return nil, err
	}

	// Установка заголовков для запроса
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Выполнение запроса
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Проверка кода состояния ответа
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Не удалось получить токен. Код состояния:", resp.StatusCode)
		return nil, status.Errorf(codes.InvalidArgument, "Неверно указан логин или пароль")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Не удалось преобразовать ответ в байты: %v", err)
		return nil, err
	}

	// Чтение и вывод ответа
	var result TokenResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Printf("Не удалось преобразовать байты в структуру: %v", err)
		return nil, err
	}
	return &result, nil
}
