package main

import (
	listsbot "github.com/skeef79/simple-lists-bot/bot"
)

func main() {

	bot, err := listsbot.NewBot()
	if err != nil {
		panic(err)
	}

	bot.Run()
}
