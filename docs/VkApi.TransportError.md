# VkApi.TransportError

```go
type TransportError struct {
	Method      string // Вызываемый метод
	Params      P // Вызываемые параметры
	Response    []byte // Ответ от сервера ВК (если был)
	Headers     TransportExternalData // Заголовки от сервера ВК (если были)
	ParentError error //Исходная ошибка (например, тут может быть ошибка сети или таймаут)
}
```

[Назад](../README.md)