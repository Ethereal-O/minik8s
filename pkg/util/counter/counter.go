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

func GetUuid() string {
	uuid.lock.Lock()
	defer uuid.lock.Unlock()
	uuid.count++
	return strconv.Itoa(uuid.count)
}

func GetMonitorCrt() string {
	monitorCertification.lock.Lock()
	defer monitorCertification.lock.Unlock()
	monitorCertification.count++
	return strconv.Itoa(monitorCertification.count)
}
