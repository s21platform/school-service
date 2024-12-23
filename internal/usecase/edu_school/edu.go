package edu_school

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/s21platform/school-service/internal/model"
)

func LoginToPlatform(email, password string) (*model.TokenResponse, error) {
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
	var result model.TokenResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Printf("Не удалось преобразовать байты в структуру: %v", err)
		return nil, err
	}
	return &result, nil
}
