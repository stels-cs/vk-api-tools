package main

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/stels-cs/vk-api-tools"
	"io/ioutil"
	"time"
)

type WsConnection struct {
	c *websocket.Conn
}

func (conn *WsConnection) Close() error {
	if conn.c == nil {
		return nil
	}
	return conn.c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func (conn *WsConnection) ReadMessage() ([]byte, error) {
	_, message, err := conn.c.ReadMessage()
	return message, err
}

func GetConnection(url string) (VkApi.ConnectionInterface, error) {
	conn, response, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		body, e := ioutil.ReadAll(response.Body)
		if e == nil {
			response.Body.Close()
			return nil, errors.New(string(body))
		}
		return nil, err
	}

	return &WsConnection{
		c: conn,
	}, nil
}

func IsCloseError(err error) bool {
	if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
		return true
	} else {
		return false
	}
}

func main() {
	// Понадобиться сервиный ключ приложения
	token := "60fb7fa360fb7fa360fb7fa34b60a5f1e7660fb60fb7fa33a685687cb8b239d9a73303c"

	// Создаем объект VkApi.StreamingClient
	streamingApi, err := VkApi.CreateStreamingClient(token)
	if err != nil {
		panic(err)
	}

	// Получаем все правила которые есть сейчас для этого ключа
	rules, err := streamingApi.GetRules()
	if err != nil {
		panic(err)
	}

	// Удаляем все правила которые были
	for _, rule := range rules {
		println("Rule:", rule.Tag, "Value:", rule.Value)
		err := streamingApi.DeleteRule(rule.Tag)
		if err != nil {
			panic(err)
		}
	}

	// Создаем свои правила
	rule1 := VkApi.StreamingRule{
		Value: "путин",
		Tag:   "putin",
	}

	rule2 := VkApi.StreamingRule{
		Value: "новости",
		Tag:   "news",
	}

	rule3 := VkApi.StreamingRule{
		Value: "биткойн",
		Tag:   "bitcoin",
	}

	err = streamingApi.AddRule(rule1)
	if err != nil {
		panic(err)
	}
	err = streamingApi.AddRule(rule2)
	if err != nil {
		panic(err)
	}
	err = streamingApi.AddRule(rule3)
	if err != nil {
		panic(err)
	}

	println("Rules created")
	go func() {
		time.Sleep(30 * time.Second)
		streamingApi.Stop() // Вот так можно остановить прослушивание сообщений
		println("Stop streaming by timeout")
	}()

	println("Start lister event")
	err = streamingApi.Start(GetConnection, IsCloseError, func(code int, event *VkApi.StreamingEvent, message *VkApi.StreamingServiceMessage) {
		if event != nil {
			println(event.Text, event.EventUrl)
		}
	})

	if err != nil {
		panic(err)
	}
}
