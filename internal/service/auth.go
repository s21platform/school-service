package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//const (
//	kсBaseUrl         = "https://auth.sberclass.ru/auth/realms/EduPowerKeycloak"
//	cookieUrlTemplate = kсBaseUrl + "/protocol/openid-connect/auth?client_id=school21&redirect_uri=https%%3A%%2F%%2Fedu.21-school.ru%%2F&state=%s&response_mode=fragment&response_type=code&scope=openid&nonce=%s"
//	tokenUrl          = kсBaseUrl + "/protocol/openid-connect/token"
//)

//var (
//	loginActionPattern = regexp.MustCompile(`(?P<LoginActionURL>https:\/\/.+?)"`)
//	oauthCodePattern   = regexp.MustCompile(`code=(?P<OAuthCode>.+)[&$]?`)
//)
//
//func getLoginActionUrl(data []byte) string {
//	rawUrl := loginActionPattern.FindStringSubmatch(string(data))[loginActionPattern.SubexpIndex("LoginActionURL")]
//
//	return strings.ReplaceAll(rawUrl, "amp;", "")
//}

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

//func LoginToPlatform(client *resty.Client, email, password string) (string, error) {
//	state := uuid.New().String()
//	nonce := uuid.New().String()
//
//	// Generate first link to sberclass
//	link := fmt.Sprintf("https://auth.sberclass.ru/auth/realms/EduPowerKeycloak/protocol/openid-connect/auth?client_id=school21&redirect_uri=https%%3A%%2F%%2Fedu.21-school.ru%%2F&state=%s&response_mode=fragment&response_type=code&scope=openid&nonce=%s", state, nonce)
//
//	// First request to get authURL
//	res, err := client.R().Get(link)
//	if err != nil || res.StatusCode() != 200 {
//		return "", status.New(400, "Error while get auth link").Err()
//	}
//
//	// Set New Cookie
//	client.SetCookies(res.Cookies())
//
//	// Get the authUrl from response
//	loginActionPattern := regexp.MustCompile(`https://.+?"`)
//	loginUrl := loginActionPattern.FindString(string(res.Body()))
//	loginUrl = loginUrl[:len(loginUrl)-1]
//	loginUrl = strings.Replace(loginUrl, "amp;", "", -1)
//	fmt.Println("Login Action URL:", loginUrl)
//
//	res, err = client.R().SetContext(context.Background()).SetFormData(map[string]string{
//		"username": email,
//		"password": password,
//	}).Post(loginUrl)
//
//	fmt.Println(res.StatusCode())
//
//	if err != nil && res.StatusCode() != 302 {
//		return "", status.New(400, "Error while send credential").Err()
//	}
//
//	client.SetCookies(res.Cookies())
//
//	location := res.Header().Get("location")
//	res, err = client.R().Post(location)
//
//	if err != nil && res.StatusCode() != 302 {
//		return "", status.New(401, "Login or password is incorrect").Err()
//	}
//
//	location = res.Header().Get("location")
//	oauthCodePattern := regexp.MustCompile(`code=([^&$]+)[&$]?`)
//
//	oauthCodeMatch := oauthCodePattern.FindStringSubmatch(location)
//	oauthCode := ""
//	if len(oauthCodeMatch) > 1 {
//		oauthCode = oauthCodeMatch[1]
//	}
//
//	client.SetCookies(res.Cookies())
//
//	res, err = client.R().SetFormData(map[string]string{
//		"code":         oauthCode,
//		"grant_type":   "authorization_code",
//		"client_id":    "school21",
//		"redirect_uri": "https://edu.21-school.ru/",
//	}).Post("https://auth.sberclass.ru/auth/realms/EduPowerKeycloak/protocol/openid-connect/token")
//
//	tokStruct := tokenResponse{}
//	err = json.Unmarshal(res.Body(), &tokStruct)
//	return tokStruct.AccessToken, nil
//}

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
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
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

//func LoginToPlatform(client *resty.Client, baseURL, email, password string) (string, error) {
//	state := uuid.New().String()
//	nonce := uuid.New().String()
//
//	// Generate first link to sberclass
//	link := fmt.Sprintf("%s/auth/realms/EduPowerKeycloak/protocol/openid-connect/auth?client_id=school21&redirect_uri=https%%3A%%2F%%2Fedu.21-school.ru%%2F&state=%s&response_mode=fragment&response_type=code&scope=openid&nonce=%s", baseURL, state, nonce)
//
//	// First request to get authURL
//	res, err := client.R().Get(link)
//	if err != nil || res.StatusCode() != 200 {
//		return "", status.New(400, "Error while get auth link").Err()
//	}
//
//	// Set New Cookie
//	client.SetCookies(res.Cookies())
//
//	// Get the authUrl from response
//	loginActionPattern := regexp.MustCompile(`https://.+?"`)
//	loginUrl := loginActionPattern.FindString(string(res.Body()))
//	loginUrl = loginUrl[:len(loginUrl)-1]
//	loginUrl = strings.Replace(loginUrl, "amp;", "", -1)
//	fmt.Println("Login Action URL:", loginUrl)
//
//	res, err = client.R().SetFormData(map[string]string{
//		"username": email,
//		"password": password,
//	}).Post(loginUrl)
//
//	fmt.Println("statusCode", res.StatusCode())
//
//	if err != nil && res.StatusCode() != 302 {
//		return "", status.New(400, "Error while send credential").Err()
//	}
//
//	client.SetCookies(res.Cookies())
//
//	location := res.Header().Get("location")
//	res, err = client.R().Post(location)
//
//	if err != nil && res.StatusCode() != 302 {
//		return "", status.New(401, "Login or password is incorrect").Err()
//	}
//
//	location = res.Header().Get("location")
//	oauthCodePattern := regexp.MustCompile(`code=([^&$]+)[&$]?`)
//
//	oauthCodeMatch := oauthCodePattern.FindStringSubmatch(location)
//	oauthCode := ""
//	if len(oauthCodeMatch) > 1 {
//		oauthCode = oauthCodeMatch[1]
//	}
//
//	client.SetCookies(res.Cookies())
//
//	res, err = client.R().SetFormData(map[string]string{
//		"code":         oauthCode,
//		"grant_type":   "authorization_code",
//		"client_id":    "school21",
//		"redirect_uri": "https://edu.21-school.ru/",
//	}).Post(fmt.Sprintf("%s/auth/realms/EduPowerKeycloak/protocol/openid-connect/token", baseURL))
//
//	tokStruct := tokenResponse{}
//	err = json.Unmarshal(res.Body(), &tokStruct)
//	return tokStruct.AccessToken, nil
//}
