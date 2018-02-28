# VkApi

Выполняет запрс к API ВКонтакте, в случае сетевых ошибок будет повторят запрос с интервалом в 1 секунду не более 30 раз

```go
func VkApi.Exec(method string, params VkApi.P, i interface{}) error
```

- ```method``` - имя метода, например ```"users.get"```
- ```params``` - параметры, напримре ```VkApi.P{"user_ids":"2050"}```
- ```i``` - указатель на объект в который будет записан результат выполнения

Выполняет запрс к API ВКонтакте, в случае сетевых ошибок будет повторят запрос с интервалом в 1 секнду не более 30 раз.
В случае успеха вернет объект [VkApi.Response](VkApi.Response.md)

```go
func VkApi.Call(method string, params VkApi.P) (VkApi.Response, error)
```

- ```method``` - имя метода, например ```"users.get"```
- ```params``` - параметры, напримре ```VkApi.P{"user_ids":"2050"}```


### Обработка ошибок

Метод VkApi.Call и VkApi.Exec могут вернуть три типа ошибок

- [VkApi.TransportError](VkApi.TransportError.md) - Произошла ошибка сети, или сервреа ВКонтакте временно недоступны 
- [VkApi.ApiError](VkApi.ApiError.md) - API ВКонтакте вернуло ошибку
- ```все остальные типы``` - в основном ошибки парсинга json, такое можно вернуть только VkApi.Exec

[Назад](../README.md)