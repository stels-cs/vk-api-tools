# Golang API Вконтакте

```bash
go get git@github.com:stels-cs/vk-api-tools.git
```

### Базовый пример использования:

````go
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
}
````
Результат:
```bash
Катя Лебедева #2050
Андрей Рогозов #6492
```

### Пример использования без структур

```go
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
}
``` 
Результат:
```bash
Катя – Санкт-Петербург
Александр – Москва
```

**Важно!** По умолчанию VkApi.Exec и VkApi.Call будут повторят запрос до 30 раз в случае любых сетевых ошибок или если от API придет один из следующих кодов ошибки: 1, 9, 6, 9, 10, 603 [vk.com/dev/errors](https://vk.com/dev/errors).
Чтобы отключить это посмотрите пример ниже.

### Пример создания объекта api

```go
package main

import (
    "github.com/stels-cs/vk-api-tools"
    "strconv"
)

func main() {
    token := ""
    v := "5.71"
    retryTimesIfFailed := 0 //Не повторять запросы в случае любых ошибок

    api := VkApi.CreateApi(token, v, VkApi.GetHttpTransport(), retryTimesIfFailed)

    users := make([]struct{
        Id int `json:"id"`
        FirstName string `json:"first_name"`
        LastName string `json:"last_name"`
    }, 0)
    err := api.Exec("users.get", VkApi.P{"user_ids": "2050,andrew"}, &users)
    if err != nil {
        panic(err)
    }

    for _, u := range users {
        println(u.FirstName + " " + u.LastName + " #" + strconv.Itoa(u.Id))
    }
}
```
Результат:
```bash
Катя Лебедева #2050
Андрей Рогозов #6492
```