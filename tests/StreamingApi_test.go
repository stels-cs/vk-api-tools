package VkApiTest

import (
	"github.com/stels-cs/vk-api-tools"
	"testing"
	"time"
)

func TestStreamingApi(t *testing.T) {
	if token := getToken(); token != "" {
		s, err := VkApi.CreateStreamingClient(token)
		if err != nil {
			t.Error(err)
		}

		rules, err := s.GetRules()
		if err != nil {
			t.Error(err)
		}

		for _, rule := range rules {
			err := s.DeleteRule(rule.Tag)
			if err != nil {
				t.Error(err)
			}
		}

		rule1 := VkApi.StreamingRule{
			Value: "коты",
			Tag:   "cats",
		}

		rule2 := VkApi.StreamingRule{
			Value: "кек",
			Tag:   "kek",
		}

		err = s.AddRule(rule1)
		if err != nil {
			t.Error(err)
		}
		err = s.AddRule(rule2)
		if err != nil {
			t.Error(err)
		}

		rules, err = s.GetRules()
		if len(rules) != 2 {
			t.Error("Expected 2 rulse, got", len(rules))
		}
		for _, rule := range rules {
			if rule.Tag != rule1.Tag && rule.Tag != rule2.Tag {
				t.Error("Expected created tag rulse, got", rule.Tag)
			}
		}

		err = s.ClearAllRules()

		rules, err = s.GetRules()

		if len(rules) != 0 {
			t.Error("Expected no rulse, got", len(rules))
		}
	} else {
		skip(t)
	}
}

func TestStreamingReading(t *testing.T) {
	if token := getToken(); token != "" {
		s, err := VkApi.CreateStreamingClient(token)
		if err != nil {
			t.Error(err)
		}

		err = s.ClearAllRules()
		if err != nil {
			t.Error(err)
		}

		s.AddRule(VkApi.StreamingRule{"путин", "cat"})
		s.AddRule(VkApi.StreamingRule{"кандидат", "what"})

		c := make(chan bool, 1)

		go func() {
			err = s.Start(func(code int, event *VkApi.StreamingEvent, msg *VkApi.StreamingServiceMessage) {
				if code == 100 {
					s.Stop()
				}
			})

			if err != nil {
				t.Error(err)
			}
			c <- false
		}()

		select {
		case <-c:
			return
		case <-time.After(10 * time.Second):
			t.Skip("No message received for 10 second")
		}

	} else {
		skip(t)
	}
}
