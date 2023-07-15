package cache

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRU_Add_existElementInFullQueue_moveToBack(t *testing.T) {
	//Arrange
	lru := NewLRU(3)
	lru.Add("someKey1", []int{1, 2})
	lru.Add("someKey2", "someValue2")
	lru.Add("someKey3", 333)

	//Act
	lru.Add("someKey1", 100)

	//Assert
	frontItem := lru.queue.Front().Value.(*Item)
	assert.Equal(t, "someKey2", frontItem.Key)
	assert.Equal(t, "someValue2", frontItem.Value)
	backItem := lru.queue.Back().Value.(*Item)
	assert.Equal(t, "someKey1", backItem.Key)
	assert.Equal(t, 100, backItem.Value)
	assert.Equal(t, 3, lru.queue.Len())
}

func TestLRU_Add_existElement_moveToBack(t *testing.T) {
	//Arrange
	lru := NewLRU(3)
	lru.Add("someKey1", []int{1, 2})
	lru.Add("someKey2", "someValue2")

	//Act
	lru.Add("someKey1", 100)

	//Assert
	frontItem := lru.queue.Front().Value.(*Item)
	assert.Equal(t, "someKey2", frontItem.Key)
	assert.Equal(t, "someValue2", frontItem.Value)
	backItem := lru.queue.Back().Value.(*Item)
	assert.Equal(t, "someKey1", backItem.Key)
	assert.Equal(t, 100, backItem.Value)
	assert.Equal(t, 2, lru.queue.Len())
}

func TestLRU_Add_newElementInFullQueue_delFrontAndPushToBack(t *testing.T) {
	//Arrange
	lru := NewLRU(3)
	lru.Add("someKey1", []int{1, 2})
	lru.Add("someKey2", "someValue2")
	lru.Add("someKey3", 333)

	//Act
	lru.Add("someKey4", 100)

	//Assert
	frontItem := lru.queue.Front().Value.(*Item)
	assert.Equal(t, "someKey2", frontItem.Key)
	assert.Equal(t, "someValue2", frontItem.Value)
	backItem := lru.queue.Back().Value.(*Item)
	assert.Equal(t, "someKey4", backItem.Key)
	assert.Equal(t, 100, backItem.Value)
	assert.Equal(t, 3, lru.queue.Len())
}

func TestLRU_Add_newElement_pushToBack(t *testing.T) {
	//Arrange
	lru := NewLRU(3)
	lru.Add("someKey1", []int{1, 2})
	lru.Add("someKey2", "someValue2")

	//Act
	lru.Add("someKey3", 100)

	//Assert
	frontItem := lru.queue.Front().Value.(*Item)
	assert.Equal(t, "someKey1", frontItem.Key)
	assert.Equal(t, []int{1, 2}, frontItem.Value)
	backItem := lru.queue.Back().Value.(*Item)
	assert.Equal(t, "someKey3", backItem.Key)
	assert.Equal(t, 100, backItem.Value)
	assert.Equal(t, 3, lru.queue.Len())
}

func TestLRU_Add_newElementAsync_allKeysExists(t *testing.T) {
	//Arrange
	wg := sync.WaitGroup{}
	lru := NewLRU(3)
	wg.Add(3)

	//Act
	go func() {
		lru.Add("someKey1", []int{1, 2})
		wg.Done()
	}()
	go func() {
		lru.Add("someKey2", "someValue2")
		wg.Done()
	}()
	go func() {
		lru.Add("someKey3", 100)
		wg.Done()
	}()
	wg.Wait()

	//Assert
	_, ok := lru.items["someKey1"]
	assert.Equal(t, true, ok)
	_, ok = lru.items["someKey2"]
	assert.Equal(t, true, ok)
	_, ok = lru.items["someKey3"]
	assert.Equal(t, true, ok)
	assert.Equal(t, 3, lru.queue.Len())
}

func TestLRU_Get_hasElement_returnItAndMoveToBack(t *testing.T) {
	//Arrange
	lru := NewLRU(3)
	lru.Add("someKey1", []int{1, 2})
	lru.Add("someKey2", "someValue2")
	lru.Add("someKey3", 333)

	//Act
	v := lru.Get("someKey2")

	//Assert
	assert.Equal(t, "someValue2", v)
	backItem := lru.queue.Back().Value.(*Item)
	assert.Equal(t, "someKey2", backItem.Key)
	assert.Equal(t, "someValue2", backItem.Value)
	assert.Equal(t, 3, lru.queue.Len())
}

func TestLRU_Get_hasNotElement_returnNil(t *testing.T) {
	//Arrange
	lru := NewLRU(3)
	lru.Add("someKey1", []int{1, 2})
	lru.Add("someKey2", "someValue2")
	lru.Add("someKey3", 333)

	//Act
	v := lru.Get("someKey4")

	//Assert
	assert.Nil(t, v)
}

func TestLRU_Len_Nothing_returnLen(t *testing.T) {
	//Arrange
	lru := NewLRU(3)
	lru.Add("someKey1", []int{1, 2})
	lru.Add("someKey2", "someValue2")
	lru.Add("someKey3", 333)

	//Act
	len := lru.Len()

	//Assert
	assert.Equal(t, 3, len)
}

func TestLRU_Delete_hasElement_deleteIt(t *testing.T) {
	//Arrange
	lru := NewLRU(3)
	lru.Add("someKey1", []int{1, 2})
	lru.Add("someKey2", "someValue2")
	lru.Add("someKey3", 333)

	//Act
	result := lru.Delete("someKey3")

	//Assert
	assert.True(t, result)
	assert.Equal(t, 2, lru.queue.Len())
	backItem := lru.queue.Back().Value.(*Item)
	assert.Equal(t, "someKey2", backItem.Key)
}

func TestLRU_Delete_hasNotElement_returnFalse(t *testing.T) {
	//Arrange
	lru := NewLRU(3)
	lru.Add("someKey1", []int{1, 2})
	lru.Add("someKey2", "someValue2")
	lru.Add("someKey3", 333)

	//Act
	result := lru.Delete("someKey4")

	//Assert
	assert.False(t, result)
	assert.Equal(t, 3, lru.queue.Len())
}
