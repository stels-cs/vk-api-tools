package main

import (
	"github.com/stels-cs/vk-api-tools"
	"strconv"
)

func main() {
	token := "529d99ca66f012327d76df9a9691bd082c92658d7c12addf77882c388fcfa6936c88bf19da3b0e717e7df"

	ep := VkApi.ExecutePack{}

	p1 := VkApi.CreateMethod("groups.getMembers", VkApi.P{"group_id": "mudakoff", "sort": "id_asc", "offset": "0", "count": "1000"})
	p2 := VkApi.CreateMethod("groups.getMembers", VkApi.P{"group_id": "mudakoff", "sort": "id_asc", "offset": "1000", "count": "1000"})
	p3 := VkApi.CreateMethod("groups.getMembers", VkApi.P{"group_id": "mudakoff", "sort": "id_asc", "offset": "2000", "count": "1000"})

	ep.Add(p1)
	ep.Add(p2)
	ep.Add(p3)

	println(ep.GetCode())

	res, err := VkApi.Call("execute", VkApi.P{"code": ep.GetCode(), "access_token": token})
	if err != nil {
		panic(err)
	}

	slice, err := res.Slice()
	if err != nil {
		panic(err)
	}

	userIdsInGroup := make([]int, 0)

	for _, pack := range slice {
		p := make([]int, 0)
		s, _ := pack.GetAny("items")
		err := s.Unmarshal(&p)
		if err != nil {
			panic(err)
		}
		userIdsInGroup = append(userIdsInGroup, p...)
	}

	println("Loaded ids: " + strconv.Itoa(len(userIdsInGroup)))
}
