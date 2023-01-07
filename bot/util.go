package bot

import (
	"fmt"

	storage "github.com/skeef79/simple-lists-bot/storage"
)

const (
	AddItemMessage    = "Add item"
	DeleteItemMessage = "Delete item"
	BackMessage       = "Back"
)

func createListMessage(list *storage.List) string {
	message := fmt.Sprintf("<b> %s </b>\n", list.Name)
	for idx, item := range list.Items {
		message += fmt.Sprintf("%d) %s\n", idx+1, item)
	}
	return message
}
