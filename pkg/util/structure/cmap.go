package structure

import (
	"sync"
)

type Cmap struct {
	lock       sync.Mutex
	record     map[string]int
	capability int
}

func (r *Cmap) Init(cap int) {
	r.record = make(map[string]int)
	r.capability = cap
}

func (r *Cmap) Put(s string) bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.record[s] < r.capability {
		r.record[s]++
		return true
	} else {
		return false
	}
}

func (r *Cmap) Get(s string) bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.record[s] > 0 {
		r.record[s]--
		return true
	} else {
		return false
	}
}
