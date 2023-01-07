package bot

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skeef79/simple-lists-bot/storage"
	"github.com/skeef79/simple-lists-bot/user"
)

func (b *bot) StartCmd(upd tgbotapi.Update) {
	name := upd.Message.From.UserName
	if name == "" {
		name = fmt.Sprintf("%s %s", upd.Message.From.FirstName, upd.Message.From.LastName)
	}

	message := `
Welcome to the <b>simple-lists-bot</b>, %s!
This bot helps you to manage your lists of things.
Use /help to see available commands	
`

	reply := tgbotapi.NewMessage(upd.Message.Chat.ID, fmt.Sprintf(message, name))
	reply.ParseMode = "html"

	if err := b.apiRequest(reply); err != nil {
		fmt.Errorf("failed to send start message")
	}

	b.userStates[upd.Message.Chat.ID] = user.EmptyState
}

func (b *bot) HelpCmd(upd tgbotapi.Update) {
	message := `
This bot currently supports the following commands:
• /lists -- show all existing lists
• /add_list -- create new list
• /remove_list -- delete an existing list
`

	reply := tgbotapi.NewMessage(upd.Message.Chat.ID, message)
	if err := b.apiRequest(reply); err != nil {
		fmt.Errorf("failed to send help message")
	}

	b.userStates[upd.Message.Chat.ID] = user.EmptyState
}

