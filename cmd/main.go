package main

import (
	listsbot "skeef79.com/simple-tg-bot/bot"
)

func main() {

	bot, err := listsbot.NewBot()
	if err != nil {
		panic(err)
	}

	bot.Run()
}
