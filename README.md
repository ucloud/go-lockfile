# go-lockfile

[![GoPkg Widget]][GoPkg] [![Go Report Card](https://goreportcard.com/badge/github.com/ucloud/go-lockfile)](https://goreportcard.com/report/github.com/ucloud/go-lockfile)

A Linux go library to lock cooperating processes based on syscall [flock](https://man7.org/linux/man-pages/man2/flock.2.html).

## Install

```shell
go get github.com/ucloud/go-lockfile
```

## Usage

A simple example:

```go
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

	fmt.Println("lock!")
	// Handle your logic
	time.Sleep(time.Second)

	err = fl.Unlock()
	if err != nil {
		fmt.Printf("failed to unlock: %v\n", err)
		os.Exit(1)
	}
}
```

## Documentation

See: [go-lockfile package](https://pkg.go.dev/github.com/ucloud/go-lockfile)
