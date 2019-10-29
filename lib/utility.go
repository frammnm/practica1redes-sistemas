package p1

import (
	"encoding/gob"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
	"os"
	v "vclock"
)

type MsgI interface {
	SetFrom(string)
	GetClock() v.VClock
	GetFrom() string
}

type Msg struct {
	Clock  v.VClock
	Body   string
	From   string
	To     string
}

type Ack struct {
	Clock v.VClock
	From  string
	To    string
}

type Marker struct {
	Clock v.VClock
	From  string
	To    string
	State []MsgI
}

type ByClock []MsgI

func (a ByClock) Len() int      { return len(a) }
func (a ByClock) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByClock) Less(i, j int) bool {
	return a[i].GetClock().Compare(a[j].GetClock(), v.Descendant)
}

//New returns a new Msg
func New(clock v.VClock, b string, from string, to string) Msg {
	return Msg{clock, b, from, to}
}

func (m Msg) GetClock() v.VClock { return m.Clock }
func (m Msg) SetFrom(s string)   { m.From = s }
func (m Msg) GetFrom() string  { return m.From }

func (m Ack) GetClock() v.VClock { return m.Clock }
func (m Ack) SetFrom(s string)   { m.From = s }
func (m Ack) GetFrom() string  { return m.From }

func (m Marker) GetClock() v.VClock { return m.Clock }
func (m Marker) SetFrom(s string)   { m.From = s }
func (m Marker) GetFrom() string  { return m.From }

func stringInArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func init() {
	gob.Register(&Msg{})
	gob.Register(&Ack{})
	gob.Register(&Marker{})
}

func RunCommand(cmd string, conn *ssh.Client) {
	sess, err := conn.NewSession()
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	defer sess.Close()
	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stdout, sessStdOut)
	sessStderr, err := sess.StderrPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stderr, sessStderr)
	err = sess.Run(cmd) // eg., /usr/bin/whoami
	if err != nil {
		panic(err)
	}
}

func PublicKey(path string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
}

func MyIp(port string) string {

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
	}

	var currentIP string

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				currentIP = ipnet.IP.String()
			}
		}
	}

	return currentIP + ":" + port
}

func RemoveFromSlice(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func TrueMap(m map[string]bool) bool {
	for _, v := range(m) {
	  if !v {
	    return false
	  }
	}
	return true
}
