package Vk

import (
	"encoding/json"
	"errors"
)

type Groups struct {
	api *Api
}

type Group struct {
	Id int `json:"id"`
	Name string `json:"name"`
	ScreenName string `json:"screen_name"`
}

func (m *Groups) GetMe() (Group, error) {
	group := Group{}
	groupList := []Group{group}
	response, err := m.api.run("groups.getById", Params{}, 0)
	if err != nil {
		return group, err
	}
	err = json.Unmarshal(*response.Response, &groupList)
	if err != nil {
		return group, err
	} else {
		if len(groupList) > 0 {
			return groupList[0], nil
		} else {
			return group, errors.New("Empty response from groups.getById({})\n")
		}
	}
}
