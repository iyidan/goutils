package safemap

import "sync"

// SafeMap concurrency map
type SafeMap struct {
	l sync.RWMutex
	m map[interface{}]interface{}
}

// New return an inited concurrency map
func New() *SafeMap {
	return &SafeMap{
		l: sync.RWMutex{},
		m: make(map[interface{}]interface{}),
	}
}

// Add if k already in the map, return false
func (m *SafeMap) Add(k interface{}, v interface{}) bool {
	m.l.Lock()
	defer m.l.Unlock()
	if _, ok := m.m[k]; !ok {
		m.m[k] = v
	} else {
		return false
	}
	return true
}

// Set set key k with value v
func (m *SafeMap) Set(k interface{}, v interface{}) {
	m.l.Lock()
	defer m.l.Unlock()
	m.m[k] = v
}

// CasSet compare and set v
func (m *SafeMap) CasSet(k interface{}, v interface{}, lastv interface{}) bool {
	m.l.Lock()
	defer m.l.Unlock()
	if tmpv, ok := m.m[k]; !ok || tmpv == lastv {
		m.m[k] = v
		return true
	}
	return false
}

// CasMultiSet compare and update multiple
func (m *SafeMap) CasMultiSet(update, old map[interface{}]interface{}) bool {
	m.l.Lock()
	defer m.l.Unlock()

	// old value compare
	for k, v := range old {
		if tmpv, ok := m.m[k]; !ok || tmpv != v {
			return false
		}
	}
	// new value set
	for k, v := range update {
		m.m[k] = v
	}
	return true
}

// Get get value of the given k, if k not exists, return nil
func (m *SafeMap) Get(k interface{}) interface{} {
	m.l.RLock()
	defer m.l.RUnlock()
	if val, ok := m.m[k]; ok {
		return val
	}
	return nil
}

// GetCheck get value of the given k, if k exists.
// if k not exists, ok is false
func (m *SafeMap) GetCheck(k interface{}) (interface{}, bool) {
	m.l.RLock()
	defer m.l.RUnlock()
	val, ok := m.m[k]
	return val, ok
}

// GetAll get a copy of m.m
func (m *SafeMap) GetAll() (cmap map[interface{}]interface{}) {
	m.l.RLock()
	defer m.l.RUnlock()
	cmap = make(map[interface{}]interface{})

	for k, v := range m.m {
		cmap[k] = v
	}
	return
}

// Exist check given k is exists
func (m *SafeMap) Exist(k interface{}) bool {
	m.l.RLock()
	defer m.l.RUnlock()
	if _, ok := m.m[k]; ok {
		return true
	}
	return false
}

// Del del the given key
func (m *SafeMap) Del(k interface{}) {
	m.l.Lock()
	defer m.l.Unlock()
	delete(m.m, k)
}

// DelAll delete all keys
func (m *SafeMap) DelAll() {
	m.l.Lock()
	defer m.l.Unlock()
	for k := range m.m {
		delete(m.m, k)
	}
}

// Len get the map keys count
func (m *SafeMap) Len() int {
	m.l.RLock()
	defer m.l.RUnlock()
	return len(m.m)
}
