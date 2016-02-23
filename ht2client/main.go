package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/http2"
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
	fmt.Println("[GO-HT2CLIENT] Connecting to the remote host", remoteHost, "...")

	client := http.Client{
		Transport: &http2.Transport{},
	}

	resp, err := client.Get("https://" + remoteHost + "/")
	handleError(err)
	defer resp.Body.Close()

	_, err = io.Copy(os.Stdout, resp.Body)
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
