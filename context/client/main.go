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

	req, errReq := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080", nil)
	if errReq != nil {
		log.Fatal(errReq)
	}
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
