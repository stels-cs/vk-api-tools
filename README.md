[![License](http://img.shields.io/badge/license-MIT-brightgreen.svg)](https://tldrlegal.com/license/mit-license)
[![Build Status](https://travis-ci.org/stels-cs/vk-api-tools.svg?branch=master)](https://travis-ci.org/stels-cs/vk-api-tools)

# Golang API ВКонтакте

```bash
go get git@github.com:stels-cs/vk-api-tools.git
```

- [VkApi](docs/VkApi.md)
- [VkApi.Response](docs/VkApi.Response.md)
- [VkApi.ApiError](docs/VkApi.ApiError.md)
- [VkApi.TransportError](docs/VkApi.TransportError.md)
- [VkApi.AnyModel](docs/VkApi.AnyModel.md)

### Пример:

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

### Пример без структур:

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

Подробнее про [VkApi.Response](docs/VkApi.Response.md)
и про [QStringDef](docs/VkApi.AnyModel.md)

**Важно!** По умолчанию VkApi.Exec и VkApi.Call будут повторят запрос до 30 раз в случае любых сетевых ошибок или если API вернет кодо ошибки: 1, 9, 6, 9, 10, 603 [vk.com/dev/errors](https://vk.com/dev/errors).
Чтобы отключить это посмотрите пример ниже.

### Пример создания объекта api

```go
package main

import (
    "github.com/stels-cs/vk-api-tools"
    "strconv"
)

func main() {
    token := "3bac432bdcb1234b1...."  //API ключ доступа
    v := "5.71" //Версия API по умолчанию
    retryTimesIfFailed := 0 //Не повторять запросы в случае любых ошибок, можно поставить 5, тогда запрос будет повторен 5 раз в случае ошибок

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

### Обработка ошибок

Метод VkApi.Call и VkApi.Exec могут вернуть три типа ошибок

- [VkApi.TransportError](docs/VkApi.TransportError.md) - Произошла ошибка сети, или сервер ВКонтакте временно недоступен 
- [VkApi.ApiError](docs/VkApi.ApiError.md) - API ВКонтакте вернуло ошибку
- ```все остальные типы``` - ошибки парсинга json, только для VkApi.Exec