package main

import (
	"io/ioutil"

	"github.com/guozijn/nuntius/provider/dingtalk"
	"github.com/guozijn/nuntius/provider/telegram"
	"github.com/prometheus/alertmanager/template"
	"gopkg.in/yaml.v2"
)

type ReceiverConf struct {
	Name     string
	Provider string
	To       []string
	From     string
	Text     string
}

var providerConfig struct {
	Providers struct {
		Telegram telegram.TelegramConfig
		DingTalk dingtalk.Config
	}

	Receivers []ReceiverConf
	Templates []string
}
var tmpl *template.Template

// LoadConfig loads the specified YAML configuration file.
func LoadConfig(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &providerConfig)
	if err != nil {
		return err
	}

	tmpl, err = template.FromGlobs(providerConfig.Templates...)
	return err
}
