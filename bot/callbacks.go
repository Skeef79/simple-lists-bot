package bot

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"skeef79.com/simple-tg-bot/user"
)

type CallbackType int

const (
	Lists CallbackType = iota
	List
)

type CallbackEntity struct {
	CbType     CallbackType
	ListString string
	ListID     uint64
	ListName   string
}

type CallbackFn func(upd tgbotapi.Update, entity CallbackEntity)

func (b *bot) InitCallbacks() {
	b.callbacks = map[CallbackType]CallbackFn{
		Lists: b.ListsCallback,
		List:  b.ListCallback,
	}
}

func getListKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(AddItemMessage),
			tgbotapi.NewKeyboardButton(DeleteItemMessage),
			tgbotapi.NewKeyboardButton(BackMessage),
		),
	)
}

func (b *bot) ListsCallback(upd tgbotapi.Update, entity CallbackEntity) {

	reply := tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID, entity.ListString)
	b.userStatesInfo[upd.CallbackQuery.Message.Chat.ID] = entity.ListID

	reply.ReplyMarkup = getListKeyboard()
	reply.ParseMode = "html"
	if err := b.apiRequest(reply); err != nil {
		fmt.Errorf("failed to send api request inside lists callback")
	}
	b.userStates[upd.CallbackQuery.Message.Chat.ID] = user.ListEditState
}

func (b *bot) ListCallback(upd tgbotapi.Update, entity CallbackEntity) {

}

func marshallCb(ce CallbackEntity) string {
	return fmt.Sprintf(
		"%d;%s;%d;%s",
		ce.CbType,
		ce.ListString,
		ce.ListID,
		ce.ListName,
	)
}

func unmarshallCb(ce string) CallbackEntity {
	data := strings.Split(ce, ";")

	cbType, _ := strconv.Atoi(data[0])
	listID, _ := strconv.Atoi(data[2])
	return CallbackEntity{
		CbType:     CallbackType(cbType),
		ListString: data[1],
		ListID:     uint64(listID),
		ListName:   data[3],
	}
}
