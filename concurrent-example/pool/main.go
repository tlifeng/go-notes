package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var (
	counter int32
	wg sync.WaitGroup
)

// 定义一个DBConnection资源
type DBConnection struct {
	id int32
}

func (D DBConnection) Close() error {
	fmt.Println("database closed, #" + fmt.Sprint(D.id))
	return nil
}

func Factory() (io.Closer, error) {
	atomic.AddInt32(&counter, 1)
	return DBConnection{
		id: counter,
	}, nil
}

func (pool Pool) performQuery(query int)  {
	defer wg.Done()
	
	resource, err := pool.AcquireResource()
	if err != nil {
		log.Fatal(err)
	}
	defer pool.ReleaseResource(resource)
	
	t := rand.Int()%10 + 1
	time.Sleep(time.Duration(t) * time.Second)
	fmt.Println("finish query" + fmt.Sprint(query))
}

func main() {
	p, err := New(Factory, 5)
	if err != nil {
		log.Fatal(err)
	}
	
	num := 10
	wg.Add(num)
	for id := 0; id < num; id++ {
		go p.performQuery(id)
	}
	wg.Wait()

	p.Close()
	fmt.Println("pool model run ends")
}
