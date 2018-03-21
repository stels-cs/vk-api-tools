package main

import (
	"github.com/stels-cs/vk-api-tools"
	"time"
)

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
		Value: "кандидиат",
		Tag:   "candidate",
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
	err = streamingApi.Start(func(code int, event *VkApi.StreamingEvent, message *VkApi.StreamingServiceMessage) {
		if event != nil {
			println(event.Text, event.EventUrl)
		}
	})

	if err != nil {
		panic(err)
	}
}
