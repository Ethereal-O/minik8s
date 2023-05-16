package structure

import "sync"

type CountMap struct {
	lock   sync.Mutex
	record map[string]int
}

func (cm *CountMap) Init() {
	cm.record = make(map[string]int)
}

func (cm *CountMap) Exist(key string) bool {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	if _, ok := cm.record[key]; ok {
		return true
	}
	return false
}

func (cm *CountMap) Add(key string) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.record[key]++
}

func (cm *CountMap) Minus(key string) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.record[key]--
}

func (cm *CountMap) Get(key string) int {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	return cm.record[key]
}

func (cm *CountMap) Delete(key string) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	delete(cm.record, key)
}
