package main

import (
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

// Get default location of a private key
func privateKeyPath() string {
	return os.Getenv("HOME") + "/.ssh/test.pem"
}

// Get private key for ssh authentication
func parsePrivateKey(keyPath string) (ssh.Signer, error) {
	buff, _ := os.ReadFile(keyPath)
	return ssh.ParsePrivateKey(buff)
}

// Get ssh client config for our connection
// SSH config will use 2 authentication strategies: by key and by password
func makeSshConfig(user string) (*ssh.ClientConfig, error) {
	key, err := parsePrivateKey(privateKeyPath())
	if err != nil {
		return nil, err
	}

	config := ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return &config, nil
}

// Handle local client connections and tunnel data to the remote serverq
// Will use io.Copy - http://golang.org/pkg/io/#Copy
func handleClient(client net.Conn, remote net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	// Start remote -> local data transfer
	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			log.Println("error while copy remote->local:", err)
		}
		chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			log.Println(err)
		}
		chDone <- true
	}()

	<-chDone
}

func main() {
	// Connection settings
	sshAddr := "localhost:22"
	localAddr := "127.0.0.1:5000"
	remoteAddr := "127.0.0.1:8080"

	user := os.Getenv("USER")

	// Build SSH client configuration
	cfg, err := makeSshConfig(user)
	if err != nil {
		log.Fatalln(err)
	}

	// Establish connection with SSH server
	conn, err := ssh.Dial("tcp", sshAddr, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Printf("connected to ssh host: %v", sshAddr)
	os.Exit(0)

	// Establish connection with remote server
	remote, err := conn.Dial("tcp", remoteAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connected to remote host: %v", remoteAddr)

	// Start local server to forward traffic to remote connection
	local, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer local.Close()
	log.Printf("created local socket on: %v", localAddr)

	// Handle incoming connections
	for {
		client, err := local.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		handleClient(client, remote)
	}
}

// vim: :ts=4:sw=4:noexpandtab:ai
