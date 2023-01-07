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
)

type CallbackEntity struct {
	CbType CallbackType
	ListID uint64
}

type CallbackFn func(upd tgbotapi.Update, entity CallbackEntity)

func (b *bot) InitCallbacks() {
	b.callbacks = map[CallbackType]CallbackFn{
		Lists: b.ListsCallback,
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

	listID := entity.ListID
	list, _ := b.users[upd.CallbackQuery.Message.Chat.ID].GetListById(listID)

	reply := tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID, createListMessage(list))
	b.userStatesInfo[upd.CallbackQuery.Message.Chat.ID] = entity.ListID

	reply.ReplyMarkup = getListKeyboard()
	reply.ParseMode = "html"
	if err := b.apiRequest(reply); err != nil {
		fmt.Errorf("failed to send api request inside lists callback")
	}
	b.userStates[upd.CallbackQuery.Message.Chat.ID] = user.ListEditState
}

func marshallCb(ce CallbackEntity) string {
	return fmt.Sprintf(
		"%d;%d",
		ce.CbType,
		ce.ListID,
	)
}

func unmarshallCb(ce string) CallbackEntity {
	data := strings.Split(ce, ";")

	cbType, _ := strconv.Atoi(data[0])
	listID, _ := strconv.Atoi(data[1])
	return CallbackEntity{
		CbType: CallbackType(cbType),
		ListID: uint64(listID),
	}
}
