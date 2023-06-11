package storage

type List struct {
	Name    string
	ID      uint64
	Items   []string
	ItemIDs []uint64
}

type Storage interface {
	GetAllLists() ([]*List, error)
	GetListByName(string) (*List, error)
	GetListById(uint64) (*List, error)
	CreateList(string) (*List, error)
	DeleteList(string) error
	AddItemByListId(uint64, string) error
	DeleteItemByIndex(uint64, int) error
}
