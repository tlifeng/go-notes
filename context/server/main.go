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
	case <- time.After(time.Second * 2):
		_, err := fmt.Fprintln(w, "hello world!")
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("finish doing something")
	case <- ctx.Done():
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