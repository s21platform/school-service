package auth

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoginToPlatform(t *testing.T) {
	// TestLoginToPlatform is a test function for login to platform
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("t", r.RequestURI, r.Method)
		switch {
		case strings.Contains(r.RequestURI, "/auth/realms/EduPowerKeycloak/protocol/openid-connect/auth"):
			// First request to get authURL
			println("here")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`"https://auth.sberclass.ru/mockloginurl"`))
		case strings.Contains(r.RequestURI, "mockloginurl"):
			// Second request for login
			fmt.Println("here 2")
			http.Redirect(w, r, "https://mocklocation", http.StatusFound)
		case strings.Contains(r.RequestURI, "mocklocation"):
			// Third request to get oauth code
			http.Redirect(w, r, "https://mocklocation?code=mockcode", http.StatusFound)
		case strings.Contains(r.RequestURI, "openid-connect/token"):
			// Fourth request to get token
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token": "mocktoken"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()
	fmt.Println(mockServer.URL)

	client := resty.New()
	client.SetHostURL(mockServer.URL)

	// Call the function with the mocked client
	token, err := LoginToPlatform(client, mockServer.URL, "testemail", "testpassword")

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, "mocktoken", token)
}
