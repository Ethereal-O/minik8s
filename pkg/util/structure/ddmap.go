package structure

import "sync"

type DoubleDirectionMap struct {
	lock    sync.Mutex
	mem2key map[string]string
	key2num map[string]int
}

func (ddmap *DoubleDirectionMap) Init() {
	ddmap.mem2key = make(map[string]string)
	ddmap.key2num = make(map[string]int)
}

func (ddmap *DoubleDirectionMap) Put(key string, mem string) int {
	ddmap.lock.Lock()
	defer ddmap.lock.Unlock()
	ddmap.mem2key[mem] = key
	ddmap.key2num[key]++
	return ddmap.key2num[key]
}

func (ddmap *DoubleDirectionMap) Get(mem string) int {
	ddmap.lock.Lock()
	defer ddmap.lock.Unlock()
	key := ddmap.mem2key[mem]
	ddmap.key2num[key]--
	delete(ddmap.mem2key, mem)
	return ddmap.key2num[key]
}
