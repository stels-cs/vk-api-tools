package Vk

import (
	"io/ioutil"
	"encoding/json"
	"net/http"
	"net/url"
)

type AccessToken struct {
	Token string `json:"access_token"`
	ExpiresIn int `json:"expires_in"`
	UserId int `json:"user_id"`
}

type AuthError struct{
	Message string `json:"error"`
	Description string `json:"error_description"`
}

func (e *AuthError) Error() string {
	return e.Message + ": " + e.Description
}

func PasswordAuth(login string, password string) (AccessToken, error) {
	token := AccessToken{}
	resp, err := http.PostForm(
		"https://oauth.vk.com/token",
		url.Values{
			"grant_type":    {"password"},
			"client_id":     {"3140623"},
			"client_secret": {"VeWdmVclDCtn6ihuP1nt"},
			"username":      {login},
			"password":      {password},
			"scope":         {"12288"},
			"v":             {"5.69"},
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
		apiError := new(AuthError)
		if err := json.Unmarshal(body, &apiError); err != nil {
			return token, err
		} else {
			return token, apiError
		}
	}
	return token, nil
}