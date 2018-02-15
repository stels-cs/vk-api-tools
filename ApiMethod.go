package Vk

import "fmt"

type ApiMethod struct {
	name Method
	params Params
}

func (m *ApiMethod) toExecute() (string, error) {
	json, err:= m.params.toJson()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("API.%s(%s)", m.name, string(json)), nil
}

func GetApiMethod(name Method, params Params) ApiMethod {
	return ApiMethod{name, params}
}