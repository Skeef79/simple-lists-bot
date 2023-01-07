package user

import (
	storage "skeef79.com/simple-tg-bot/storage"
)

const (
	EmptyState int64 = iota
	CreateListState
	ListEditState
	ShowListsState
	AddItemState
	DeleteItemState
)

type User struct {
	storage.Storage
}
