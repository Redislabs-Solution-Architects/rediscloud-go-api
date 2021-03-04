package kvstore

import "strings"

type kvmap struct {
	delimiter string
	prefix    string
	store     map[string]int
}

//KVStore is a key value store of some sort, which can be specialized for different use cases
type KVStore interface {
	//Put updates the value at the given key, overwriting any value already there.
	// A zero value key is a no-op.
	// A zero value value is equivalent to a delete.
	Put(key string, value int) KVStore
	//Get returns the value at the given key.
	//A zero value key will always return 0.
	Get(key string) int
	//Delete the value at the given key.
	//for all keys and values Put(k,v).Delete(k).Get == 0
	Delete(key string) KVStore
	//Keys returns all non-zero key values.
	Keys() []string
	//New returns a new and empty KVStore whose keys are distinct
	//from any other KVStore initialized with a different prefix (including the zero valued prefix)
	//The result of using the same prefix twice is unspecified.
	//The result of using the zero value prefix is unspecified.
	New(prefix string) KVStore
	//Copy returns a new KVStore whose keys and values are copied to the new store
	Copy(prefix string) KVStore
}

//NewKVMap creates a new kvmap
func NewKVMap() KVStore {
	return kvmap{
		delimiter: "/",
		prefix:    "",
		store:     make(map[string]int),
	}
}

//KVStore interface implementation
func (m kvmap) Put(key string, value int) KVStore {
	m.store[m.prefix+m.delimiter+key] = value
	return m
}

func (m kvmap) Get(key string) (v int) {
	v = m.store[m.prefix+m.delimiter+key]
	return
}

func (m kvmap) Delete(key string) KVStore {
	delete(m.store, m.prefix+m.delimiter+key)
	return m
}

func (m kvmap) Keys() (keys []string) {
	for k := range m.store {
		parts := strings.Split(k, m.delimiter)
		ac := parts[len(parts)-1]
		keys = append(keys, ac)
	}
	return
}

//New is a no-op for the
func (m kvmap) New(prefix string) KVStore {
	return kvmap{
		delimiter: "/",
		prefix:    prefix,
		store:     make(map[string]int),
	}
}

func (m kvmap) Copy(prefix string) KVStore {
	new := m.New(prefix)
	for _, k := range m.Keys() {
		v := m.Get(k)
		new.Put(k, v)
	}
	return new
}
