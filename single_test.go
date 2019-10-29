package main

import (
	"testing"
	"golang.org/x/crypto/ssh"
	"strings"
	"fmt"
	p1 "CodigoPractica/lib"
)

func TestSingle(t *testing.T){
	
	defaultAddresses := []string{"192.168.1.70:8081", "192.168.1.70:8082", "192.168.1.70:8083"}
	roles := []string{"C", "A", "B"}
	dir := "/home/francisco/go/src/CodigoPractica"
	// dir := "/home/a794893/go/src/CodigoPractica1"

	i := 0
	for i < len(defaultAddresses) {
	 config := &ssh.ClientConfig {
	  // User: "a794893",
	  User: "francisco",
	  Auth: []ssh.AuthMethod{ p1.PublicKey("/home/francisco/.ssh/id_rsa")}, HostKeyCallback: ssh.InsecureIgnoreHostKey()}

	  // fmt.Println("address", defaultAddresses[i])

	 conn, err := ssh.Dial("tcp", strings.Split(defaultAddresses[i], ":")[0] + ":22", config)
	 if err != nil {
	 	panic(err)
	 }
 	 // defer conn.Close()

	 // Start Snode
	 fmt.Println("ssh to:", defaultAddresses[i])
	 go p1.RunCommand("cd " + dir + " && go run main.go " + roles[i] + " single", conn)
	 fmt.Println("cd " + dir + " && go run main.go " + roles[i] + " single")
	 i++
	 }

	 for{}
}


func TestMulti(t *testing.T){
	
	defaultAddresses := []string{"192.168.1.70:8081", "192.168.1.70:8082", "192.168.1.70:8083"}
	roles := []string{"C", "A", "B"}
	dir := "/home/francisco/go/src/CodigoPractica"
	// dir := "/home/a794893/go/src/CodigoPractica1"

	i := 0
	for i < len(defaultAddresses) {
	 config := &ssh.ClientConfig {
	  // User: "a794893",
	  User: "francisco",
	  Auth: []ssh.AuthMethod{ p1.PublicKey("/home/francisco/.ssh/id_rsa")}, HostKeyCallback: ssh.InsecureIgnoreHostKey()}

	  // fmt.Println("address", defaultAddresses[i])

	 conn, err := ssh.Dial("tcp", strings.Split(defaultAddresses[i], ":")[0] + ":22", config)
	 if err != nil {
	 	panic(err)
	 }
 	 // defer conn.Close()

	 // Start Snode
	 fmt.Println("ssh to:", defaultAddresses[i])
	 go p1.RunCommand("cd " + dir + " && go run main.go " + roles[i] + " multi", conn)
	 fmt.Println("cd " + dir + " && go run main.go " + roles[i] + " multi")
	 i++
	 }

	 for{}
}


func TestSnap(t *testing.T){
	
	defaultAddresses := []string{"192.168.1.70:8081", "192.168.1.70:8082", "192.168.1.70:8083"}
	roles := []string{"C", "A", "B"}
	dir := "/home/francisco/go/src/CodigoPractica"
	// dir := "/home/a794893/go/src/CodigoPractica1"

	i := 0
	for i < len(defaultAddresses) {
	 config := &ssh.ClientConfig {
	  // User: "a794893",
	  User: "francisco",
	  Auth: []ssh.AuthMethod{ p1.PublicKey("/home/francisco/.ssh/id_rsa")}, HostKeyCallback: ssh.InsecureIgnoreHostKey()}

	  // fmt.Println("address", defaultAddresses[i])

	 conn, err := ssh.Dial("tcp", strings.Split(defaultAddresses[i], ":")[0] + ":22", config)
	 if err != nil {
	 	panic(err)
	 }
 	 // defer conn.Close()

	 // Start Snode
	 fmt.Println("ssh to:", defaultAddresses[i])
	 go p1.RunCommand("cd " + dir + " && go run main.go " + roles[i] + " single/snap", conn)
	 fmt.Println("cd " + dir + " && go run main.go " + roles[i] + " single/snap")
	 i++
	 }

	 for{}
}


