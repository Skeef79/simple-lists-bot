package storage

import (
	"fmt"
	"math/rand"

	"golang.org/x/exp/slices"
)

type List struct {
	Name  string
	ID    uint64
	Items []string
}

type Storage interface {
	GetAllLists() ([]*List, error)
	GetListByName(string) (*List, error)
	GetListById(uint64) (*List, error)
	CreateList(string) (*List, error)
	DeleteList(string) error
	AddItemByListName(string, string) error
	AddItemByListId(uint64, string) error
	DeleteItemByIndex(uint64, int) error
}

type InMemStorage struct {
	Name  string
	Lists []*List
}

func (s *InMemStorage) GetAllLists() ([]*List, error) {
	return s.Lists, nil
}

func (s *InMemStorage) GetListByName(name string) (*List, error) {
	for _, list := range s.Lists {
		if list.Name == name {
			return list, nil
		}
	}
	return nil, fmt.Errorf("List not found")
}

func (s *InMemStorage) GetListById(id uint64) (*List, error) {
	for _, list := range s.Lists {
		if list.ID == id {
			return list, nil
		}
	}
	return nil, fmt.Errorf("list not found")
}

func (s *InMemStorage) CreateList(name string) (*List, error) {
	newList := &List{
		Name:  name,
		ID:    rand.Uint64(),
		Items: make([]string, 0),
	}
	s.Lists = append(s.Lists, newList)
	return newList, nil
}

func (s *InMemStorage) DeleteList(name string) error {
	for i, list := range s.Lists {
		if list.Name == name {
			s.Lists = slices.Delete(s.Lists, i, i+1)
			return nil
		}
	}
	return fmt.Errorf("list not found")
}

func (s *InMemStorage) AddItemByListName(listName string, itemName string) error {
	list, err := s.GetListByName(listName)
	if err != nil {
		return err
	}
	list.Items = append(list.Items, itemName)
	return nil
}

func (s *InMemStorage) AddItemByListId(listID uint64, itemName string) error {
	list, err := s.GetListById(listID)
	if err != nil {
		return err
	}
	list.Items = append(list.Items, itemName)
	return nil
}

func (s *InMemStorage) DeleteItemByIndex(listID uint64, index int) error {
	for _, list := range s.Lists {
		if list.ID == listID {
			list.Items = slices.Delete(list.Items, index, index+1)
			return nil
		}
	}
	return fmt.Errorf("not found")
}

func NewInMemStorage(name string) Storage {
	return &InMemStorage{
		Name:  name,
		Lists: make([]*List, 0),
	}
}
