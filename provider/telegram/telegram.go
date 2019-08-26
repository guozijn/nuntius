package telegram

import (
	"strconv"

	"github.com/guozijn/nuntius"
	"gopkg.in/telegram-bot-api.v4"
)

type TelegramConfig struct {
	Token string `yaml:"token"`
}

type Telegram struct {
	bot *tgbotapi.BotAPI
}

func NewTelegram(config TelegramConfig) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, err
	}

	return &Telegram{
		bot: bot,
	}, nil
}

func (tg *Telegram) Send(message nuntius.Message) error {
	for _, sChatID := range message.To {
		chatID, err := strconv.ParseInt(sChatID, 10, 64)
		if err != nil {
			return err
		}
		msg := tgbotapi.NewMessage(chatID, message.Text)
		_, err = tg.bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
