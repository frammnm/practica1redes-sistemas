package main

import (
	"testing"
	"golang.org/x/crypto/ssh"
	"strings"
	"fmt"
	p1 "CodigoPractica/lib"
)

func TestSingle(t *testing.T){
	
	roles := []string{"C", "A", "B"}
	defaultAddresses := []string{"192.168.1.70:8083", "192.168.1.70:8081", "192.168.1.70:8082"}
	// defaultAddresses := []string{"155.210.154.197:17433", "155.210.154.200:17431", "155.210.154.199:17432"}
	dir := "/home/francisco/go/src/CodigoPractica"
	// dir := "/home/a794893/go/src/CodigoPractica"
	rsa := "/home/francisco/.ssh/id_rsa"
	// rsa := "/home/a794893/.ssh/id_rsa"

	i := 0
	for i < len(defaultAddresses) {
	 config := &ssh.ClientConfig {
	  // User: "a794893",
	  User: "francisco",
	  Auth: []ssh.AuthMethod{ p1.PublicKey(rsa)}, HostKeyCallback: ssh.InsecureIgnoreHostKey()}

	 conn, err := ssh.Dial("tcp", strings.Split(defaultAddresses[i], ":")[0] + ":22", config)
	 if err != nil {
	 	panic(err)
	 }

	 // Start Snode
	 fmt.Println("ssh to:", defaultAddresses[i])
	 go p1.RunCommand("cd " + dir + " && go run main.go " + roles[i] + " single", conn)
	 fmt.Println("cd " + dir + " && go run main.go " + roles[i] + " single")
	 // go p1.RunCommand("cd " + dir + " && /usr/local/go/bin/go run main.go " + roles[i] + " single", conn)
	 // fmt.Println("cd " + dir + " && /usr/local/go/bin/go run main.go " + roles[i] + " single")
	 i++
	 }

	 for{}
}


func TestMulti(t *testing.T){
	
	roles := []string{"C", "A", "B"}
	defaultAddresses := []string{"192.168.1.70:8083", "192.168.1.70:8081", "192.168.1.70:8082"}
	// defaultAddresses := []string{"155.210.154.197:17433", "155.210.154.200:17431", "155.210.154.199:17432"}
	dir := "/home/francisco/go/src/CodigoPractica"
	// dir := "/home/a794893/go/src/CodigoPractica"
	rsa := "/home/francisco/.ssh/id_rsa"
	// rsa := "/home/a794893/.ssh/id_rsa"
	i := 0
	for i < len(defaultAddresses) {
	 	config := &ssh.ClientConfig {
	  // User: "a794893",
	  User: "francisco",
	  Auth: []ssh.AuthMethod{ p1.PublicKey(rsa)}, HostKeyCallback: ssh.InsecureIgnoreHostKey()}

		conn, err := ssh.Dial("tcp", strings.Split(defaultAddresses[i], ":")[0] + ":22", config)
		if err != nil {
		 panic(err)
		}

		// Start Snode
		fmt.Println("ssh to:", defaultAddresses[i])
		go p1.RunCommand("cd " + dir + " && go run main.go " + roles[i] + " multi", conn)
		fmt.Println("cd " + dir + " && go run main.go " + roles[i] + " multi")
		// go p1.RunCommand("cd " + dir + " && /usr/local/go/bin/go run main.go " + roles[i] + " multi", conn)
		// fmt.Println("cd " + dir + " && /usr/local/go/bin/go run main.go " + roles[i] + " multi")
		i++
	 }

	 for{}
}


func TestSnap(t *testing.T){
	
	roles := []string{"C", "A", "B"}
	defaultAddresses := []string{"192.168.1.70:8083", "192.168.1.70:8081", "192.168.1.70:8082"}
	// defaultAddresses := []string{"155.210.154.197:17433", "155.210.154.200:17431", "155.210.154.199:17432"}
	dir := "/home/francisco/go/src/CodigoPractica"
	// dir := "/home/a794893/go/src/CodigoPractica"
	rsa := "/home/francisco/.ssh/id_rsa"
	// rsa := "/home/a794893/.ssh/id_rsa"

	i := 0
	for i < len(defaultAddresses) {
		config := &ssh.ClientConfig {
		// User: "a794893",
		User: "francisco",
		Auth: []ssh.AuthMethod{ p1.PublicKey(rsa)}, HostKeyCallback: ssh.InsecureIgnoreHostKey()}

		conn, err := ssh.Dial("tcp", strings.Split(defaultAddresses[i], ":")[0] + ":22", config)
		if err != nil {
			panic(err)
		}

		// Start Snode
		fmt.Println("ssh to:", defaultAddresses[i])
		go p1.RunCommand("cd " + dir + " && go run main.go " + roles[i] + " single/snap", conn)
		fmt.Println("cd " + dir + " && go run main.go " + roles[i] + " single/snap")
		// go p1.RunCommand("cd " + dir + " && /usr/local/go/bin/go run main.go " + roles[i] + " single/snap", conn)
		// fmt.Println("cd " + dir + " && /usr/local/go/bin/go run main.go " + roles[i] + " single/snap")
		i++
		}

		for{}
}


