# VkApi.Response

```go
type Response struct {
	Response      *json.RawMessage `json:"response"`
	Error         ApiError         `json:"error"`
	ExecuteErrors []struct {
		Method  string `json:"method"`
		Code    int    `json:"error_code"`
		Message string `json:"error_msg"`
	} `json:"execute_errors"`
}
```

```go
func (r Response) Any() *VkApi.AnyModel
```

Создает объект [VkApi.AnyModel](VkApi.AnyModel.md) из r.Response и возворащает ссылку на него.


```go
func (r Response) FirstAny() (*VkApi.AnyModel, error)
```

Создает объект AnyModel из r.Response и возворащает ссылку перый элемент массива.


```go
res, err := VkApi.Call("users.get", VkApi.P{"user_ids": "2050,avk", "fields": "city"})
if err != nil {
    panic(err)
}
u, _ := res.FirstAny()
name, _ := u.GetString("first_name")
//name == "Катя"
```