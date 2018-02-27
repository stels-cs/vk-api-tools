package main

import (
	"github.com/stels-cs/vk-api-tools"
	"strconv"
)

type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func main() {
	users := make([]User, 0)
	err := VkApi.Exec("users.get", VkApi.P{"user_ids": "2050,andrew"}, &users)
	if err != nil {
		panic(err)
	}

	for _, u := range users {
		println(u.FirstName + " " + u.LastName + " #" + strconv.Itoa(u.Id))
	}
	//Катя Лебедева #2050
	//Андрей Рогозов #6492
}
