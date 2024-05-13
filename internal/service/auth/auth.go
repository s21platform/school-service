package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"log"
	"regexp"
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

type tokenResponse struct {
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

func LoginToPlatform(email, password string) (string, error) {
	state := uuid.New().String()
	nonce := uuid.New().String()

	// Setup client. Switch off redirects
	client := resty.New()
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	// Generate first link to sberclass
	link := fmt.Sprintf("https://auth.sberclass.ru/auth/realms/EduPowerKeycloak/protocol/openid-connect/auth?client_id=school21&redirect_uri=https%%3A%%2F%%2Fedu.21-school.ru%%2F&state=%s&response_mode=fragment&response_type=code&scope=openid&nonce=%s", state, nonce)

	// First request to get authURL
	res, err := client.R().Get(link)
	if err != nil || res.StatusCode() != 200 {
		log.Fatalln("Error while make request")
	}

	// Set New Cookie
	client.SetCookies(res.Cookies())

	// Get the authUrl from response
	loginActionPattern := regexp.MustCompile(`https://.+?"`)
	loginUrl := loginActionPattern.FindString(string(res.Body()))
	loginUrl = loginUrl[:len(loginUrl)-1]
	loginUrl = strings.Replace(loginUrl, "amp;", "", -1)
	fmt.Println("Login Action URL:", loginUrl)

	res, err = client.R().SetContext(context.Background()).SetFormData(map[string]string{
		"username": email,
		"password": password,
	}).Post(loginUrl)

	fmt.Println(res.StatusCode())

	if err != nil && res.StatusCode() != 302 {
		fmt.Println("not 302", err)
	}

	client.SetCookies(res.Cookies())

	location := res.Header().Get("location")
	res, err = client.R().Post(location)

	if err != nil && res.StatusCode() != 302 {
		log.Fatalf("Error")
	}

	location = res.Header().Get("location")
	oauthCodePattern := regexp.MustCompile(`code=([^&$]+)[&$]?`)

	oauthCodeMatch := oauthCodePattern.FindStringSubmatch(location)
	oauthCode := ""
	if len(oauthCodeMatch) > 1 {
		oauthCode = oauthCodeMatch[1]
	}

	client.SetCookies(res.Cookies())

	res, err = client.R().SetFormData(map[string]string{
		"code":         oauthCode,
		"grant_type":   "authorization_code",
		"client_id":    "school21",
		"redirect_uri": "https://edu.21-school.ru/",
	}).Post("https://auth.sberclass.ru/auth/realms/EduPowerKeycloak/protocol/openid-connect/token")

	tokStruct := tokenResponse{}
	err = json.Unmarshal(res.Body(), &tokStruct)
	return tokStruct.AccessToken, nil
}
