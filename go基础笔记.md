# Goroutine

1. 常见的并发模型： 多线程和消息传递，goroutine是基于CSP的消息传递模型
2. goroutine能动态伸缩栈空间，2kb到1GB。相比线程2MB开销小很多
3. 同一个goroutine中顺序一致性的内存模型能得到保证，不同的goroutine不能保证，需要定义明确的同步事件来保证
4. 控制一个或多个goroutine退出，一是close(channel)搭配select case <- channel 一起用，二是context cancel方案和select case <- ctx.Done()一起用
```go
//这个例子mu的加锁和解锁不在同一个goroutine，所以可能会出错，可能先执行了main goroutine中的unlock
func main() {
    var mu sync.Mutex

    go func(){
        fmt.Println("你好, 世界")
        mu.Lock()
    }()

    mu.Unlock()

```

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

1. 无缓存的channel每次发送操作都有对应的接送动作，不然会死锁
2. 无缓存channel在同一goroutine使用发送和接收操作，容易造成死锁
3. 无缓存的Channel上的发送操作总在对应的接收操作完成前发生
4. 可以通过带缓存channel的大小来控制并发量
5. 可以通过带缓存Channel的使用量和最大容量比例来判断程序运行的并发率。当管道为空的时候可以认为是空闲状态，当管道满了时任务是繁忙状态，这对于后台一些低级任务的运行是有参考价值的
6. 对于带缓冲的Channel，对于Channel的第K个接收完成操作发生在第K+C个发送操作完成之前，其中C是Channel的缓存大小
```go
//利用带缓存channel控制并发量
var limit = make(chan int, 3)

func main() {
    for _, w := range work {
        go func() {
            limit <- 1
            w()
            <-limit
        }()
    }
    select{}
}
```

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

## 关闭channel
1. 向一个已经关闭的通道发送值是不允许的,会报错
2. 从一个已经关闭但是里面还有值的通道取值是允许的,可以正常获取到值
3. 从一个已经关闭但是为空的通道取值是允许的,会获取通道类型元素的零值
4. 不可以再次关闭一个已经关闭的通道,会报错
5. 已经关闭的通道无法再次打开

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

