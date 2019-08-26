package main

import (
	"io/ioutil"
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

var config struct {
	Providers struct {
		MessageBird messagebird.MessageBirdConfig
		Nexmo       nexmo.NexmoConfig
		Twilio      twilio.TwilioConfig
		Infobip     infobip.InfobipConfig
		Exotel      exotel.ExotelConfig
		CM          cm.CMConfig
		Telegram    telegram.TelegramConfig
		Turbosms    turbosms.TurbosmsConfig
		OTC         otc.OTCConfig
		MediaBurst  mediaburst.MediaBurstConfig
		FreeMobile  freemobile.Config
		AspSms      aspsms.Config
		Sipgate     sipgate.Config
		Pushbullet  pushbullet.Config
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

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return err
	}

	tmpl, err = template.FromGlobs(config.Templates...)
	return err
}
