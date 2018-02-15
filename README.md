# Go tools for work with Vk Api

Example:

````go
package main

import (
	"github.com/stels-cs/vk-api-tools"
)

func main() {
	api := Vk.GetApi(Vk.AccessToken{}, Vk.GetHttpTransport(), nil)
	users, err := api.Users.GetByIds([]int{1})
	if err != nil {
        println(err)
        return
    }
    println(users[0].FirstName)
}
````

