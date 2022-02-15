# Goroutine

可以使用`sync.WaitGroup`等待go协程的运行

```go
wg.Add(num)	// 添加num个任务

wg.Done()	// 响应一个任务完成

wg.Wait()	//等待任务
```



使用事例:

```go
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

/*
输出:
Hello
I am done! id: world
exit

*/
```



# Channel

使用goroutine会产生竞者条件(race condition), 解决的三种办法

1. 使用atomic包的方法
2. 使用排他锁mutex
3. 通过管道获取值

**注意： 无缓存channel要获取完所有值**

例子:

```go
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

	go ChannelIntCounter()
	go ChannelIntCounter()

	ch <- 0 //这里使用有缓存的channel,否则最后一个值没拿会发生死锁。假如是一个无缓存channel，ch赋值不能放在go程序前面，否则会一直阻塞发生死锁
	wg.Wait()

	fmt.Println(<-ch)
}
```



# Context

创建context首先创建空的context

```go
ctx := context.Background()
ctx := context.TODO() //不知道用不用时候可以使用todo创建
```

控制context的三种方法

```go
//context的三种创建方式，第二第三种不用手动取消
ctx, cancel := context.WithCancel(ctx)
ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
ctx, cancel := context.WithDeadline(ctx, time.Date(2022, 2, 10, 8, 0, 0, 0, time.Local))
```

例子：

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func job(ctx context.Context) {
	//doing something
	select {
	case <-ctx.Done(): //要通过select，ctx.done来知道ctx已经结束
		errContext := ctx.Err()
		if errContext != nil {
			fmt.Println(errContext)
		}
	case <-time.After(2 * time.Second):
		fmt.Println("job completed!")
	}
}

func main() {
	ctx := context.Background()
	//ctx := context.TODO() //不知道用不用时候可以使用todo创建

	//context的三种创建方式，第二第三种不用手动取消
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		time.Sleep(time.Second)
		cancel()
	}()
	job(ctx)
}

```

## http使用context

server端代码：

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("handler start")

	ctx := r.Context()

	select {
	case <- time.After(time.Second * 2):	//两秒后返回hello world!
		_, err := fmt.Fprintln(w, "hello world!")
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("finish doing something")
	case <- ctx.Done():	//客户端context断开，输出错误原因
		errCtx := ctx.Err()
		if errCtx != nil {
			fmt.Println(errCtx)
		}
	}

	fmt.Printf("handler end\n\n")
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

client端代码:

```go
package main

import (
	"context"
	"io/ioutil"
	"log"
	http "net/http"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*2)
	defer cancelFunc()

    //创建一个带context的请求
	req, errReq := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080", nil)
	if errReq != nil {
		log.Fatal(errReq)
	}
    //发送请求
	resp, errResp := http.DefaultClient.Do(req)
	if errResp != nil {
		log.Fatal(errResp)
	}

	defer resp.Body.Close()

	respBytes, errReadResp := ioutil.ReadAll(resp.Body)
	if errReadResp != nil {
		log.Fatal(errReadResp)
	}
	log.Fatal(string(respBytes))
}

```

