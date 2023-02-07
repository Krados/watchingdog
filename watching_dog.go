package watchingdog

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type Role string

const (
	ROLE_FOLLOWER Role = "follower"
	ROLE_LEADER   Role = "leader"
)

type WatchingDog struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	name       string
	lt         LeaderTaker
	watcher    Watcher
	alerter    Alerter
	wg         *sync.WaitGroup
	role       Role
}

func New(lt LeaderTaker, watcher Watcher, alerter Alerter) *WatchingDog {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &WatchingDog{
		name:       uuid.New().String(),
		ctx:        ctx,
		cancelFunc: cancelFunc,
		lt:         lt,
		watcher:    watcher,
		alerter:    alerter,
		wg:         new(sync.WaitGroup),
		role:       ROLE_FOLLOWER,
	}
}

func (w *WatchingDog) Start() {
	w.wg.Add(1)
	defer w.wg.Done()
	for {
		select {
		case <-w.ctx.Done():
			if w.role == ROLE_LEADER {
				w.lt.Revoke()
			}
			return
		default:
		}
		// if can not take leader then sleep
		isLeader, err := w.lt.TakeLeader()
		if err != nil { // err occur just continue
			continue
		}
		if !isLeader { // not the leader
			w.lt.Wait()
			continue
		}
		// leadership gained
		w.leaderShipGained(w.ctx)
	}
}

func (w *WatchingDog) leaderShipGained(ctx context.Context) {
	w.role = ROLE_LEADER
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		// extend leadership duration
		err := w.lt.ExtendDuration()
		// if extend failed then revoke leadership
		if err != nil {
			w.lt.Revoke()
			return
		}
		w.watch()
	}
}

func (w *WatchingDog) Stop() {
	w.cancelFunc()
	w.wg.Wait()
}

func (w *WatchingDog) Name() string {
	return w.name
}

func (w *WatchingDog) Role() Role {
	return w.role
}

func (w *WatchingDog) watch() {
	err := w.watcher.Watch()
	if err != nil {
		w.alert()
	}
}

func (w *WatchingDog) alert() {
	w.alerter.Alert()
}
