package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	log.Printf("example usage: host:port 3 5 16384 (host:port, threads, GiB, buffer bytes)")
	wg := &sync.WaitGroup{}
	host := "localhost:12345"
	threads := 1
	gigs := 10
	maxBuf := 16 * 1024
	if len(os.Args) > 1 {
		host = os.Args[1]
		if len(os.Args) > 2 {
			var err error
			if threads, err = strconv.Atoi(os.Args[2]); err != nil {
				log.Fatalf("2nd argument (threads) should be an integer: %v", os.Args[2])
			}
			if len(os.Args) > 3 {
				if gigs, err = strconv.Atoi(os.Args[3]); err != nil {
					log.Fatalf("3rd argument (GB) should be an integer: %v", os.Args[3])
				}
				if len(os.Args) > 4 {
					if maxBuf, err = strconv.Atoi(os.Args[4]); err != nil {
						log.Fatalf("4th argument (buffer bytes) should be an integer: %v", os.Args[4])
					}
				}
			}
		}
	}
	log.Printf("using %v total 'threads'", threads)
	log.Printf("using %v total gigabytes per thread", gigs)
	log.Printf("using %v bytes buffer", maxBuf)
	var buf = make([]byte, maxBuf)
	for i := 0; i < maxBuf; i++ {
		buf[i] = byte(rand.Intn(256))
	}
	log.Printf("using %v bytes buffer", maxBuf)
	for t := 0; t < threads; t++ {
		wg.Add(1)
		go func(c int) {
			var n, total int
			log.Printf("(%v) Connecting to %v", c, host)
			conn, err := d.DialContext(ctx, "tcp", host)
			if err != nil {
				log.Fatalf("Failed to dial: %v", err)
			}
			defer conn.Close()

			if _, err = conn.Write([]byte("Hello, World!")); err != nil {
				log.Fatal(err)
			}
			for j := 0; j < gigs*1024*1024*1024/maxBuf; j++ {
				if n, err = conn.Write(buf[:]); err != nil {
					log.Fatal(err)
				}
				total += n
			}
			log.Printf("(%v) Written %v bytes", c, total)
			wg.Done()
		}(t + 1)
	}
	wg.Wait()
}
