package storage

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DbStorage struct {
	DB     *gorm.DB
	ChatID int64
}

type ListItems struct {
	ID     uint64 `gorm:"column:id"`
	ListID uint64 `gorm:"column:list_id"`
	Name   string `gorm:"column:name"`
}

type Lists struct {
	ID     uint64 `gorm:"column:id"`
	Name   string `gorm:"column:name"`
	ChatID int64  `gorm:"column:chat_id"`
}

//TODO: pass db config here or load it from config.json

func NewDbStorage(chatID int64) Storage {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/lists?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_LISTS_BOT_USER"), os.Getenv("MYSQL_LISTS_BOT_PASSWORD"))
	s := &DbStorage{}
	var err error
	s.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	s.ChatID = chatID
	log.Println("Successfully connected to mysql db")
	return s
}

func getAllLists(db *gorm.DB, chatID int64) ([]*List, error) {
	var rows []*Lists
	if err := db.Where("chat_id = ?", chatID).Find(&rows).Error; err != nil {
		return nil, err
	}
	res := make([]*List, 0, len(rows))
	for _, row := range rows {
		res = append(res, &List{
			ID:    row.ID,
			Name:  row.Name,
			Items: make([]string, 0),
		})
	}

	return res, nil
}

func getListByID(db *gorm.DB, id uint64) (*List, error) {
	var row Lists
	if err := db.First(&row, id).Error; err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, fmt.Errorf("list not found")
	}

	return &List{
		ID:   row.ID,
		Name: row.Name,
	}, nil
}

func getListByName(db *gorm.DB, name string) (*List, error) {
	var row Lists
	if err := db.Where("name = ?", name).First(&row).Error; err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, fmt.Errorf("list not found")
	}

	return &List{
		ID:   row.ID,
		Name: row.Name,
	}, nil
}

func getListItems(db *gorm.DB, id uint64) ([]string, []uint64, error) {
	var row Lists
	if err := db.First(&row, id).Error; err != nil {
		return nil, nil, err
	}
	if row.ID == 0 {
		return nil, nil, fmt.Errorf("list not found")
	}

	rows := make([]*ListItems, 0)
	if err := db.Where("list_id = ?", id).Find(&rows).Error; err != nil {
		return nil, nil, err
	}
	listNames := make([]string, 0, len(rows))
	listIDs := make([]uint64, 0, len(rows))
	for _, row := range rows {
		listNames = append(listNames, row.Name)
		listIDs = append(listIDs, row.ID)
	}

	return listNames, listIDs, nil
}

func (s *DbStorage) GetAllLists() ([]*List, error) {
	return getAllLists(s.DB, s.ChatID)
}

func (s *DbStorage) GetListByName(name string) (*List, error) {
	list, err := getListByName(s.DB, name)
	if err != nil {
		return nil, err
	}
	list.Items, list.ItemIDs, err = getListItems(s.DB, list.ID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *DbStorage) GetListById(id uint64) (*List, error) {
	list, err := getListByID(s.DB, id)
	if err != nil {
		return nil, err
	}
	list.Items, list.ItemIDs, err = getListItems(s.DB, list.ID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func newList(name string) Lists {
	return Lists{
		Name: name,
	}
}

func newListItem(listID uint64, itemName string) ListItems {
	return ListItems{
		ListID: listID,
		Name:   itemName,
	}
}

func (s *DbStorage) CreateList(name string) (*List, error) {
	list := newList(name)
	if err := s.DB.Create(&list).Error; err != nil {
		return nil, err
	}
	res, err := getListByID(s.DB, list.ID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *DbStorage) DeleteList(id uint64) error {
	list, err := getListByID(s.DB, id)
	if err != nil {
		return err
	}

	listID := list.ID
	if err := s.DB.Delete(&Lists{}, &list.ID).Error; err != nil {
		return err
	}

	if err := s.DB.Where("list_id = ?", listID).Delete(&ListItems{}).Error; err != nil {
		return err
	}

	return nil
}

func (s *DbStorage) AddItemByListId(listID uint64, itemName string) error {
	_, err := getListByID(s.DB, listID)
	if err != nil {
		return err
	}

	listItem := newListItem(listID, itemName)
	if err := s.DB.Create(&listItem).Error; err != nil {
		return err
	}

	return nil
}

func (s *DbStorage) DeleteItemByIndex(listID uint64, index int) error {
	list, err := getListByID(s.DB, listID)
	if err != nil {
		return err
	}

	list.Items, list.ItemIDs, err = getListItems(s.DB, listID)
	if err != nil {
		return err
	}

	itemID := list.ItemIDs[index]

	if err := s.DB.Delete(ListItems{}, itemID).Error; err != nil {
		return err
	}
	return nil
}
