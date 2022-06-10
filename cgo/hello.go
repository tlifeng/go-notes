//package main
//
///*
//#include <hello.h>
//*/
//import "C"
//
//func main() {
//	C.SayHello(C.CString("hello world\n"))
//}

package main

//void SayHello(_GoString_ s);
import "C"

import (
	"fmt"
	"time"
)

func main() {
	C.SayHello("1\n")
	C.SayHello("2\n")
	C.SayHello("3\n")
	C.SayHello("4\n")
	C.SayHello("5\n")
	fmt.Println(6)
	fmt.Println(7)
	fmt.Println(8)
	time.Sleep(1 * time.Second)
}

//export SayHello
func SayHello(s string) {
	fmt.Print(s)
}
