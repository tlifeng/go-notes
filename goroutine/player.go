package main

import (
	"fmt"
	"math/rand"
	"sync"
)

var wg sync.WaitGroup

func player(name string, ch chan int) {
	defer wg.Done()

	for {
		ball, ok := <- ch

		if !ok {
			fmt.Printf("channel is closed! %s wins!\n", name)
			return
		}

		n := rand.Intn(100)

		if n%10 == 0 {
			close(ch)
			return
		}

		ball++
		fmt.Printf("%s recives ball %d\n", name, ball)
		ch <- ball
	}
}

func main() {
	wg.Add(2)

	ch := make(chan int, 0) //unbuffered channel

	go player("ming",ch)
	go player("hong",ch)

	ch <- 0

	wg.Wait()
}
