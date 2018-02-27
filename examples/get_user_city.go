package main

import (
	"github.com/stels-cs/vk-api-tools"
)

func main() {
	res, err := VkApi.Call("users.get", VkApi.P{"user_ids": "2050,avk", "fields": "city"})
	if err != nil {
		panic(err)
	}

	print(res.QStringDef("0.first_name", "") + " – ")
	println(res.QStringDef("0.city.title", ""))

	print(res.QStringDef("1.first_name", "") + " – ")
	println(res.QStringDef("1.city.title", ""))
	//Катя – Санкт-Петербург
	//Александр – Москва
}
