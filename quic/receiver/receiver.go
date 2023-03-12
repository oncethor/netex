package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/quic-go/quic-go"
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

	l, err := quic.ListenAddr(addrPort, generateTLSConfig(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	log.Println("Accepting connections on " + addrPort)
	for i := 1; ; i++ {
		conn, err := l.Accept(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		go func(conn quic.Connection, i int) {
			for j := 1; ; j++ {
				log.Printf("Accepting streams for (%v)", i)
				stream, err := conn.AcceptStream(context.Background())
				if err != nil {
					log.Println(err)
					return
				}

				go func(s quic.Stream, i, j int) {
					var buf = make([]byte, maxBuf)
					var n int
					var total int64
					var err error
					start := time.Now()
					log.Printf("(%v, %v) Started", i, j)
					for {
						if n, err = stream.Read(buf[:]); err != nil {
							log.Println(err)
							break
						}
						total += int64(n)
					}
					dur := time.Since(start)
					log.Printf("(%v, %v) Read %v bytes in %v (%.3f MiB/s)", i, j, total, dur, float64(total)/dur.Seconds()/(1024*1024))
					s.Close()
				}(stream, i, j)
			}
		}(conn, i)
	}
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-exchanger"},
	}
}
