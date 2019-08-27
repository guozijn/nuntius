package dingtalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	
	"github.com/guozijn/nuntius"
)

var dingtalkHTTPClient = &http.Client{Timeout: time.Second * 20}

type Config struct {
	Token string `yaml:"token"`
	URL   string `yaml:"url"`
	MsgType string `yaml:"msgtype"`
	AtMobiles []string `yaml:"atmobiles"`
	IsAtAll bool `yaml:"isatall"`
}

type DingTalk struct {
	Config
}

type DingTalkPayload struct {
	Msgtype string `json:"msgtype"`
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
	At struct {
		AtMobiles []string `json:"atMobiles"`
		IsAtAll   bool     `json:"isAtAll"`
	} `json:"at"`
}

func NewDingTalk(config Config) *DingTalk {
	if config.URL == "" {
		config.URL = "https://oapi.dingtalk.com/robot/send"
	}
	if config.AtMobiles == nil {
		config.AtMobiles = []string{}
	}
	return &DingTalk{config}
}

// Send sends SMS to user registered in configuration
func (c *DingTalk) Send(message nuntius.Message) error {
	payload := &DingTalkPayload{}
	payload.Msgtype = c.MsgType
	payload.Text.Content = message.Text
	payload.At.AtMobiles = c.AtMobiles
	payload.At.IsAtAll = c.IsAtAll

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", c.URL + "?access_token=" + c.Token, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Nuntius")

	response, err := dingtalkHTTPClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var body []byte
	response.Body.Read(body)
	if response.StatusCode == http.StatusOK && err == nil {
		return nil
	}

	return fmt.Errorf("Failed sending message. statusCode: %d", response.StatusCode)
}
