package Vk

import (
	"strconv"
	"encoding/json"
)

type Message struct {
	api *Api
}

func (m *Message) GetLongPollServer(needPts bool, lpVersion int) (LongPollServerResponse, error) {
	server := LongPollServerResponse{}
	response, err := m.api.run(
		"messages.getLongPollServer",
		Params{
			"lp_version": strconv.Itoa(lpVersion),
			"need_pts":   iif(needPts, "1", "0"),
		}, 0)
	if err != nil {
		return server, err
	}

	err = json.Unmarshal(*response.Response, &server)
	if err != nil {
		return server, err
	} else {
		return server, nil
	}
}

func (m *Message) SendTextRequest(peerId int, text string) ApiMethod {
	return GetApiMethod("messages.send", Params{
		"peer_id": strconv.Itoa(peerId),
		"message": text,
	})
}
