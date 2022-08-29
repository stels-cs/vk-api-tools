// Пакет VkApi предоствляет инструменты для упрощенного взамодействия с API ВКонтакте.
//
//   - Обычные зпросы к API
//   - LongPoll
//   - Группировка запросов в execute
//   - Создание очереди из запросов с огрпниченимем кол-ва звпросов в секунду
//   - StreamingAPI
//
// # Пример использования
//
// Получение списка пользователей
//
//	users := make([]struct{
//		Id        int    `json:"id"`
//		FirstName string `json:"first_name"`
//		LastName  string `json:"last_name"`
//	}, 0)
//	err := api.Exec("users.get", VkApi.P{"user_ids": "2050,andrew"}, &users)
//	if err != nil {
//		panic(err)
//	}
//	for _, u := range users {
//		println(u.FirstName + " " + u.LastName + " #" + strconv.Itoa(u.Id))
//	}
//	// Output:
//	// Катя Лебедева #2050
//	// Андрей Рогозов #6492
//
// # Пример без использования структуры
//
// Получение города пользователя
//
//	res, err := VkApi.Call("users.get", VkApi.P{"user_ids": "2050,avk", "fields": "city"})
//	if err != nil {
//		panic(err)
//	}
//
//	print(res.QStringDef("0.first_name", "") + " – ")
//	println(res.QStringDef("0.city.title", ""))
//
//	print(res.QStringDef("1.first_name", "") + " – ")
//	println(res.QStringDef("1.city.title", ""))
//	//Катя – Санкт-Петербург
//	//Александр – Москва
//
// Пример с LongPoll
//
// Типа бот
//
//	type Bot struct {
//		token string
//	}
//	func (b *Bot) NewMessage(event VkApi.MessageEvent) {
//		if event.IsOutMessage() {
//			return
//		}
//		u := make([]struct {
//			FirstName string `json:"first_name"`
//		}, 0)
//		err := VkApi.Exec("users.get", VkApi.P{"user_id": strconv.Itoa(event.PeerId)}, &u)
//		if err != nil {
//			println(err.Error())
//			return
//		}
//		msg := "Привет, " + u[0].FirstName
//		_, err = VkApi.Call("messages.send", VkApi.P{
//			"peer_id":      strconv.Itoa(event.PeerId),
//			"message":      msg,
//			"access_token": b.token,
//		})
//		if err != nil {
//			println(err.Error())
//			return
//		}
//	}
//	func (b *Bot) EditMessage(event VkApi.MessageEvent) {
//		msg := "Не редактируйте сообщения прлиз"
//		_, err := VkApi.Call("messages.send", VkApi.P{
//			"peer_id":      strconv.Itoa(event.PeerId),
//			"message":      msg,
//			"access_token": b.token,
//		})
//		if err != nil {
//			println(err.Error())
//			return
//		}
//	}
//	func main() {
//		token := "529d99ca66f0....da3b0e717e7df"
//		bot := &Bot{token}
//		logger := log.New(os.Stdout, "LongPoll", log.Lshortfile)
//		api := VkApi.CreateApi(token, "5.71", VkApi.GetHttpTransport(), 30)
//		lp := VkApi.GetLongPollServer(api, logger)
//		lp.SetListener(bot)
//		go lp.Start()
//		println("Press Ctl+C for quit")
//		c := make(chan os.Signal, 1)
//		signal.Notify(c, os.Interrupt)
//		<-c
//		lp.Stop()
//	}
//
// # Пример создания очереди запросов
//
// Запросы автоматически упаковываются в execute вызовы, в случае если поток запросов больше чем максимальный RPS очереди.
// В данном примере максимум 1 вызов в секунду.  В примере делается 4 запроса, но к API будет сделано всего два execute запроса
//
//	token := "529d99c....c88bf19da3b0e717e7df"
//	v := "5.71"
//	api := VkApi.CreateApi(token, v, VkApi.GetHttpTransport(), 30)
//	requestsPerSecond := 1
//	rq := VkApi.GetRequestQueue(api, requestsPerSecond)
//	go rq.Start()
//	defer rq.Stop()
//
//	l1 := make(chan int)
//	l2 := make(chan int)
//
//	var diff1 int64
//	var diff3 int64
//	var r0, r1, r2, r3 VkApi.RequestResult
//
//	go func() {
//		ts1 := time.Now().UnixNano()
//		c0 := rq.Call(VkApi.CreateMethod("users.get", VkApi.P{"user_ids": "1"}))
//		r0 = <-c0
//		diff1 = time.Now().UnixNano() - ts1
//		l1 <- 1
//	}()
//
//	go func() {
//		ts3 := time.Now().UnixNano()
//		c1 := rq.Call(VkApi.CreateMethod("users.get", VkApi.P{"user_ids": "2050"}))
//		c2 := rq.Call(VkApi.CreateMethod("users.get", VkApi.P{"user_ids": "avk"}))
//		c3 := rq.Call(VkApi.CreateMethod("users.get", VkApi.P{"user_ids": "andrew"}))
//
//		r1 = <-c1
//		r2 = <-c2
//		r3 = <-c3
//		diff3 = time.Now().UnixNano() - ts3
//		l2 <- 1
//	}()
//
//	<-l1
//	<-l2
//	println("One request: " + strconv.Itoa(int(diff1/int64(time.Nanosecond))) + "ns")
//	println(r0.Res.QStringDef("0.first_name", "") + "\n")
//	println("Three request: " + strconv.Itoa(int(diff3/int64(time.Nanosecond))) + "ns")
//	println(r1.Res.QStringDef("0.first_name", ""))
//	println(r2.Res.QStringDef("0.first_name", ""))
//	println(r3.Res.QStringDef("0.first_name", ""))
//	// Output
//	// One request: 3483161000ns
//	// Павел
//	//
//	// Three request: 3483171000ns
//	// Катя
//	// Александр
//	// Андрей
//
// # Пример StreamingAPI
//
// для работы с websocket используется github.com/gorilla/websocket
//
//	 Необходимо реализовать интерфейсы для работы с WebSockets
//
//	 type WsConnection struct {
//	 	c *websocket.Conn
//	 }
//
//	 func (conn *WsConnection) Close() error {
//	 	if conn.c == nil {
//	 		return nil
//			}
//			return conn.c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
//	 }
//
//	 func (conn *WsConnection) ReadMessage() ([]byte, error) {
//	 	_, message, err := conn.c.ReadMessage()
//	 	return message, err
//	 }
//
//	 func GetConnection(url string) (VkApi.ConnectionInterface, error) {
//	 	conn, response, err := websocket.DefaultDialer.Dial(url, nil)
//	 	if err != nil {
//	 		body, e := ioutil.ReadAll(response.Body)
//	 		if e == nil {
//	 			response.Body.Close()
//	 			return nil, errors.New(string(body))
//	 		}
//	 		return nil, err
//	 	}
//
//	 	return &WsConnection{
//	 	c: conn,
//	 	}, nil
//	 }
//
//	 func IsCloseError(err error) bool {
//	 	if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
//	 		return true
//	 	} else {
//	 		return false
//	 	}
//		}
//
//
//	 // Понадобиться сервиный ключ приложения
//	 token := "60fb7fa360fb7fa360fb7fa34b60a5f1e7660fb60fb7fa33a685687cb8b239d9a73303c"
//
//	 // Создаем объект VkApi.StreamingClient
//	 streamingApi, err := VkApi.CreateStreamingClient(token)
//	 if err != nil {
//	 	panic(err)
//	 }
//
//	 // Получаем все правила которые есть сейчас для этого ключа
//	 rules, err := streamingApi.GetRules()
//	 if err != nil {
//	 	panic(err)
//	 }
//
//	 // Удаляем все правила которые были
//	 for _, rule := range rules {
//	 	println("Rule:", rule.Tag, "Value:", rule.Value)
//	 	err := streamingApi.DeleteRule(rule.Tag)
//	 	if err != nil {
//	 		panic(err)
//	 	}
//	 }
//
//	 // Создаем свои правила
//	 rule1 := VkApi.StreamingRule{
//	 	Value: "путин",
//	 	Tag:   "putin",
//	 }
//
//	 rule2 := VkApi.StreamingRule{
//	 	Value: "кандидиат",
//	 	Tag:   "candidate",
//	 }
//
//	 rule3 := VkApi.StreamingRule{
//	 	Value: "биткойн",
//	 	Tag:   "bitcoin",
//	 }
//
//	 err = streamingApi.AddRule(rule1)
//	 if err != nil {
//	 	panic(err)
//	 }
//	 err = streamingApi.AddRule(rule2)
//	 if err != nil {
//	 	panic(err)
//	 }
//	 err = streamingApi.AddRule(rule3)
//	 if err != nil {
//	 	panic(err)
//	 }
//
//	 println("Rules created")
//	 go func () {
//	 	time.Sleep(30 * time.Second)
//	 	streamingApi.Stop() // Вот так можно остановить прослушивание сообщений
//	 	println("Stop streaming by timeout")
//	 } ()
//
//	 println("Start lister event")
//	 err = streamingApi.Start(GetConnection, IsCloseError, func(code int, event *VkApi.StreamingEvent, message *VkApi.StreamingServiceMessage) {
//	 	if event != nil {
//	 		println(event.Text, event.EventUrl)
//	 	}
//	 } )
//
//	 if err != nil {
//	 	panic(err)
//	 }
package VkApi

