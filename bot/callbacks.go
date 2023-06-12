package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skeef79/simple-lists-bot/user"
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
	userID := upd.CallbackQuery.Message.Chat.ID
	listID := entity.ListID
	list, _ := b.users[upd.CallbackQuery.Message.Chat.ID].GetListById(listID)

	reply := tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID, createListMessage(list))
	b.userStatesInfo[upd.CallbackQuery.Message.Chat.ID] = entity.ListID

	reply.ReplyMarkup = getListKeyboard()
	reply.ParseMode = "html"
	if err := b.apiRequest(reply); err != nil {
		log.Fatalf("failed to send api request inside lists callback")
	}
	b.userStates[userID] = user.ListEditState
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

	cbType, _ := strconv.ParseInt(data[0], 10, 64)
	listID, _ := strconv.ParseInt(data[1], 10, 64)
	return CallbackEntity{
		CbType: CallbackType(cbType),
		ListID: uint64(listID),
	}
}
