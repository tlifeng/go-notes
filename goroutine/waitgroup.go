package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func say(id string) {
	time.Sleep(time.Second)
	fmt.Println("I am done! id: " + id)
	wg.Done()
}

func main() {
	wg.Add(2)

	// goroutine可以用在匿名函数上
	go func() {
		fmt.Println("Hello")
		wg.Done()
	}()
	go say("world")

	wg.Wait()
	fmt.Println("exit")
}