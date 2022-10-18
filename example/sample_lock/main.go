package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ucloud/go-lockfile"
)

func main() {
	var name string
	if len(os.Args) >= 2 {
		name = os.Args[1]
	}
	if name == "" {
		name = "default"
	}
	lf := lockfile.New("test.lock")

	fmt.Printf("%q begin to acquire lock...\n", name)
	retries := 10000
	for {
		retries--
		if retries <= 0 {
			fmt.Printf("%q failed to acquire lock, lock is too busy\n", name)
			return
		}
		err := lf.TryLock()
		switch err {
		case lockfile.ErrBusy:
			time.Sleep(time.Millisecond * 100)

		case nil:
			fmt.Printf("%q acquired lock, begin to handle\n", name)
			hanle(name)
			if strings.HasPrefix(name, "crash") {
				fmt.Printf("%q crashed, without releasing lock\n", name)
				return
			}
			err = lf.Unlock()
			if err != nil {
				fmt.Printf("%q failed to unlock: %v.\n", name, err)
				return
			}
			fmt.Printf("%q released lock\n", name)
			return

		default:
			fmt.Printf("%q failed to acquire lock: %v\n", name, err)
			return
		}
	}
}

func hanle(name string) {
	total := time.NewTimer(time.Second * 5)
	tk := time.NewTicker(time.Millisecond * 300)
	var idx int
	for {
		select {
		case <-tk.C:
			fmt.Printf("%q handle step %d\n", name, idx)
			idx++
		case <-total.C:
			return
		}
	}
}
