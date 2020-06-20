package VkApi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type BotLongPollServer struct {
	Key        string
	Server     string
	Ts         int
	httpClient *http.Client
	Logger     *log.Logger
	intSize    int
	listener   CallbackEventListener
	stop       chan bool
	api        *Api
	GroupId    int
}

type BotLongPollResponse struct {
	Ts         json.Number     `json:"ts,Number"`
	Updates    []CallbackEvent `json:"updates"`
	Failed     int             `json:"failed"`
	MinVersion int             `json:"min_version"`
	MaxVersion int             `json:"nax_version"`
}

func GetBotLongPollServer(api *Api, logger *log.Logger, groupId int) *BotLongPollServer {
	return &BotLongPollServer{
		api:     api,
		stop:    make(chan bool, 1),
		Logger:  logger,
		GroupId: groupId,
		Ts:      0,
	}
}

func (s *BotLongPollServer) SetListener(listener CallbackEventListener) {
	s.listener = listener
}

func (s *BotLongPollServer) Listen(listener CallbackEventListener) error {
	s.SetListener(listener)
	return s.Start()
}

func (s *BotLongPollServer) getUpdates() (*BotLongPollResponse, error) {

	path := s.Server + "?act=a_check&key=" + s.Key + "&ts=" + strconv.Itoa(s.Ts) + "&wait=80"

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(req.Context(), 300*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	if s.httpClient == nil {
		s.httpClient = s.getHttpClient()
	}

	requestError := make(chan error, 1)
	requestResponse := make(chan *http.Response, 1)
	go func() {
		res, err := s.httpClient.Do(req)
		if err != nil {
			requestError <- err
		} else {
			requestResponse <- res
		}
	}()

	for {
		select {
		case <-s.stop:
			cancel()
			return nil, &LPCloseError{}
		case err := <-requestError:
			return nil, err
		case res := <-requestResponse:
			data, err := ioutil.ReadAll(res.Body)
			_ = res.Body.Close()
			if err != nil {
				return nil, &TransportError{
					string(s.Server),
					P{
						"ts":  strconv.Itoa(s.Ts),
						"key": s.Key,
					},
					data,
					TransportExternalData(res.Header),
					nil,
				}
			}
			response := BotLongPollResponse{}
			err = json.Unmarshal(data, &response)
			if err != nil {
				return nil, errors.New("BotLongPoll json decode error: " + err.Error() + " on string: " + string(data))
			}
			if response.Failed == 0 && (response.Ts == "" || response.Ts == "0") {
				return nil, errors.New("BotLongPoll bad response: " + string(data))
			}
			return &response, nil
		}
	}
}

func (s *BotLongPollServer) up(updateTs bool) error {
	data, err := s.api.run(
		"groups.getLongPollServer",
		P{
			"group_id": strconv.Itoa(s.GroupId),
		}, 0)
	if err != nil {
		return err
	}
	var e error
	if updateTs {
		s.Ts, e = data.GetInt("ts")
	}
	s.Key, e = data.GetString("key")
	s.Server, e = data.GetString("server")
	if e != nil {
		return e
	}
	return nil
}

func (s *BotLongPollServer) Stop() {
	s.stop <- true
}

func (s *BotLongPollServer) Start() error {
	err := s.up(s.Ts == 0)
	if err != nil {
		return err
	}
	for {
		updates, err := s.getUpdates()

		if err != nil {
			if err.Error() == LongPollClosed {
				return nil
			}
			s.onError(err)
		} else if updates.Failed == 0 {
			newTs, err := strconv.Atoi(string(updates.Ts))
			if err != nil {
				return errors.New("Cant get new ts from updates, ints not INT: " + err.Error())
			}
			if newTs > s.Ts {
				s.Ts = newTs
			}
			s.onUpdate(updates.Updates)
		} else if updates.Failed == 1 {
			newTs, err := strconv.Atoi(string(updates.Ts))
			if err != nil {
				return errors.New("Cant get new ts from updates, ints not INT: " + err.Error())
			}
			s.onErrorS(fmt.Sprintf("failed: 1, new ts %d old ts %d", newTs, s.Ts))
			s.Ts = newTs
		} else if updates.Failed == 2 {
			s.onErrorS("failed: 2")
			err := s.up(false)
			if err != nil {
				return err
			}
		} else if updates.Failed == 3 {
			s.onErrorS("failed: 3")
			err := s.up(true)
			if err != nil {
				return err
			}
		} else if updates.Failed == 4 {
			return errors.New("Invalid version min:" + strconv.Itoa(updates.MinVersion) + " max:" + strconv.Itoa(updates.MaxVersion))
		}

		select {
		case <-s.stop:
			return nil
		default:

		}
	}
}

func (s *BotLongPollServer) onError(err error) {
	s.onErrorS(printError(err))
}

func (s *BotLongPollServer) onErrorS(msg string) {
	if s.Logger != nil {
		s.Logger.Println(s.GetName() + ": " + msg)
	}
}

func (s *BotLongPollServer) onUpdate(updates []CallbackEvent) {
	for _, upd := range updates {
		s.listener(&upd)
	}
}

func (s *BotLongPollServer) getHttpClient() *http.Client {
	c := http.DefaultClient
	c.Timeout = 300 * time.Second
	return c
}

func (s *BotLongPollServer) GetName() string {
	return "BotLongPoll"
}
