package VkApi

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

// Тип для работы с JSON неизвестной структуры
type AnyModel json.RawMessage

func keyNotExsits(key string) string {
	return "Key " + key + " not exist"
}

func (m *AnyModel) Int() (int, error) {
	var i int
	err := json.Unmarshal([]byte(*m), &i)
	if err != nil {
		if err.Error() == "json: cannot unmarshal string into Go value of type int" {
			str, err := m.String()
			if err != nil {
				return 0, err
			} else {
				i, err := strconv.Atoi(str)
				if err != nil {
					return 0, err
				} else {
					return i, nil
				}
			}
		} else {
			return 0, err
		}
	} else {
		return i, nil
	}
}

// Кастует текущий JSON в Bool
//  any := VkApi.AnyModel("true")
//  b, err := any.Bool() // b == true err == nil
//  any = VkApi.AnyModel("krokodail")
//  b, err = any.Bool() // b == false err == error
func (m *AnyModel) Bool() (bool, error) {
	var b bool
	err := json.Unmarshal([]byte(*m), &b)
	if err != nil {
		return false, err
	} else {
		return b, nil
	}
}

// Кастует текущий JSON в Float
//  any := VkApi.AnyModel("5.666")
//  f, err := any.Float() // f == 5.666, err == nil
func (m *AnyModel) Float() (float64, error) {
	var f float64
	err := json.Unmarshal([]byte(*m), &f)
	if err != nil {
		return 0, err
	} else {
		return f, nil
	}
}

// Кастует текущий JSON в String
//  any := VkApi.AnyModel("\"Hello world\"")
//  s, err := any.Float() // s == "Hello world", err == nil
//  any = VkApi.AnyModel("1000")
//  s, err = any.Float() // s == "1000", err == nil
func (m *AnyModel) String() (string, error) {
	var i string
	err := json.Unmarshal([]byte(*m), &i)
	if err != nil {
		if err.Error() == "json: cannot unmarshal number into Go value of type string" {
			i, err := m.Int()
			if err != nil {
				return "", err
			} else {
				return strconv.Itoa(i), nil
			}
		} else {
			return "", err
		}
	} else {
		return i, nil
	}
}

// Если в текущем объекте лежит JSON объект то врнет новый объект со значением ключа
// any := VkApi.AnyModel(`{"name":"Ivan","age":15}`)
// name, err := any.GetAny("name")
// s, _ := name.String() // s == "Ivan"
func (m *AnyModel) GetAny(key string) (*AnyModel, error) {
	x := make(map[string]json.RawMessage, 0)
	err := json.Unmarshal([]byte(*m), &x)
	if err != nil {
		return nil, err
	}
	data, has := x[key]
	if has {
		am := AnyModel(data)
		return &am, nil
	} else {
		return nil, errors.New(keyNotExsits(key))
	}
}

// Если в текущем объекте лежит JSON объект то врнет строкове значение ключа
//  any := VkApi.AnyModel(`{"name":"Ivan","age":15}`)
//  name, err := any.GetString("name") // name == "Ivan" err == nil
func (m *AnyModel) GetString(key string) (string, error) {
	x := make(map[string]json.RawMessage, 0)
	err := json.Unmarshal([]byte(*m), &x)
	if err != nil {
		return "", err
	}
	data, has := x[key]
	if has {
		am := AnyModel(data)
		return am.String()
	} else {
		return "", errors.New(keyNotExsits(key))
	}
}

// Если в текущем объекте лежит JSON объект то врнет строкове значение ключа, если ключа нет, то вернет def значение
//  any := VkApi.AnyModel(`{"name":"Ivan","age":15}`)
//  name, err := any.GetStringDef("last_name", "Popovich") // name == "Popovich" err == nil
func (m *AnyModel) GetStringDef(key string, def string) (string, error) {
	res, err := m.GetString(key)
	if err == nil {
		return res, nil
	} else if err.Error() == keyNotExsits(key) {
		return def, nil
	} else {
		return def, err
	}
}

