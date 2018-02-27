package VkApi

import "fmt"

type Method struct {
	name   string
	params P
}

func (m *Method) toExecute() (string, error) {
	json, err := m.params.toJson()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("API.%s(%s)", m.name, string(json)), nil
}

func GetApiMethod(name string, params P) Method {
	return Method{name, params}
}

func CreateMethod(name string, params P) Method {
	return Method{name, P(params)}
}