func (b *bot) getListsKeyboard(ID int64) tgbotapi.InlineKeyboardMarkup {
	_, ok := b.users[ID]
	if !ok {
		b.users[ID] = user.User{
			Storage: storage.NewInMemStorage(fmt.Sprintf("%d", ID)),
		}
	}

	lists, _ := b.users[ID].GetAllLists()

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(lists))
	for _, list := range lists {
		button := tgbotapi.NewInlineKeyboardButtonData(list.Name, marshallCb(CallbackEntity{
			CbType: Lists,
			ListID: list.ID,
		}))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func (b *bot) ShowListsCmd(upd tgbotapi.Update) {
	message := "<b>Your lists</b>:\n"
	keyboard := b.getListsKeyboard(upd.Message.Chat.ID)
	if len(keyboard.InlineKeyboard) == 0 {
		message = "You don't have any lists yet\n"
	}

	reply := tgbotapi.NewMessage(upd.Message.Chat.ID, message)

	if len(keyboard.InlineKeyboard) != 0 {
		reply.ReplyMarkup = keyboard
	}

	reply.ParseMode = "html"

	if err := b.apiRequest(reply); err != nil {
		fmt.Errorf("failed to send show lists message")
	}

	b.userStates[upd.Message.Chat.ID] = user.ShowListsState
}

func (b *bot) AddListCmd(upd tgbotapi.Update) {
	message := "Type the name of a list:\n"
	reply := tgbotapi.NewMessage(upd.Message.Chat.ID, message)
	b.userStates[upd.Message.Chat.ID] = user.CreateListState

	if err := b.apiRequest(reply); err != nil {
		fmt.Errorf("failed to send add list")
	}
}

func (b *bot) HandleMessage(upd tgbotapi.Update) {
	ID := upd.Message.Chat.ID
	state, ok := b.userStates[ID]
	if !ok || state == user.EmptyState {
		b.HelpCmd(upd)
		return
	}

	text := upd.Message.Text
	if state == user.CreateListState {
		_, ok := b.users[ID]
		if !ok {
			b.users[ID] = user.User{
				Storage: storage.NewInMemStorage(fmt.Sprintf("%d", ID)),
			}
		}
		b.users[ID].CreateList(text)
		b.userStates[ID] = user.EmptyState

		message := fmt.Sprintf("List '%s' was successfully created!", text)
		reply := tgbotapi.NewMessage(ID, message)

		if err := b.apiRequest(reply); err != nil {
			fmt.Errorf("failed to send create list")
		}

		return
	}

	if state == user.ShowListsState {
		message := "You can choose a list to see/add/delete items in it"
		reply := tgbotapi.NewMessage(ID, message)

		if err := b.apiRequest(reply); err != nil {
			fmt.Errorf("failed to send show list message")
		}
		return
	}

	if state == user.ListEditState {
		switch text {
		case AddItemMessage:
			message := "Type the name of an item"
			reply := tgbotapi.NewMessage(ID, message)
			b.userStates[ID] = user.AddItemState
			reply.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

			if err := b.apiRequest(reply); err != nil {
				fmt.Errorf("failed to send message")
			}
		case DeleteItemMessage:
			message := "Type the number of an item to delete"
			reply := tgbotapi.NewMessage(ID, message)
			b.userStates[ID] = user.DeleteItemState
			reply.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

			if err := b.apiRequest(reply); err != nil {
				fmt.Errorf("failed to send message")
			}
		case BackMessage:
			reply := tgbotapi.NewMessage(ID, "Gettings back to the lists...")
			b.userStates[ID] = user.ShowListsState
			reply.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

			if err := b.apiRequest(reply); err != nil {
				fmt.Errorf("failed to send message")
			}

			b.ShowListsCmd(upd)

		default:
			message := "Choose one of the options, please"
			reply := tgbotapi.NewMessage(ID, message)
			reply.ReplyMarkup = getListKeyboard()
			reply.ParseMode = "html"
			if err := b.apiRequest(reply); err != nil {
				fmt.Errorf("failed to send message")
			}
		}
		return
	}

	//There is a limitaion of 64 bytes on button data!
	//So we should only send an ID of a list
	//Current solution is super shit

	if state == user.AddItemState {
		itemName := text
		listID := b.userStatesInfo[ID]
		b.users[ID].AddItemByListId(listID, itemName)

		list, _ := b.users[ID].GetListById(listID)

		reply := tgbotapi.NewMessage(ID, createListMessage(list))
		reply.ReplyMarkup = getListKeyboard()
		reply.ParseMode = "html"
		if err := b.apiRequest(reply); err != nil {
			fmt.Errorf("failed to send api request inside lists callback")
		}
		b.userStates[ID] = user.ListEditState
	}

	if state == user.DeleteItemState {
		listID := b.userStatesInfo[ID]
		list, _ := b.users[ID].GetListById(listID)

		if len(list.Items) == 0 {
			message := "Can't delete item when list is empty"
			reply := tgbotapi.NewMessage(ID, message)
			if err := b.apiRequest(reply); err != nil {
				fmt.Errorf("failed to send api request to delete an item")
			}

			reply = tgbotapi.NewMessage(ID, createListMessage(list))
			reply.ReplyMarkup = getListKeyboard()
			reply.ParseMode = "html"
			if err := b.apiRequest(reply); err != nil {
				fmt.Errorf("failed to send api request inside lists callback")
			}
			b.userStates[ID] = user.ListEditState
			return
		}

		if itemIndex, err := strconv.Atoi(text); err == nil && itemIndex >= 1 && itemIndex <= len(list.Items) {
			err := b.users[ID].DeleteItemByIndex(listID, itemIndex-1)
			if err != nil {
				fmt.Errorf("%s", err)
			}

			reply := tgbotapi.NewMessage(ID, createListMessage(list))
			reply.ReplyMarkup = getListKeyboard()
			reply.ParseMode = "html"
			if err := b.apiRequest(reply); err != nil {
				fmt.Errorf("failed to send api request inside lists callback")
			}
			b.userStates[ID] = user.ListEditState
			return
		} else {
			message := fmt.Sprintf("Item number should be from 1 to %d", len(list.Items))
			reply := tgbotapi.NewMessage(ID, message)
			if err := b.apiRequest(reply); err != nil {
				fmt.Errorf("failed to send api request to delete an item")
			}
			reply = tgbotapi.NewMessage(ID, createListMessage(list))
			reply.ReplyMarkup = getListKeyboard()
			reply.ParseMode = "html"
			if err := b.apiRequest(reply); err != nil {
				fmt.Errorf("failed to send api request inside lists callback")
			}
			b.userStates[ID] = user.ListEditState
			return
		}

	}

}
