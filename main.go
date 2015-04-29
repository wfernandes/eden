package main

import (
	"fmt"
	"github.com/fzzy/radix/extra/cluster"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	//	"github.com/fzzy/radix/redis"
)

var b = make([]rune, 100)
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var i int64 = 0

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	closeChan := make(chan os.Signal, 1)
	signal.Notify(closeChan, os.Interrupt)

	var numWrites int64 = 0

	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Printf("Writes per second %d\n", numWrites)
			atomic.StoreInt64(&numWrites, 0)
		}
	}()

	var wg sync.WaitGroup

	doneChan := make(chan struct{})

	go func() {
		select {
		case <-closeChan:
			close(doneChan)
		}
	}()

	for j := 0; j < 50; j++ {
		wg.Add(1)
		go func(closeCh chan os.Signal) {
			c, err := cluster.NewCluster("10.10.60.50:7000")
			//			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			defer c.Close()

			if err != nil {
				fmt.Println("error:", err)
				return
			}
		loop:
			for {
				key := rand.Int63()
				value := getNextValue()
				appId := fmt.Sprintf("d49bea3a-6fc5-4001-83e0-279019b9%d", rand.Intn(10000))
				r := c.Cmd("zadd", appId, key, value)

				if r.Err != nil {
					log.Printf("Error setting %d to %s: %s\n", key, value, r.Err.Error())
				} else {
					atomic.AddInt64(&numWrites, 1)
				}

				select {
				case <-doneChan:
					break loop
				default:
				}
			}
			wg.Done()
		}(closeChan)
	}
	wg.Wait()
}

func getNextValue() string {
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
