package consensus

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitQueue(t *testing.T) {
	list := NewQueue()
	assert.Equal(t, 0, list.Len())
	assert.Nil(t, list.Head())
	assert.Nil(t, list.Tail())
}

func TestQueue_AddNode(t *testing.T) {
	list := NewQueue()

	block1 := NewBlock("blockHash1")
	block2 := NewBlock("blockHash2")
	block3 := NewBlock("blockHash3")
	list.Add(&block1)
	assert.Equal(t, 1, list.Len())
	assert.Equal(t, "blockHash1", list.Head().Value.Hash())
	assert.Nil(t, list.Head().Next)
	assert.Equal(t, "blockHash1", list.Tail().Value.Hash())

	list.Add(&block2)
	assert.Equal(t, 2, list.Len())
	assert.Equal(t, "blockHash1", list.Head().Value.Hash())
	assert.NotNil(t, list.Head().Next)
	assert.Equal(t, "blockHash2", list.Tail().Value.Hash())
	assert.Equal(t, "blockHash1", list.Tail().Prev.Value.Hash())

	list.Add(&block3)
	assert.Equal(t, 3, list.Len())
	assert.Equal(t, "blockHash3", list.Tail().Value.Hash())
	assert.Equal(t, "blockHash2", list.Tail().Prev.Value.Hash())
}

func TestQueue_RemoveHead_HeadNil(t *testing.T) {
	list := NewQueue()
	actual := list.Pop()
	assert.Equal(t, 0, list.Len())
	assert.Nil(t, list.Head())
	assert.Nil(t, list.Tail())
	assert.Nil(t, actual)
}

func TestQueue_RemoveHead_OneItem(t *testing.T) {
	block1 := NewBlock("blockHash1")
	list := NewQueue()
	list.Add(&block1)
	actual := list.Pop()
	assert.Equal(t, 0, list.Len())
	assert.Nil(t, list.Head())
	assert.Nil(t, list.Tail())
	assert.Equal(t, "blockHash1", actual.Value.Hash())
}

func TestQueue_RemoveHead(t *testing.T) {

	block1 := NewBlock("blockHash1")
	block2 := NewBlock("blockHash2")
	list := NewQueue()
	list.Add(&block1)
	list.Add(&block2)
	actual := list.Pop()
	assert.Equal(t, 1, list.Len())
	assert.Equal(t, "blockHash2", list.Head().Value.Hash())
	assert.Nil(t, list.Head().Prev)
	assert.Nil(t, list.Head().Next)

	assert.Equal(t, "blockHash2", list.Tail().Value.Hash())
	assert.Nil(t, list.Head().Prev)
	assert.Nil(t, list.Head().Next)

	assert.Equal(t, "blockHash1", actual.Value.Hash())
}
