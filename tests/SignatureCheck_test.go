package VkApiTest

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	VkApi "vk-api-tools"
)

func TestSignatureCalculation(t *testing.T) {

	appSecret := "UURSsxO59uTyHVvSzHgW"
	request := "/?api_url=https://api.vk.com/api.php&api_id=6196804&api_settings=1&viewer_id=19039187&viewer_type=0&sid=e211a8bf9bad808a2a95d75721071b874ba82d07a8b0b6aaeb98f2d220deca8fd591c89a2dca1c6165b8e&secret=9c3f105f93&access_token=064affc04d119ad5798e9e8e2b24012fcad249be99712151047532d53f2dd107f24195f6d7309bceb0274&user_id=0&is_app_user=1&auth_key=7eb1471c6341ba56ff0c0dad0f8dba6b&language=0&parent_language=0&is_secure=1&ads_app_id=6196804_e7d36e80a3155f8eb0&referrer=unknown&lc_name=abe9e425&sign=17b0427e7a43f60d081487c36170ff6d052516d06341457668391a22fd7732c1&hash="

	if VkApi.IsCorrectRequest(request, appSecret) == false {
		t.Error("Expected correct request, got incorrent")
	}

}

func TestMiniAppsSignature(t *testing.T) {

	secret := "rkwdOT04kUh28RDEC9zr"
	request := "vk_access_token_settings=friends%2Cgroups&vk_app_id=6825462&vk_are_notifications_enabled=0&vk_is_app_user=1&vk_language=ru&vk_platform=desktop_web&vk_user_id=19039187&sign=vBBPIysvzccFUn_e55JCGxZBnmxpXeh92XpiAY9gcv8"
	sign := "vBBPIysvzccFUn_e55JCGxZBnmxpXeh92XpiAY9gcv8"

	v, err := url.ParseQuery(request)
	assert.Nil(t, err)
	assert.Equal(t, sign, VkApi.CalculateMiniAppsSignature(v, secret))

	assert.Equal(t, true, VkApi.IsCorrectMiniAppsRequest(request, secret))

}
