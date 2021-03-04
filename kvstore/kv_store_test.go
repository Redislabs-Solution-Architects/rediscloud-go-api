package kvstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	var uut1 KVStore
	uut1 = NewKVMap()
	uut1.Put("this", int(1))
	uut1.Put("that", 2)
	assert.EqualValues(t, []string{"this", "that"}, uut1.Keys())
}
func TestCopy(t *testing.T) {
	var uut1 KVStore
	uut1 = NewKVMap()
	uut1.Put("this", int(1))
	uut2 := uut1.Copy("fargle")
	assert.Equal(t, 1, uut2.Get("this"))
	uut2.Put("this", 5)
	assert.Equal(t, 1, uut1.Get("this"))
	assert.Equal(t, 5, uut2.Get("this"))
	// assert.EqualValues(t, []string{"this"}, uut1.Keys())
	assert.EqualValues(t, []string{"this"}, uut2.Keys())
}
func TestNewKVStore(t *testing.T) {
	var uut KVStore
	uut = NewKVMap()
	uut.Put("this", int(1))
	expected := 1
	actual := uut.Get("this")
	assert.Equal(t, expected, actual)
}

func TestNewNewKVStore(t *testing.T) {
	var uut KVStore
	uut = NewKVMap().New("prefix")
	uut.Put("this", int(1))
	expected := 1
	actual := uut.Get("this")
	assert.Equal(t, expected, actual)
}

func TestNoInterference(t *testing.T) {
	var uut1, uut2 KVStore
	uut1 = NewKVMap()
	uut2 = uut1.New("prefix")
	uut1.Put("this", 1)
	uut2.Put("this", 2)
	assert.Equal(t, 1, uut1.Get("this"))
	assert.Equal(t, 2, uut2.Get("this"))
}
