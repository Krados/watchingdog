package watchingdog

import (
	"sync"
	"time"
)

type DummyLeaderTaker struct {
	leaderFlag bool
	sync.Mutex
}

func NewDummyLeaderTaker() LeaderTaker {
	return &DummyLeaderTaker{}
}

func (d *DummyLeaderTaker) TakeLeader() (isLeader bool, err error) {
	d.Lock()
	defer d.Unlock()
	if !d.leaderFlag {
		d.leaderFlag = true
		return true, nil
	}
	return false, nil
}

func (d *DummyLeaderTaker) ExtendDuration() (err error) {
	return nil
}

func (d *DummyLeaderTaker) Revoke() (err error) {
	d.Lock()
	defer d.Unlock()
	d.leaderFlag = false
	return nil
}

func (d *DummyLeaderTaker) Wait() {
	time.Sleep(time.Second * 5)
}