import (
	"bytes"
	"encoding/json"
	"time"
)

type P map[string]string

func (p *P) toString() string {
	str := ""
	for k, v := range *p {
		if str != "" {
			str += " "
		}
		str += k + "=" + v
	}
	return str
}

func (p *P) toJson() ([]byte, error) {
	return json.Marshal(p)
}

func (p *P) JSON() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(p)
	return buffer.Bytes(), err
}

func (p *P) toMap() map[string][]string {
	out := map[string][]string{}
	for k, v := range *p {
		out[k] = []string{v}
	}
	return out
}

type Api struct {
	token         string
	transport     Transport
	Version       string
	maxRetryCount int
}

// Создает новый объект api
//
//	token := "API TOKEN"
//	version := "5.71"
//	tr := VkApi
//	api := VkApi.CreateApi("31223dbcda...", "5.71", tr, 30)
func CreateApi(t, v string, transport Transport, maxRetryCount int) *Api {
	api := &Api{
		t,
		transport,
		v,
		maxRetryCount,
	}
	return api
}

// TODO: написать тест когда params == nil
func (api *Api) run(method string, params P, retryCount int) (Response, error) {
	if params == nil {
		params = P{}
	}
	response := Response{}
	b := json.RawMessage(``)
	response.Response = &b

	if retryCount == 0 {
		retryCount = api.maxRetryCount
	}

	if len(api.token) > 0 {
		params["access_token"] = api.token
	}
	if _, has := params["v"]; has == false {
		params["v"] = api.Version
	}

	data, err, external := api.transport.call(method, params)
	if err != nil {
		if retryCount > 0 {
			return api.retryCall(method, params, retryCount)
		} else {
			return response, &TransportError{
				method,
				params,
				data,
				external,
				err,
			}
		}
	}

	err = json.Unmarshal(data, &response)
	if err != nil {
		if retryCount > 0 {
			return api.retryCall(method, params, retryCount)
		} else {
			return response, &TransportError{
				method,
				params,
				data,
				external,
				nil,
			}
		}
	}

	if response.canRetry() && retryCount > 0 {
		return api.retryCall(method, params, retryCount)
	}

	response.Error.CallMethod = &method
	response.Error.CallParams = &params

	if response.success() {
		return response, nil
	} else {
		return response, &response.Error
	}
}

func (api *Api) retryCall(method string, params P, retryCount int) (Response, error) {
	time.Sleep(time.Second)
	res, err := api.run(method, params, retryCount-1)
	return res, err
}

// Выполняет запрос к API ВКонтакте
//
//	users := make([]struct{
//		Id        int    `json:"id"`
//		FirstName string `json:"first_name"`
//		LastName  string `json:"last_name"`
//	}, 0)
//	res, err := api.Call("users.get", VkApi.P{"user_ids": "2050,andrew"})
//	name := res.QStringDef("0.first_name") // name == "Катя"
func (api *Api) Call(name string, params P) (Response, error) {
	return api.run(name, params, 0)
}

// Выполняет запрос к API ВКонтакте, записывает результат в структуру s
//
//	users := make([]struct{
//		Id        int    `json:"id"`
//		FirstName string `json:"first_name"`
//		LastName  string `json:"last_name"`
//	}, 0)
//	err := api.Exec("users.get", VkApi.P{"user_ids": "2050,andrew"}, &users)
func (api *Api) Exec(name string, params P, s interface{}) error {
	res, err := api.Call(name, params)
	if err != nil {
		return err
	}
	err = res.Unmarshal(s)
	if err != nil {
		return err
	}
	return nil
}
