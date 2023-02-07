package watchingdog

import (
	"errors"
	"time"
)

type DummyWatcher struct {
	MsgBag  []string
	ErrFlag bool
}

func NewDummyWatcher(errFlag bool) Watcher {
	return &DummyWatcher{
		MsgBag:  make([]string, 0),
		ErrFlag: errFlag,
	}
}

func (d *DummyWatcher) Watch() error {
	if d.ErrFlag {
		return d.errWork()
	}
	return d.successWork()
}

func (d *DummyWatcher) errWork() error {
	// simulate err work
	msg := "some error occur"
	d.MsgBag = append(d.MsgBag, msg)
	time.Sleep(time.Second)
	return errors.New(msg)
}

func (d *DummyWatcher) successWork() error {
	// simulate some work
	msg := "do some dummy work"
	d.MsgBag = append(d.MsgBag, msg)
	time.Sleep(time.Second * 2)
	return nil
}
