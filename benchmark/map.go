package benchmark

import (
	"sync"

	cmap "github.com/orcaman/concurrent-map"
)

// HashMap is thread-safe map
type HashMap interface {
	Set(key string, val interface{})
	Get(key string) interface{}
	Del(key string)
	Len() int
	ForEach(func(key string, val interface{}) bool)
}

// BuildInMap is not a thread-safe type
type BuildInMap struct {
	items map[string]interface{}
}

func NewBuildInMap() *BuildInMap {
	m := &BuildInMap{}
	m.items = make(map[string]interface{}, 0)
	return m
}

func (m *BuildInMap) Set(key string, val interface{}) {
	m.items[key] = val
}

func (m *BuildInMap) Get(key string) interface{} {
	return m.items[key]
}

func (m *BuildInMap) Del(key string) {
	delete(m.items, key)
}

func (m *BuildInMap) Len() int {
	return len(m.items)
}

func (m *BuildInMap) ForEach(f func(key string, val interface{}) bool) {
	for k, v := range m.items {
		if !f(k, v) {
			return
		}
	}
}

// RWMutexMap
type RWMutexMap struct {
	items map[string]interface{}
	sync.RWMutex
}

func NewRWMutexMap() *RWMutexMap {
	m := &RWMutexMap{}
	m.items = make(map[string]interface{}, 0)
	return m
}

func (m *RWMutexMap) Set(key string, val interface{}) {
	m.Lock()
	defer m.Unlock()
	m.items[key] = val
}

func (m *RWMutexMap) Get(key string) interface{} {
	m.RLock()
	defer m.RUnlock()
	return m.items[key]
}

func (m *RWMutexMap) Del(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.items, key)
}

func (m *RWMutexMap) Len() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.items)
}

func (m *RWMutexMap) ForEach(f func(key string, val interface{}) bool) {
	m.RLock()
	defer m.RUnlock()
	for k, v := range m.items {
		if !f(k, v) {
			return
		}
	}
}


// MutexMap
type MutexMap struct {
	items map[string]interface{}
	sync.Mutex
}

func NewMutexMap() *MutexMap {
	m := &MutexMap{}
	m.items = make(map[string]interface{}, 0)
	return m
}

func (m *MutexMap) Set(key string, val interface{}) {
	m.Lock()
	defer m.Unlock()
	m.items[key] = val
}

func (m *MutexMap) Get(key string) interface{} {
	m.Lock()
	defer m.Unlock()
	return m.items[key]
}

func (m *MutexMap) Del(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.items, key)
}

func (m *MutexMap) Len() int {
	m.Lock()
	defer m.Unlock()
	return len(m.items)
}

func (m *MutexMap) ForEach(f func(key string, val interface{}) bool) {
	m.Lock()
	defer m.Unlock()
	for k, v := range m.items {
		if !f(k, v) {
			return
		}
	}
}

// SyncMap
type SyncMap struct {
	sync.Map
}

func (m *SyncMap) Set(key string, val interface{}) {
	m.Store(key, val)
}

func (m *SyncMap) Get(key string) interface{} {
	val, _ := m.Load(key)
	return val
}

func (m *SyncMap) Del(key string) {
	m.Delete(key)
}

func (m *SyncMap) Len() int {
	return m.Len()
}

func (m *SyncMap) ForEach(f func(key string, val interface{}) bool) {
	m.ForEach(f)
}

// ConcurrentMap
type ConcurrentMap struct {
	m cmap.ConcurrentMap
}

func NewConcurrentMap() *ConcurrentMap {
	cm := ConcurrentMap{}
	cm.m = cmap.New()
	return &cm
}

func (cm *ConcurrentMap) Set(key string, val interface{}) {
	cm.m.Set(key, val)
}

func (cm *ConcurrentMap) Get(key string) interface{} {
	val, _ := cm.m.Get(key)
	return val
}

func (cm *ConcurrentMap) Del(key string) {
	cm.m.Remove(key)
}

func (cm *ConcurrentMap) Len() int {
	return cm.m.Count()
}

func (cm *ConcurrentMap) ForEach(f func(key string, val interface{}) bool) {
	for item := range cm.m.IterBuffered() {
		f(item.Key, item.Val)
	}
}

