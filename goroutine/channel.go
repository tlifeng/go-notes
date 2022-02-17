package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var wg sync.WaitGroup
var counter int32
var mux sync.Mutex
var ch = make(chan int, 1)

func AtomicIntCounter() {
	defer wg.Done()
	for i := 0; i < 10000; i++ {
		atomic.AddInt32(&counter, 1)
	}
}

func MutexIntCounter()  {
	defer wg.Done()
	for i := 0; i < 10000; i++ {
		mux.Lock()
		counter++
		mux.Unlock()
	}
}

func ChannelIntCounter() {
	defer wg.Done()
	for i := 0; i < 10000; i++ {
		count := <- ch
		count++
		ch <- count
	}
}


func main() {
	wg.Add(2)

	ch <- 0

	go ChannelIntCounter()
	go ChannelIntCounter()

	wg.Wait()

	fmt.Println(<-ch)
}