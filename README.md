# netex
Measure network speed (in Go)

This is a simple network speed measurer
    
receiver: ```go run receiver.go [port [buffer bytes]]```

sender:   ```go run sender.go [host:port [threads [total GiB to send [buffer bytes]]]]```

```go run sender.go myhost:12345 3 2``` will start 3 threads and send 2 GiB per each thread

First start the receiver, then run the sender. Default port is 12345.