// Если в текущем объекте лежит JSON объект то врнет числовое значение ключа
//  any := VkApi.AnyModel(`{"name":"Ivan","age":15}`)
//  age, err := any.GetInt("age") // age == 15 err == nil
func (m *AnyModel) GetInt(key string) (int, error) {
	x := make(map[string]json.RawMessage, 0)
	err := json.Unmarshal([]byte(*m), &x)
	if err != nil {
		return 0, err
	}
	data, has := x[key]
	if has {
		am := AnyModel(data)
		return am.Int()
	} else {
		return 0, errors.New(keyNotExsits(key))
	}
}

func (m *AnyModel) GetSlice(key string) ([]*AnyModel, error) {
	x := make(map[string]json.RawMessage, 0)
	err := json.Unmarshal([]byte(*m), &x)
	if err != nil {
		return []*AnyModel{}, err
	}
	data, has := x[key]
	if has {
		am := AnyModel(data)
		return am.Slice()
	} else {
		return []*AnyModel{}, errors.New(keyNotExsits(key))
	}
}

func (m *AnyModel) GetIntDef(key string, def int) (int, error) {
	res, err := m.GetInt(key)
	if err == nil {
		return res, nil
	} else if err.Error() == keyNotExsits(key) {
		return def, nil
	} else {
		return def, err
	}
}

func (m *AnyModel) Slice() ([]*AnyModel, error) {
	i := make([]json.RawMessage, 0)
	r := make([]*AnyModel, 0)
	err := json.Unmarshal([]byte(*m), &i)

	if err != nil {
		return r, err
	} else {
		for _, raw := range i {
			m := AnyModel(raw)
			r = append(r, &m)
		}
		return r, nil
	}
}

func (m *AnyModel) QAny(s string) *AnyModel {
	parts := strings.Split(s, ".")
	x := m
	for _, key := range parts {
		index, indexErr := strconv.Atoi(key)
		if indexErr == nil {
			slice, err := x.Slice()
			if err == nil && (index >= 0 && len(slice) > index) || (index < 0 && len(slice)+index >= 0) {
				if index < 0 {
					index = len(slice) + index
				}
				x = slice[index]
				continue
			}
		}

		ret, err := x.GetAny(key)
		if err != nil {
			return nil
		} else {
			x = ret
		}
	}
	return x
}

func (m *AnyModel) QString(s string) *string {
	any := m.QAny(s)
	if any == nil {
		return nil
	}
	s, err := any.String()
	if err == nil {
		return &s
	} else {
		return nil
	}
}

func (m *AnyModel) QStringDef(s, def string) string {
	res := m.QString(s)
	if res == nil {
		return def
	}
	return *res
}

func (m *AnyModel) QInt(s string) *int {
	any := m.QAny(s)
	if any == nil {
		return nil
	}
	i, err := any.Int()
	if err == nil {
		return &i
	} else {
		return nil
	}
}

func (m *AnyModel) QIntDef(s string, def int) int {
	res := m.QInt(s)
	if res == nil {
		return def
	}
	return *res
}

func (m *AnyModel) QBool(s string) *bool {
	any := m.QAny(s)
	if any == nil {
		return nil
	}
	i, err := any.Bool()
	if err == nil {
		return &i
	} else {
		return nil
	}
}

func (m *AnyModel) QBoolDef(s string, def bool) bool {
	res := m.QBool(s)
	if res == nil {
		return def
	}
	return *res
}

func (m *AnyModel) QFloat(s string) *float64 {
	any := m.QAny(s)
	if any == nil {
		return nil
	}
	i, err := any.Float()
	if err == nil {
		return &i
	} else {
		return nil
	}
}

func (m *AnyModel) QFloatDef(s string, def float64) float64 {
	res := m.QFloat(s)
	if res == nil {
		return def
	}
	return *res
}

func (m *AnyModel) QSlice(s string) *[]*AnyModel {
	any := m.QAny(s)
	if any == nil {
		return nil
	}
	a, err := any.Slice()
	if err == nil {
		return &a
	} else {
		return nil
	}
}

func (m *AnyModel) QSliceDef(s string, def []*AnyModel) []*AnyModel {
	res := m.QSlice(s)
	if res == nil {
		return def
	}
	return *res
}

func (m *AnyModel) Unmarshal(i interface{}) error {
	return json.Unmarshal([]byte(*m), i)
}
