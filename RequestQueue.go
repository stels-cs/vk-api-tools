package Vk

import (
	"sync"
	"time"
	"errors"
	"encoding/json"
)

const requestPerSecond = 3

type RequestResult struct {
	Res *ApiResponse
	Err error
}

type RequestItem struct {
	r     ApiMethod
	ch    chan RequestResult
	index int
}

type RequestQueue struct {
	api              *Api
	queue            []RequestItem
	lock             sync.Mutex
	stop             chan bool
	item             chan RequestItem
	lastExecuteTime  int64
	lastExecuteCount int
	timer            *time.Timer
	timerUp          bool
}

func GetRequestQueue(api *Api) *RequestQueue {
	return &RequestQueue{
		api:   api,
		lock:  sync.Mutex{},
		stop:  make(chan bool, 1),
		item:  make(chan RequestItem, 25),
		timer: time.NewTimer(0),
	}
}

func (rq *RequestQueue) Call(m ApiMethod) chan RequestResult {
	ch := make(chan RequestResult, 1)
	rq.item <- RequestItem{m, ch, 0}
	return ch
}

func (rq *RequestQueue) Start() {
	for {
		select {
		case <-rq.stop:
			return
		case item := <-rq.item:
			rq.toQueue(item)
			rq.execute()
		case <-rq.timer.C:
			rq.execute()
		}
	}
}

func (rq *RequestQueue) Stop() {
	rq.stop <- true
}

func (rq *RequestQueue) pass() bool {
	now := time.Now().Unix()
	if rq.lastExecuteTime != now {
		rq.lastExecuteTime = now
		rq.lastExecuteCount = 1
		return true
	} else {
		rq.lastExecuteCount++
		if rq.lastExecuteCount > requestPerSecond {
			return false
		} else {
			return true
		}
	}
}
func (rq *RequestQueue) execute() {
	if len(rq.queue) == 0 {
		return
	}
	if !rq.pass() {
		rq.upTimer()
		return
	}
	rq.timerUp = false
	pack := ExecutePack{}
	responseMap := map[int]RequestItem{}
	for len(rq.queue) > 0 {
		item := rq.queue[0]
		index, err := pack.Add(item.r)
		if err != nil {
			item.ch <- RequestResult{nil, err}
			rq.queue = rq.queue[1:]
		} else if index != -1 {
			item.index = index
			rq.queue = rq.queue[1:]
			responseMap[item.index] = item
		} else if pack.Count() == 0 {
			item.ch <- RequestResult{nil, errors.New("Request too large ")}
			rq.queue = rq.queue[1:]
		} else {
			break
		}
	}
	if len(rq.queue) > 0 {
		rq.upTimer()
	}
	go func() {
		code := pack.GetCode()
		response, err := rq.api.Execute(code)
		if err != nil {
			for _, item := range responseMap {
				item.ch <- RequestResult{nil, err}
			}
			return
		}
		errList := response.ExecuteErrors
		var resList []json.RawMessage
		err = json.Unmarshal(*response.Response, &resList)
		if err != nil {
			for _, item := range responseMap {
				item.ch <- RequestResult{nil, err}
			}
			return
		}
		for index, res := range resList {
			item, ok := responseMap[index]
			if !ok {
				continue
			}
			if isBoolAndFalse(&res) {
				if len(errList) > 0 {
					err := errList[0]
					errList = errList[1:]
					apiError := &ApiError{
						Code:       err.Code,
						Message:    err.Method + " " + err.Message,
						Params:     []RequestedParams{},
						CallMethod: &item.r.name,
						CallParams: &item.r.params,
					}
					item.ch <- RequestResult{nil, apiError}
				} else {
					item.ch <- RequestResult{nil, errors.New("Cant get execute error from response " + string(*response.Response))}
				}
			} else {
				copyRes := make(json.RawMessage, len(res))
				copy(copyRes, res)
				apiResponse := ApiResponse{
					Response: &copyRes,
					Error: ApiError{
						Code: 0,
					},
				}
				item.ch <- RequestResult{&apiResponse, nil}
			}
		}
	}()
}
func (rq *RequestQueue) toQueue(item RequestItem) {
	rq.queue = append(rq.queue, item)
}
func (rq *RequestQueue) upTimer() {
	if rq.timerUp {
		return
	}
	rq.timerUp = true
	if rq.lastExecuteCount > requestPerSecond {
		now := time.Now().Nanosecond()
		rq.timer.Reset(time.Duration(int64(time.Second) - int64(now)))
	} else {
		rq.timer.Reset((1000 / requestPerSecond) * time.Millisecond)
	}
}
