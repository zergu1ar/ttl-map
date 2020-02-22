package ttlmap

import (
	"sync"
	"time"
)

type TTLMap struct {
	mm    map[string]*entry
	mutex *sync.RWMutex
}

func New() *TTLMap {
	return &TTLMap{
		mm:    make(map[string]*entry),
		mutex: new(sync.RWMutex),
	}
}

func (slf *TTLMap) Add(key string, value interface{}, ttl time.Duration) bool {
	slf.mutex.Lock()
	defer slf.mutex.Unlock()

	_, found := slf.mm[key]
	if found {
		return false // Failed to add
	}

	slf.mm[key] = &entry{
		value:  value,
		expire: time.Now().Add(ttl),
	}

	return true
}

func (slf *TTLMap) Get(key string) (interface{}, bool) {
	slf.mutex.RLock()
	entry, exists := slf.mm[key]
	slf.mutex.RUnlock()

	if !exists {
		return nil, false
	} else {
		if time.Now().After(entry.expire) {
			defer slf.Del(key)
			return nil, false
		} else {
			return entry.value, true
		}
	}
}

func (slf *TTLMap) Del(key string) {
	slf.mutex.Lock()
	defer slf.mutex.Unlock()
	delete(slf.mm, key)
}

func (slf *TTLMap) Exists(key string) bool {
	_, hit := slf.Get(key)
	return hit
}

func (slf *TTLMap) Range(iterator func(key string, value interface{})) {
	willDelete := make(map[string]struct{})

	defer func() {
		slf.mutex.Lock()
		for key := range willDelete {
			delete(slf.mm, key)
		}
		slf.mutex.Unlock()
	}()

	slf.mutex.RLock()
	defer slf.mutex.RUnlock()

	for key, entry := range slf.mm {
		if time.Now().After(entry.expire) {
			willDelete[key] = struct{}{}
		} else {
			iterator(key, entry.value)
		}
	}
}

type entry struct {
	value  interface{}
	expire time.Time
}
