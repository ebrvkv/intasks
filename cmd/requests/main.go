package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ebrvkv/intasks/internal/counter"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const initialAmountReq = 1

var (
	client                      = http.DefaultClient
	cnt                         = counter.NewReqCounter()
	url                         string
	multiplier, timeout, period int
)

func makeRequests(ctx context.Context, wg *sync.WaitGroup, amount int, reqChan chan int64) {
	wg.Add(1)
	defer wg.Done()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		panic(err)
	}
	t := time.NewTicker(time.Duration(period) * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			amount *= multiplier
			for i := 0; i < amount; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					cnt.Inc()
					resp, err := client.Do(req)
					cnt.Reduce()
					if err != nil && strings.Contains(err.Error(), "Timeout") && !cnt.Stopped() {
						cnt.Stop()
						reqChan <- cnt.Get()
					}
					if resp != nil {
						if err := resp.Body.Close(); err != nil {
							log.Println(err)
						}
					}
				}()
			}
		}
	}
}

func init() {
	flag.StringVar(&url, "url", "https://ya.ru", "url to which GET requests will be sent")
	flag.IntVar(&multiplier, "m", 2, "int value by witch we will multiply amount of requests made on previous iteration")
	flag.IntVar(&timeout, "t", 100, "timeout in milliseconds from net.Dialer till end of response from remote end")
	flag.IntVar(&period, "p", 1000, "how often in milliseconds we need to increase amount of HTTP requests")

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 1000
	t.MaxConnsPerHost = 1000
	t.MaxIdleConnsPerHost = 1000
	client = &http.Client{
		Transport: t,
	}
}

func main() {
	flag.Parse()
	if url == "" {
		fmt.Println("URL can't be empty")
		return
	}

	client.Timeout = time.Duration(timeout) * time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())
	reqChan := make(chan int64)

	wg := &sync.WaitGroup{}

	go makeRequests(ctx, wg, initialAmountReq, reqChan)

	count := <-reqChan
	cancel()
	wg.Wait()
	fmt.Printf("Amount of parallel/concurrent requests on exit: %d\n", count)
}
