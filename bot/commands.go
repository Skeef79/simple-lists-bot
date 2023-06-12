package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type commandKey string

type commandEntity struct {
	Key    commandKey
	Desc   string
	Action func(upd tgbotapi.Update)
}

const (
	StartCmdKey      = commandKey("start")
	HelpCmdKey       = commandKey("help")
	ShowListsCmdKey  = commandKey("lists")
	AddListCmdKey    = commandKey("add_list")
	DeleteListCmdKey = commandKey("delete_list")
	// AddItemCmdKey    = commandKey("add_item")
	// RemoveItemCmdKey = commandKey("remove_item")
)

func (b *bot) initCommands() error {
	commands := []commandEntity{
		{
			Key:    StartCmdKey,
			Desc:   "Run bot",
			Action: b.StartCmd,
		},
		{
			Key:    HelpCmdKey,
			Desc:   "Help",
			Action: b.HelpCmd,
		},
		{
			Key:    ShowListsCmdKey,
			Desc:   "Show all lists",
			Action: b.ShowListsCmd,
		},
		{
			Key:    AddListCmdKey,
			Desc:   "Add new list",
			Action: b.AddListCmd,
		},
		{
			Key:    DeleteListCmdKey,
			Desc:   "Delete list",
			Action: b.DeleteListCmd,
		},
	}

	tgCommands := make([]tgbotapi.BotCommand, 0, len(commands))

	for _, cmd := range commands {
		b.commands[cmd.Key] = cmd
		tgCommands = append(tgCommands, tgbotapi.BotCommand{
			Command:     "/" + string(cmd.Key),
			Description: cmd.Desc,
		})
	}

	cfg := tgbotapi.NewSetMyCommands(tgCommands...)
	return b.apiRequest(cfg)
}
