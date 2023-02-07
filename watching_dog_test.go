package watchingdog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsLeader(t *testing.T) {
	dog := New(NewDummyLeaderTaker(), NewDummyWatcher(false), NewDummyAlerter())
	go func() {
		dog.Start()
	}()
	time.Sleep(time.Second)
	role := dog.Role()
	dog.Stop()
	assert.Equal(t, ROLE_LEADER, role)
}

func TestIsFollower(t *testing.T) {
	lt := NewDummyLeaderTaker()
	lt.TakeLeader()
	dog := New(lt, NewDummyWatcher(false), NewDummyAlerter())
	go func() {
		dog.Start()
	}()
	time.Sleep(time.Second)
	role := dog.Role()
	dog.Stop()
	assert.Equal(t, ROLE_FOLLOWER, role)
}

func TestConcurrency(t *testing.T) {
	lt := NewDummyLeaderTaker()
	dogList := make([]*WatchingDog, 0)
	dog1 := New(lt, NewDummyWatcher(false), NewDummyAlerter())
	dogList = append(dogList, dog1)
	dog2 := New(lt, NewDummyWatcher(false), NewDummyAlerter())
	dogList = append(dogList, dog2)
	dog3 := New(lt, NewDummyWatcher(false), NewDummyAlerter())
	dogList = append(dogList, dog3)

	startSignal := make(chan struct{})
	go func() {
		<-startSignal
		dog1.Start()
	}()
	go func() {
		<-startSignal
		dog2.Start()
	}()
	go func() {
		<-startSignal
		dog3.Start()
	}()
	time.Sleep(time.Second)
	close(startSignal)
	time.Sleep(time.Second * 2)

	leaderCount := 0
	for i := 0; i < len(dogList); i++ {
		if dogList[i].Role() == ROLE_LEADER {
			leaderCount++
		}
	}
	assert.Equal(t, 1, leaderCount)
	dog1.Stop()
	dog2.Stop()
	dog3.Stop()
}

func TestExecuteWatch(t *testing.T) {
	lt := NewDummyLeaderTaker()
	watcher := &DummyWatcher{MsgBag: make([]string, 0), ErrFlag: false}
	dog := New(lt, watcher, NewDummyAlerter())
	go func() {
		dog.Start()
	}()
	time.Sleep(time.Second * 2)
	dog.Stop()
	watchGotHit := false
	if len(watcher.MsgBag) > 0 {
		watchGotHit = true
	}
	assert.Equal(t, true, watchGotHit)
}

func TestExecuteAlert(t *testing.T) {
	lt := NewDummyLeaderTaker()
	alerter := &DummyAlerter{MsgBag: make([]string, 0)}
	dog := New(lt, NewDummyWatcher(true), alerter)
	go func() {
		dog.Start()
	}()
	time.Sleep(time.Second * 2)
	dog.Stop()
	alertGotHit := false
	if len(alerter.MsgBag) > 0 {
		alertGotHit = true
	}
	assert.Equal(t, true, alertGotHit)
}
