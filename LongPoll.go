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

type LongPollListener interface {
	NewMessage(msg MessageEvent)
	EditMessage(msg MessageEvent)
}

type LongPollDefaultListener struct{}

func (l *LongPollDefaultListener) NewMessage(msg MessageEvent)  {}
func (l *LongPollDefaultListener) EditMessage(msg MessageEvent) {}

type LPCloseError struct {
}

func (e *LPCloseError) Error() string {
	return "LONG_POLL_CLOSED"
}

type LongPollServer struct {
	Key        string
	Server     string
	Ts         int
	httpClient *http.Client
	Logger     *log.Logger
	intSize    int
	listener   LongPollListener
	stop       chan bool
	api        *Api
}

type LongPollServerResponse struct {
	Key    string `json:"key"`
	Server string `json:"server"`
	Ts     int    `json:"ts"`
}

type Update []interface{}

type LongPollResponse struct {
	Ts         int      `json:"ts"`
	Updates    []Update `json:"updates"`
	Failed     int      `json:"failed"`
	MinVersion int      `json:"min_version"`
	MaxVersion int      `json:"nax_version"`
}

func GetLongPollServer(api *Api, logger *log.Logger) *LongPollServer {
	return &LongPollServer{
		api:    api,
		stop:   make(chan bool, 1),
		Logger: logger,
	}
}

func (s *LongPollServer) SetListener(listener LongPollListener) {
	s.listener = listener
}

func (s *LongPollServer) getUpdates() (*LongPollResponse, error) {

	path := "https://" + s.Server + "?act=a_check&key=" + s.Key + "&ts=" + strconv.Itoa(s.Ts) + "&wait=25&mode=42&version=2"

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
			res.Body.Close()
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
			response := LongPollResponse{}
			err = json.Unmarshal(data, &response)
			if err != nil {
				return nil, err
			}
			if response.Failed == 0 && response.Ts == 0 {
				return nil, errors.New("Long poll bad response: " + string(data))
			}
			return &response, nil
		}
	}
}

func (s *LongPollServer) up(updateTs bool) error {
	data, err := s.api.run(
		"messages.getLongPollServer",
		P{
			"lp_version": "2",
			"need_pts":   "1",
		}, 30)
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

func (s *LongPollServer) Stop() {
	s.stop <- true
}

func (s *LongPollServer) Start() error {
	err := s.up(true)
	if err != nil {
		return err
	}
	for {
		updates, err := s.getUpdates()

		if err != nil {
			if err.Error() == "LONG_POLL_CLOSED" {
				return nil
			}
			if s.Logger != nil {
				s.Logger.Println(printError(err))
			}
		} else if updates.Failed == 0 {
			if updates.Ts > s.Ts {
				s.Ts = updates.Ts
			}
			s.onUpdate(updates.Updates)
		} else if updates.Failed == 1 {
			if s.Logger != nil {
				s.Logger.Printf("Long poll failed 1, new ts %d old ts %d", updates.Ts, s.Ts)
			}
			s.Ts = updates.Ts
		} else if updates.Failed == 2 {
			if s.Logger != nil {
				s.Logger.Println("Long poll failed 2")
			}
			err := s.up(false)
			if err != nil {
				return err
			}
		} else if updates.Failed == 3 {
			if s.Logger != nil {
				s.Logger.Println("Long poll failed 3")
			}
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

func (s *LongPollServer) onError(err error) {
	if s.Logger != nil {
		s.Logger.Println(printError(err))
	}
}

func (s *LongPollServer) onUpdate(updates []Update) {
	for _, upd := range updates {
		if len(upd) > 0 {
			if typeId, ok := upd[0].(float64); ok {
				s.onEvent(int(typeId), upd)
			} else {
				s.onBadEvent(0, upd, errors.New("Cant detect type id\n"))
			}
		} else {
			s.onBadEvent(0, upd, errors.New("Data size is zero\n"))
		}
	}
}

func (s *LongPollServer) onBadEvent(typeId int, data Update, err error) {
	s.onError(errors.New(fmt.Sprintf("Bad message %s type: %d data %+v", err.Error(), typeId, data)))
}

func (s *LongPollServer) onEvent(typeId int, data Update) {
	if typeId == LPNewMessage || typeId == LPEditMessage {
		message := MessageEvent{}
		err := message.Fill(data)
		if err != nil {
			s.onBadEvent(typeId, data, err)
		} else if s.listener != nil {
			if typeId == LPEditMessage {
				message.IsEditMessage = true
				if s.listener != nil {
					s.listener.EditMessage(message)
				}
			} else {
				if s.listener != nil {
					s.listener.NewMessage(message)
				}
			}
		}
	}
}
func (s *LongPollServer) getHttpClient() *http.Client {
	c := http.DefaultClient
	c.Timeout = 300 * time.Second
	return c
}

func (s *LongPollServer) GetName() string {
	return "Long poll server"
}
