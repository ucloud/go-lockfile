package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

const defaultCount = 20

func main() {
	var cntStr string
	if len(os.Args) >= 2 {
		cntStr = os.Args[1]
	}
	var taskCnt int
	var err error
	if cntStr != "" {
		taskCnt, err = strconv.Atoi(cntStr)
		if err != nil {
			panic(err)
		}
	}
	if taskCnt <= 0 {
		taskCnt = defaultCount
	}

	var wg sync.WaitGroup
	wg.Add(taskCnt)
	for i := 0; i < taskCnt; i++ {
		idx := i
		go func() {
			name := strconv.Itoa(idx)
			cmd := exec.Command("./out/sample_lock", name)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				fmt.Printf("failed to run %q: %v\n", name, err)
				os.Exit(1)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
