package Vk

import (
	"testing"
	"github.com/stels-cs/quiz-bot/Vk"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

func checkKatrya(userArr []Vk.User) error {
	if len(userArr) != 1 {
		return errors.New(fmt.Sprintf("Expect 1 user in array but got %d", len(userArr)))
	}

	if userArr[0].Id != 2050 {
		return errors.New(fmt.Sprintf("Incorrect user id expect 2050 got %d", userArr[0].Id))
	}

	return nil
}

func TestExecutePack(t *testing.T) {
	pack := Vk.ExecutePack{}
	index, err := pack.Add(Vk.GetApiMethod("users.get", Vk.Params{"user_ids": "2050"}))
	if err != nil {
		t.Error(err)
	}
	if index != 0 {
		t.Errorf("Bad return index, expect 0 but got %d", index)
	}

	code := pack.GetCode()

	trueCode := `return[API.users.get({"user_ids":"2050"})];`

	if code != trueCode {
		t.Errorf("Bad return code, expect \nTRUE: %s\nRETN: %s", trueCode, code)
	}
}

func TestExecutePackCall(t *testing.T) {
	pack := Vk.ExecutePack{}
	index, err := pack.Add(Vk.GetApiMethod("users.get", Vk.Params{"user_ids": "2050"}))
	if err != nil {
		t.Error(err)
	}
	api := Vk.GetApi(Vk.AccessToken{Token: "7e3d0d54ec8d5050ef540d16fd978fe269f615de0ced94b9f97553b245727ab3844b08572efdbaa8dde8a"}, Vk.GetHttpTransport(), nil)
	res, err := api.Execute(pack.GetCode())
	if err != nil {
		t.Error(err)
		return
	}
	var data []json.RawMessage
	err = json.Unmarshal(*res.Response, &data)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 1 {
		t.Errorf("Invalid response expect 1 got %d", len(data))
	}
	if len(res.ExecuteErrors) != 0 {
		t.Errorf("There are %d execure errors on good request", len(res.ExecuteErrors))
	}
	userRes := data[index]
	var users []Vk.User
	err = json.Unmarshal(userRes, &users)
	if err != nil {
		t.Error(err)
	}
	err = checkKatrya(users)
	if err != nil {
		t.Error(err)
	}
}


func request(rq *Vk.RequestQueue, t *testing.T) {
	res := <-rq.Call(Vk.GetApiMethod("users.get", Vk.Params{"user_ids": "2050"}))
	if res.Err != nil {
		t.Error(res.Err)
		return
	}

	var users []Vk.User
	err := json.Unmarshal(*res.Res.Response, &users)

	if err != nil {
		t.Error(err)
		return
	}

	err = checkKatrya(users)
	if err != nil {
		t.Error(err)
	}
}


func TestRequestQueue(t *testing.T) {
	rq := Vk.GetRequestQueue(Vk.GetApi(Vk.AccessToken{Token: "7e3d0d54ec8d5050ef540d16fd978fe269f615de0ced94b9f97553b245727ab3844b08572efdbaa8dde8a"}, Vk.GetHttpTransport(), nil))
	go rq.Start()
	defer rq.Stop()

	request(rq, t)
}

func TestRequestQueueManyRequest(t *testing.T) {
	rq := Vk.GetRequestQueue(Vk.GetApi(Vk.AccessToken{Token: "7e3d0d54ec8d5050ef540d16fd978fe269f615de0ced94b9f97553b245727ab3844b08572efdbaa8dde8a"}, Vk.GetHttpTransport(), nil))
	go rq.Start()
	defer rq.Stop()

	end1 := make(chan bool, 1)
	end2 := make(chan bool, 1)
	end3 := make(chan bool, 1)

	go request(rq, t)
	time.Sleep(time.Millisecond)
	go request(rq, t)
	time.Sleep(time.Millisecond)
	go request(rq, t)
	time.Sleep(time.Millisecond)

	go func () {
		request(rq, t)
		end1 <- true
	} ()

	go func () {
		request(rq, t)
		end2 <- true
	} ()

	go func () {
		request(rq, t)
		end3 <- true
	} ()

	<-end1
	ts := time.Now().UnixNano()
	<-end2
	<-end3
	diff := time.Now().UnixNano() - ts
	if diff > int64(100 * time.Millisecond) {
		t.Errorf("Multi request not stacked %dns", diff)
	}
}

func TestRequestQueueOneGoodOneFail(t *testing.T) {
	rq := Vk.GetRequestQueue(Vk.GetApi(Vk.AccessToken{Token: "7e3d0d54ec8d5050ef540d16fd978fe269f615de0ced94b9f97553b245727ab3844b08572efdbaa8dde8a"}, Vk.GetHttpTransport(), nil))
	go rq.Start()
	defer rq.Stop()

	var end1 chan Vk.RequestResult
	var end2 chan Vk.RequestResult

	go request(rq, t)
	time.Sleep(time.Millisecond)
	go request(rq, t)
	time.Sleep(time.Millisecond)
	go request(rq, t)
	time.Sleep(time.Millisecond)

	end1 = rq.Call(Vk.GetApiMethod("users.get", Vk.Params{"user_ids": "2050"}))
	end2 = rq.Call(Vk.GetApiMethod("messages.send", Vk.Params{"peer_id": "1", "message":"test"}))

	res1 := <- end1

	if res1.Err != nil {
		t.Error(res1.Err)
		return
	}

	var users []Vk.User
	err := json.Unmarshal(*res1.Res.Response, &users)

	if err != nil {
		t.Error(err)
		return
	}

	err = checkKatrya(users)
	if err != nil {
		t.Error(err)
	}

	res2 := <- end2

	if res2.Err == nil {
		t.Error("Not error on bad request")
		return
	}

	if api, ok:= res2.Err.(*Vk.ApiError); ok {
		if api.Code != 15 {
			t.Errorf("Wrong error code expext 15 got %d %s", api.Code, api.Error())
		}
	} else {
		t.Error("Wrong error type")
	}
}

func TestRequestQueueOneFailOneGood(t *testing.T) {
	rq := Vk.GetRequestQueue(Vk.GetApi(Vk.AccessToken{Token: "7e3d0d54ec8d5050ef540d16fd978fe269f615de0ced94b9f97553b245727ab3844b08572efdbaa8dde8a"}, Vk.GetHttpTransport(), nil))
	go rq.Start()
	defer rq.Stop()

	var end1 chan Vk.RequestResult
	var end2 chan Vk.RequestResult

	go request(rq, t)
	time.Sleep(time.Millisecond)
	go request(rq, t)
	time.Sleep(time.Millisecond)
	go request(rq, t)
	time.Sleep(time.Millisecond)

	end2 = rq.Call(Vk.GetApiMethod("messages.send", Vk.Params{"peer_id": "1", "message":"test"}))
	end1 = rq.Call(Vk.GetApiMethod("users.get", Vk.Params{"user_ids": "2050"}))

	res1 := <- end1

	if res1.Err != nil {
		t.Error(res1.Err)
		return
	}

	var users []Vk.User
	err := json.Unmarshal(*res1.Res.Response, &users)

	if err != nil {
		t.Error(err)
		return
	}

	err = checkKatrya(users)
	if err != nil {
		t.Error(err)
	}

	res2 := <- end2

	if res2.Err == nil {
		t.Error("Not error on bad request")
		return
	}

	if api, ok:= res2.Err.(*Vk.ApiError); ok {
		if api.Code != 15 {
			t.Errorf("Wrong error code expext 15 got %d %s", api.Code, api.Error())
		}
	} else {
		t.Error("Wrong error type")
	}
}

func TestRequestQueueOneFail(t *testing.T) {
	rq := Vk.GetRequestQueue(Vk.GetApi(Vk.AccessToken{Token: "7e3d0d54ec8d5050ef540d16fd978fe269f615de0ced94b9f97553b245727ab3844b08572efdbaa8dde8a"}, Vk.GetHttpTransport(), nil))
	go rq.Start()
	defer rq.Stop()

	var end2 chan Vk.RequestResult

	go request(rq, t)
	time.Sleep(time.Millisecond)
	go request(rq, t)
	time.Sleep(time.Millisecond)
	go request(rq, t)
	time.Sleep(time.Millisecond)

	end2 = rq.Call(Vk.GetApiMethod("messages.send", Vk.Params{"peer_id": "1", "message":"test"}))
	//end1 = rq.run(Vk.GetApiMethod("users.get", Vk.Params{"user_ids": "2050"}))

	res2 := <- end2

	if res2.Err == nil {
		t.Error("Not error on bad request")
		return
	}

	if api, ok:= res2.Err.(*Vk.ApiError); ok {
		if api.Code != 15 {
			t.Errorf("Wrong error code expext 15 got %d %s", api.Code, api.Error())
		}
	} else {
		t.Error("Wrong error type")
	}
}

func TestRequestQueueTwoFailOneGood(t *testing.T) {
	rq := Vk.GetRequestQueue(Vk.GetApi(Vk.AccessToken{Token: "7e3d0d54ec8d5050ef540d16fd978fe269f615de0ced94b9f97553b245727ab3844b08572efdbaa8dde8a"}, Vk.GetHttpTransport(), nil))
	go rq.Start()
	defer rq.Stop()

	var end1 chan Vk.RequestResult
	var end2 chan Vk.RequestResult
	var end3 chan Vk.RequestResult

	go request(rq, t)
	time.Sleep(time.Millisecond)
	go request(rq, t)
	time.Sleep(time.Millisecond)
	go request(rq, t)
	time.Sleep(time.Millisecond)

	end2 = rq.Call(Vk.GetApiMethod("messages.send", Vk.Params{"peer_id": "1", "message":"test"}))
	end3 = rq.Call(Vk.GetApiMethod("messages.send", Vk.Params{"peer_id": "1", "message":"test"}))
	end1 = rq.Call(Vk.GetApiMethod("users.get", Vk.Params{"user_ids": "2050"}))

	res1 := <- end1

	if res1.Err != nil {
		t.Error(res1.Err)
		return
	}

	var users []Vk.User
	err := json.Unmarshal(*res1.Res.Response, &users)

	if err != nil {
		t.Error(err)
		return
	}

	err = checkKatrya(users)
	if err != nil {
		t.Error(err)
	}

	res2 := <- end2

	if res2.Err == nil {
		t.Error("Not error on bad request")
		return
	}

	if api, ok:= res2.Err.(*Vk.ApiError); ok {
		if api.Code != 15 {
			t.Errorf("Wrong error code expext 15 got %d %s", api.Code, api.Error())
		}
	} else {
		t.Error("Wrong error type")
	}

	res3 := <- end3

	if res3.Err == nil {
		t.Error("Not error on bad request")
		return
	}

	if api, ok:= res3.Err.(*Vk.ApiError); ok {
		if api.Code != 15 {
			t.Errorf("Wrong error code expext 15 got %d %s", api.Code, api.Error())
		}
	} else {
		t.Error("Wrong error type")
	}
}
