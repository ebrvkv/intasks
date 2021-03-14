package counter

import "sync"

type Request interface {
	Get() int64
	Inc()
	Reduce()
	Stop() bool
	Stopped() bool
}

type requestImpl struct {
	sync.Mutex
	cnt     int64
	stopped bool
}

func NewReqCounter() Request {
	return &requestImpl{}
}

func (rc *requestImpl) Stop() bool {
	rc.Lock()
	defer rc.Unlock()
	if rc.stopped {
		return false
	}
	rc.stopped = true
	return rc.stopped
}

func (rc *requestImpl) Stopped() bool {
	rc.Lock()
	defer rc.Unlock()
	return rc.stopped
}

func (rc *requestImpl) Get() int64 {
	rc.Lock()
	defer rc.Unlock()
	return rc.cnt
}

func (rc *requestImpl) Inc() {
	rc.Lock()
	rc.cnt++
	rc.Unlock()
}

func (rc *requestImpl) Reduce() {
	rc.Lock()
	rc.cnt--
	rc.Unlock()
}
