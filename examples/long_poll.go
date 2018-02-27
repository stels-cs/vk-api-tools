package main

import (
	"github.com/stels-cs/vk-api-tools"
	"log"
	"os"
	"os/signal"
	"strconv"
)

type Bot struct {
	token string
}

func (b *Bot) NewMessage(event VkApi.MessageEvent) {
	if event.IsOutMessage() {
		return
	}
	u := make([]struct {
		FirstName string `json:"first_name"`
	}, 1)
	err := VkApi.Exec("users.get", VkApi.P{"user_id": strconv.Itoa(event.PeerId)}, &u)
	if err != nil {
		println(err.Error())
		return
	}
	msg := "Привет, " + u[0].FirstName
	_, err = VkApi.Call("messages.send", VkApi.P{
		"peer_id":      strconv.Itoa(event.PeerId),
		"message":      msg,
		"access_token": b.token,
	})
	if err != nil {
		println(err.Error())
		return
	}
}

func (b *Bot) EditMessage(event VkApi.MessageEvent) {
	msg := "Не редактируйте сообщения прлиз"
	_, err := VkApi.Call("messages.send", VkApi.P{
		"peer_id":      strconv.Itoa(event.PeerId),
		"message":      msg,
		"access_token": b.token,
	})
	if err != nil {
		println(err.Error())
		return
	}
}

func main() {
	token := "529d99ca66f012327d76df9a9691bd082c92658d7c12addf77882c388fcfa6936c88bf19da3b0e717e7df"
	bot := &Bot{token}
	logger := log.New(os.Stdout, "LongPoll", log.Lshortfile)
	api := VkApi.CreateApi(token, "5.71", VkApi.GetHttpTransport(), 30)
	lp := VkApi.GetLongPollServer(api, logger)
	lp.SetListener(bot)
	go lp.Start()
	println("Press Ctl+C for quit")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	lp.Stop()
}
