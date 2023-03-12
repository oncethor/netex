package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	port := 12345
	maxBuf := 16 * 1024
	if len(os.Args) > 1 {
		var err error
		if port, err = strconv.Atoi(os.Args[1]); err != nil {
			log.Fatalf("1st argument (port) should be an integer: %v", os.Args[1])
		}
		if len(os.Args) > 2 {
			if maxBuf, err = strconv.Atoi(os.Args[4]); err != nil {
				log.Fatalf("2nd argument (buffer bytes) should be an integer: %v", os.Args[2])
			}
		}
	}
	log.Printf("using %v port", port)
	log.Printf("using %v bytes buffer", maxBuf)

	addrPort := fmt.Sprintf(":%v", port)
	l, err := net.Listen("tcp", addrPort)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	log.Println("Accepting connections on " + addrPort)
	for i := 1; ; i++ {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn, i int) {
			var buf = make([]byte, maxBuf)
			var n int
			var total int64
			var err error
			start := time.Now()
			log.Printf("(%v) Started", i)
			for {
				if n, err = c.Read(buf[:]); err != nil {
					log.Println(err)
					break
				}
				total += int64(n)
			}
			dur := time.Since(start)
			log.Printf("(%v) Read %v bytes in %v (%.3f MiB/s)", i, total, dur, float64(total)/dur.Seconds()/(1024*1024))
			c.Close()
		}(conn, i)
	}
}
