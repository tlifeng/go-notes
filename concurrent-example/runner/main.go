package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func createTask() func(int)  {
	return func(id int) {
		time.Sleep(time.Second)
		fmt.Printf("Task complete #%d \n", id)
	}
}

func main() {
	r := New(4 * time.Second)

	r.AddTask(createTask(), createTask(), createTask())

	err := r.Start()

	switch err {
	case ErrInterrupt:
		fmt.Println("tasks interrupted")
	case ErrTimeout:
		fmt.Println("tasks timeout")
	default:
		fmt.Println("all tasks finished")
	}
	atomic.LoadInt64()
}
