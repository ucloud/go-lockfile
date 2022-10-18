package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/ucloud/go-lockfile"
)

func main() {
	fl := lockfile.New("busy.lock")
	var err error
	for i := 0; i < 20; i++ {
		err = fl.TryLock()
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 50)
	}
	if err != nil {
		if errors.Is(err, lockfile.ErrBusy) {
			fmt.Println("The lock is busy")
		} else {
			fmt.Printf("Unexpected error: %v\n", err)
		}
		os.Exit(1)
	}

	fmt.Println("lock!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	go func() {
		<-sc
		err := fl.Unlock()
		if err != nil {
			fmt.Printf("failed to unlock: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("unlock!")
		os.Exit(0)
	}()

	select {}
}
