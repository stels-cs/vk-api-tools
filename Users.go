package Vk

import (
	"encoding/json"
	"errors"
)

type Users struct {
	api *Api
}

type User struct {
	Id int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Sex int `json:"sex"`
}

func (m *Users) GetMe() (User, error) {
	user := User{}
	userList := []User{user}
	response, err := m.api.run("users.get", Params{
		"fields":"sex",
	}, 0)
	if err != nil {
		return user, err
	}
	err = json.Unmarshal(*response.Response, &userList)
	if err != nil {
		return user, err
	} else {
		if len(userList) > 0 {
			return userList[0], nil
		} else {
			return user, errors.New("Empty response from users.get({})\n")
		}
	}
}

func (m *Users) GetByIds(ids []int) ([]User, error) {
	var items []User
	response, err := m.api.run("users.get", Params{
		"user_ids": intToString(ids),
		"fields":"sex",
	}, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal(*response.Response, &items)
	if err != nil {
		return items, err
	} else {
		return items, nil
	}
}
