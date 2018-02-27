package main

import (
	"github.com/stels-cs/vk-api-tools"
	"strconv"
	"time"
)

func main() {

	token := "529d99ca66f012327d76df9a9691bd082c92658d7c12addf77882c388fcfa6936c88bf19da3b0e717e7df"
	v := "5.71"
	api := VkApi.CreateApi(token, v, VkApi.GetHttpTransport(), 30)
	requestsPerSecond := 1
	rq := VkApi.GetRequestQueue(api, requestsPerSecond)
	go rq.Start()
	defer rq.Stop()

	l1 := make(chan int)
	l2 := make(chan int)

	var diff1 int64
	var diff3 int64
	var r0, r1, r2, r3 VkApi.RequestResult

	go func() {
		ts1 := time.Now().UnixNano()
		c0 := rq.Call(VkApi.CreateMethod("users.get", VkApi.P{"user_ids": "1"}))
		r0 = <-c0
		diff1 = time.Now().UnixNano() - ts1
		l1 <- 1
	}()

	go func() {
		ts3 := time.Now().UnixNano()
		c1 := rq.Call(VkApi.CreateMethod("users.get", VkApi.P{"user_ids": "2050"}))
		c2 := rq.Call(VkApi.CreateMethod("users.get", VkApi.P{"user_ids": "avk"}))
		c3 := rq.Call(VkApi.CreateMethod("users.get", VkApi.P{"user_ids": "andrew"}))

		r1 = <-c1
		r2 = <-c2
		r3 = <-c3
		diff3 = time.Now().UnixNano() - ts3
		l2 <- 1
	}()

	<-l1
	<-l2
	println("One request: " + strconv.Itoa(int(diff1/int64(time.Nanosecond))) + "ns")
	println(r0.Res.QStringDef("0.first_name", "") + "\n")
	println("Three request: " + strconv.Itoa(int(diff3/int64(time.Nanosecond))) + "ns")
	println(r1.Res.QStringDef("0.first_name", ""))
	println(r2.Res.QStringDef("0.first_name", ""))
	println(r3.Res.QStringDef("0.first_name", ""))

}
