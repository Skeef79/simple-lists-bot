package bot

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	user "github.com/skeef79/simple-lists-bot/user"
)

type bot struct {
	BotApi         *tgbotapi.BotAPI
	commands       map[commandKey]commandEntity
	users          map[int64]user.User
	userStates     map[int64]int64
	userStatesInfo map[int64]uint64 //actually this used only to store ID of a current list
	callbacks      map[CallbackType]CallbackFn
}

func (b *bot) apiRequest(c tgbotapi.Chattable) error {
	_, err := b.BotApi.Request(c)
	return err
}

func NewBot() (*bot, error) {
	api, aErr := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if aErr != nil {
		return nil, aErr
	}

	api.Debug = true

	b := &bot{
		BotApi:         api,
		commands:       make(map[commandKey]commandEntity),
		users:          make(map[int64]user.User),
		userStatesInfo: make(map[int64]uint64),
		userStates:     make(map[int64]int64),
	}

	if err := b.initCommands(); err != nil {
		return nil, err
	}

	b.InitCallbacks()

	log.Default().Print("bot created")
	return b, nil
}

func (b *bot) Run() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	for upd := range b.BotApi.GetUpdatesChan(updateConfig) {
		if upd.Message != nil {
			if upd.Message.IsCommand() {
				key := upd.Message.Command()
				if cmd, ok := b.commands[commandKey(key)]; ok {
					go cmd.Action(upd)
				} else {
					fmt.Errorf("Command handler fot %s not found", key)
				}
				continue
			}

			go b.HandleMessage(upd)
		}

		if upd.CallbackQuery != nil {
			data := upd.CallbackData()
			entity := unmarshallCb(data)

			callback := tgbotapi.NewCallback(upd.CallbackQuery.ID, "")
			b.apiRequest(callback)

			b.callbacks[entity.CbType](upd, entity)
		}
	}
}
