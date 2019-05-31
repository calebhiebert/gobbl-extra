package tel

import (
	"github.com/calebhiebert/gobbl"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Integration is the basis for the telegram integration
type Integration struct {
	API    *tgbotapi.BotAPI
	gobblr *gbl.Bot
}

// New creates a new telegram integration with a given token
func New(telegramToken string, gobblr *gbl.Bot) (*Integration, error) {
	integration := &Integration{}

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		return nil, err
	}

	integration.API = bot
	integration.gobblr = gobblr

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updatesChannel, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		return nil, err
	}

	go integration.handleUpdates(updatesChannel)

	return integration, nil
}

func (i *Integration) handleUpdates(updateChannel tgbotapi.UpdatesChannel) {
	for update := range updateChannel {
		i.handleUpdate(&update)
	}
}

func (i *Integration) handleUpdate(update *tgbotapi.Update) {
	inputContext := gbl.InputContext{
		RawRequest:  update,
		Integration: i,
	}

	go i.gobblr.Execute(&inputContext)
}
