package counter

import (
	"strconv"
	"sync"
)

type Counter struct {
	count int
	lock  sync.Mutex
}

var uuid = Counter{count: 10000}
var monitorCertification = Counter{count: 10000}
var rrPolicyCounter = Counter{count: 0}

func (counter *Counter) fetchAndAdd() int {
	counter.lock.Lock()
	defer counter.lock.Unlock()
	count := counter.count
	counter.count++
	return count
}

func GetUuid() string {
	return strconv.Itoa(uuid.fetchAndAdd())
}

func GetMonitorCrt() string {
	return strconv.Itoa(monitorCertification.fetchAndAdd())
}

func GetRRPolicy() int {
	return rrPolicyCounter.fetchAndAdd()
}
