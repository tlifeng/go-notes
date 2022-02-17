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
	//ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	//ctx, cancel := context.WithDeadline(ctx, time.Date(2022, 2, 10, 8, 0, 0, 0, time.Local))

	go func() {
		time.Sleep(time.Second)
		cancel()
	}()
	job(ctx)
}
