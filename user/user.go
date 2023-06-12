package user

import (
	storage "github.com/skeef79/simple-lists-bot/storage"
)

const (
	EmptyState int64 = iota
	CreateListState
	ListEditState
	ShowListsState
	AddItemState
	DeleteItemState
	DeleteListState
)

type User struct {
	storage.Storage
}
