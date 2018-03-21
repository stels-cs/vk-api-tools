package VkApi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type AccessToken struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
	UserId    int    `json:"user_id"`
}

/*
	scope := VkApi.U_OFFLINE | VkApi.U_GROUPS | VkApi.U_WALL
	clientId := "3140623"
	clientSecret := "VeWdmVclDCtn6ihuP1nt"
	v := "5.73"
	login := "example@vk.com"
	password := "test-password-123"
	accessToken, err := VkApi.PasswordAuth(login, password, scope, clientId, clientSecret, v, nil)
	// accessToken.Token - string
*/
func PasswordAuth(login, password string, scope int, clientId, clientSecret, v string, h *http.Client) (AccessToken, error) {
	token := AccessToken{}
	if h == nil {
		h = &http.Client{Timeout: 30 * time.Second}
	}
	resp, err := h.PostForm(
		"https://oauth.vk.com/token",
		url.Values{
			"grant_type":    {"password"},
			"client_id":     {clientId},
			"client_secret": {clientSecret},
			"username":      {login},
			"password":      {password},
			"scope":         {strconv.Itoa(scope)},
			"v":             {v},
			"2fa_supported": {"1"},
		})
	if err != nil {
		return token, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return token, err
	}

	if err := json.Unmarshal(body, &token); err != nil {
		return token, err
	}
	if token.Token == "" {
		apiError := new(ApiError)
		if err := json.Unmarshal(body, &apiError); err != nil {
			return token, err
		} else {
			return token, apiError
		}
	}
	return token, nil
}
