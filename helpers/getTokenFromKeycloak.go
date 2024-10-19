package helpers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func GetTokenFromKeycloak(username string, password string) (string, error) {
	tokenUrl := os.Getenv("TOKEN_URL")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "password")
	data.Set("username", username)
	data.Set("password", password)
	fmt.Println(data)
	resp, err := http.Post(tokenUrl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	if err != nil {
		return "", err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			return
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
