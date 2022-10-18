package main

import (
	"errors"
	"fmt"
	"os"
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

	time.Sleep(time.Second * 3)
	panic("crashed")
}
