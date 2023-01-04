package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateList(t *testing.T) {
	storage := NewInMemStorage("test-storage")
	list1, err := storage.CreateList("list1")

	require.NoError(t, err)
	require.Equal(t, list1.Name, "list1")

	allLists, err := storage.GetAllLists()
	require.NoError(t, err)
	require.Equal(t, len(allLists), 1)
	require.Equal(t, allLists[0].Name, "list1")

	list2, err := storage.CreateList("list2")
	require.NoError(t, err)
	require.Equal(t, list2.Name, "list2")

	allLists, err = storage.GetAllLists()
	require.NoError(t, err)
	require.Equal(t, len(allLists), 2)
	require.Equal(t, allLists[0].Name, "list1")
	require.Equal(t, allLists[1].Name, "list2")
}

// TODO: Add more tests, remember about TDD!
