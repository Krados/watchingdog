package main

import (
	"fmt"
	"time"

	"github.com/Krados/watchingdog"
)

func main() {
	lt := watchingdog.NewDummyLeaderTaker()
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
