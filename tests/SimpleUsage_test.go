package VkApiTest

import (
	"github.com/stels-cs/vk-api-tools"
	"os"
	"testing"
)

var envTokenTag = "VK_ACCESS_TOKEN"

func getToken() string {
	token := os.Getenv(envTokenTag)
	return token
}

func skip(t *testing.T) {
	t.Skip("No access_token passed, pass token in enviroment variable", envTokenTag, "to run this test")
}

func TestSimpleUsage1(t *testing.T) {
	if token := getToken(); token != "" {
		res, err := VkApi.Call("users.get", VkApi.P{"user_ids": "stelscs", "access_token": token})
		if err != nil {
			t.Error(err)
		}
		user, _ := res.FirstAny()
		if i, _ := user.GetInt("id"); i != 19039187 {
			t.Error("Expected 19039187, got", i)
		}
		if i, _ := user.GetString("first_name"); i != "Иван" {
			t.Error("Expected Иван, got", i)
		}
	} else {
		skip(t)
	}
}

func TestSimpleUsage2(t *testing.T) {
	if token := getToken(); token != "" {
		res, err := VkApi.Call("users.get", VkApi.P{"user_ids": "stelscs", "access_token": token})
		if err != nil {
			t.Error(err)
		}
		user, _ := res.FirstAny()
		if i, _ := user.GetInt("id"); i != 19039187 {
			t.Error("Expected 19039187, got", i)
		}
		if i, _ := user.GetString("first_name"); i != "Иван" {
			t.Error("Expected Иван, got", i)
		}
	} else {
		skip(t)
	}
}
