// Package main TODO
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	u := flag.String("url", "", "URL do serviço a ser testado.")
	r := flag.Int("requests", 1, "Número total de requests.")
	c := flag.Int("concurrency", 1, "Número de chamadas simultâneas.")

	flag.Parse()

	if *u == "" {
		log.Fatal("url is required")
	}

	_, err := http.NewRequest(http.MethodGet, *u, nil)
	if err != nil {
		log.Fatal("creating request:", err.Error())
	}

	sig := make(chan struct{}, *c)
	res := make(chan int, *r)
	ers := make(chan error, *r)

	wg := sync.WaitGroup{}
	wg.Add(*r)

	start := time.Now()
	for range *r {
		sig <- struct{}{}
		go func() {
			defer func() { <-sig }()
			defer wg.Done()

			// não é possível retornar erro aqui
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, *u, nil)

			response, err := http.DefaultClient.Do(req)
			if err != nil {
				ers <- err
				return
			}
			defer response.Body.Close()
			res <- response.StatusCode
		}()
	}

	wg.Wait()
	close(res)

	delta := time.Since(start).String()

	codes := make(map[int]int, len(res))
	for code := range res {
		codes[code]++
	}

	successes := codes[http.StatusOK]
	delete(codes, http.StatusOK)

	fmt.Printf("Delta: %s Total: %d Errs:%d StatusOk: %d Misc.: %+v\n", delta, *r, len(ers), successes, codes)

}
