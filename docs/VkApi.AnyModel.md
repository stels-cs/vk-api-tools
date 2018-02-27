# AnyModel

Тип для работы с JSON произвольной структуры

```go
type AnyModel json.RawMessage
```

[Посмотрите тесты, чтобы лучше понять как это работает](../tests/AnyModel_test.go)

Например, есть json вида:

```json
{
  "first_name":"Ivan",
  "last_name":"Ivanov",
  "city":{
    "id":19,
    "title":"Moscow"
  },
  "friends":[ {"id":1,"name":"Jon","verify":true},{"id":2,"name":"Bob","verify":false},{"id":3,"name":"Alexandra"} ]
}
```

и мы хотим узнать из него

- значение ключа ```first_name```
- название города
- проверить есть ли ключ ```age```
- получить ```name``` первого и ```id``` последнего друга
- вывести имена всех друзей у кого есть ключ ```verify``` и его значение ```true``` или этого ключа нет

### Создадим объект AnyModel

```go
any := VkApi.AnyModel(`{
  "first_name":"Ivan",
  "last_name":"Ivanov",
  "city":{
    "id":19,
    "title":"Moscow"
  },
  "friends":[ {"id":1,"name":"Jon","verify":true},{"id":2,"name":"Bob","verify":false},{"id":3,"name":"Alexandra"} ]
}`)
```

### Получим значение ключа ```first_name``` и ```title``` из объекта ```city```

```go
firstName := any.QStringDef("first_name", "")
cityName := any.QStringDef("city.title", "")
```

вторым аргуметом мы передали значение по умолчанмю, можно использовать такой подход

```go
firstName := any.QString("first_name")
```

так мы получим указатель на строку или nil если такого ключа нет

### Проверим есть ли ключ ```age```

```go
hasAge := any.QString("age") != nil
```


### Получим ```name``` первого и ```id``` последнего друга

```go
firstFriendName := any.QStringDef("friends.0.name", "")
lastFriendId := any.QIntDef("friends.-1.id", -1)
```


### Выведем имена всех друзей у который есть ключ ```verify``` и его значение ```true``` или этого ключа нет

```go
friends := any.QSlice("friends") // Указатель на массив *AnyModel

buff := ""
if friends != nil {
    for _, friend := range *friends {
        v := friend.QBool("verify")

        if v == nil || *v == true { // если v == nil то такого ключа небыло
            println(friend.QStringDef("name", ""))
        }
    }
}
```

[Назад](../README.md)