package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		fmt.Println("1111111")
		fmt.Println("2222222")
	}()
	time.Sleep(time.Second)
}
