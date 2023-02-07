package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Krados/watchingdog"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.168.13:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer cli.Close()

	lt := NewEtcdLeaderTaker(cli, "leaderKey", time.Second*10, time.Second*5)
	watcher := watchingdog.NewDummyWatcher(false)
	alerter := watchingdog.NewDummyAlerter()
	dog := watchingdog.New(lt, watcher, alerter)

	go func() {
		dog.Start()
	}()

	go func() {
		tc := time.NewTicker(time.Second)
		for range tc.C {
			fmt.Println(dog.Name(), dog.Role())
		}
	}()

	<-time.After(time.Second * 10)
	dog.Stop()
}

type EtcdLeaderTaker struct {
	cli           *clientv3.Client
	leaderKey     string
	leaseDuration time.Duration
	leaseID       clientv3.LeaseID
	waitDuration  time.Duration
}

func NewEtcdLeaderTaker(cli *clientv3.Client, leaderKey string, leaseDuration time.Duration, waitDuration time.Duration) watchingdog.LeaderTaker {
	return &EtcdLeaderTaker{
		cli:           cli,
		leaderKey:     leaderKey,
		leaseDuration: leaseDuration,
		waitDuration:  waitDuration,
	}
}

func (e *EtcdLeaderTaker) TakeLeader() (isLeader bool, err error) {
	leaseResp, err := e.cli.Grant(context.Background(), int64(e.leaseDuration.Seconds()))
	if err != nil {
		return
	}
	txResp, err := clientv3.NewKV(e.cli).Txn(context.Background()).
		If(clientv3.Compare(clientv3.Version(e.leaderKey), "=", 0)).
		Then(clientv3.OpPut(e.leaderKey, e.leaderKey, clientv3.WithLease(leaseResp.ID))).
		Commit()
	if err != nil {
		return
	}

	if !txResp.Succeeded {
		return
	}
	e.leaseID = leaseResp.ID

	return true, nil
}

func (e *EtcdLeaderTaker) ExtendDuration() (err error) {
	resp, err := e.cli.KeepAliveOnce(context.Background(), e.leaseID)
	_ = resp
	return
}

func (e *EtcdLeaderTaker) Revoke() (err error) {
	_, err = e.cli.Revoke(context.Background(), e.leaseID)
	return
}

func (e *EtcdLeaderTaker) Wait() {
	time.Sleep(e.waitDuration)
}
