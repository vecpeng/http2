package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dgrr/http2/fasthttp2"
	"github.com/valyala/fasthttp"
)

func main() {
	c := &fasthttp.HostClient{
		Addr:  "api.binance.com:443",
		IsTLS: true,
	}

	if err := fasthttp2.ConfigureClient(c); err != nil {
		panic(err)
	}

	count := int32(0)
	var wg sync.WaitGroup
	for i := 0; i < 1; i++ {
		for atomic.LoadInt32(&count) >= 4 {
			time.Sleep(time.Millisecond * 100)
		}

		wg.Add(1)
		atomic.AddInt32(&count, 1)
		go func() {
			defer wg.Done()
			defer atomic.AddInt32(&count, -1)

			req := fasthttp.AcquireRequest()
			res := fasthttp.AcquireResponse()

			res.Reset()

			req.Header.SetMethod("GET")
			// TODO: Use SetRequestURI
			req.URI().Update("https://api.binance.com/api/v3/exchangeInfo")

			err := c.Do(req, res)
			if err != nil {
				log.Fatalln(err)
			}

			body := res.Body()

			fmt.Printf("%d: %d\n", res.Header.StatusCode(), len(body))
			res.Header.VisitAll(func(k, v []byte) {
				fmt.Printf("%s: %s\n", k, v)
			})

			a := make(map[string]interface{})
			if err = json.Unmarshal(body, &a); err != nil {
				panic(err)
			}

			fmt.Println("------------------------")
		}()
	}

	wg.Wait()

	// fmt.Printf("%s\n", body)
}
