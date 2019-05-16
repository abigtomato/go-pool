package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/abigtomato/go-pool/pkg/pool"
)

func main() {
	var wg sync.WaitGroup
	defer pool.Close()

	runTimes := 100000
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		if err := pool.Go(func() {
			time.Sleep(10 * time.Millisecond)
			fmt.Println("Hello Goroutine Pool!")
			wg.Done()
		}); err != nil {
			log.Println(err.Error())
		}
	}
	wg.Wait()
}
