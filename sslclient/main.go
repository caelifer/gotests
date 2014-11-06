package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var (
	host = flag.String("host", "", "Remote SSL host")
	port = flag.String("port", "443", "Remote SSL port")
)

func main() {
	flag.Parse()

	if *host == "" {
		fmt.Fprintln(os.Stderr, "Must provide remote host name")
		flag.Usage()
		os.Exit(1)
	}

	remoteHost := *host + ":" + *port

	quit := make(chan bool)
	defer close(quit)

	go prompt("Connecting to "+remoteHost+"... ", quit)
	conn, err := tls.Dial("tcp", remoteHost, &tls.Config{})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	quit <- true

	// Verify remote host
	if err = conn.VerifyHostname(*host); err != nil {
		log.Fatal("Failed host name verification")
	}

	// Set timeout
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// busy message
	go prompt("Sending request... ", quit)
	_, err = conn.Write([]byte("GET / HTTP/1.0\r\nHost: " + *host + "\r\n\r\n"))
	quit <- true
	if err != nil {
		log.Fatal(err)
	}

	// read response
	var nbytes int
	buf := make([]byte, 1024)

	// Actuall read loop
	for {
		// read
		n, err := conn.Read(buf)
		nbytes += n

		if n > 0 {
			fmt.Print(string(buf[:n])) // Always print what was read
		}

		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}
	}
	fmt.Printf("\n\nRead %d bytes\n", nbytes)
}

func prompt(p string, q chan bool) {
	fmt.Print(p)
	for {
		select {
		case <-time.Tick(time.Second):
			fmt.Print("+")
		case <-q:
			fmt.Println(" done")
			return
		}
	}
}
