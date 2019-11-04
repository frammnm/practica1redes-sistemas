package p1

import (
	"encoding/gob"
	"fmt"
	"github.com/dmichael/go-multicast/multicast"
	"log"
	"net"
	"sort"
	v "vclock"
	"time"
)

func SendGroupM(msg MsgI, defaultMulticastAddress string, addrs []string, resendLat int, port string) {

	conn, err := multicast.NewBroadcaster(defaultMulticastAddress)
	encoder := gob.NewEncoder(conn)
	if err != nil {
		log.Fatal(err)
	}
	
	//Listen ACKs
	addrChecked := make(map[string]bool)
	for _, v := range addrs {
		addrChecked[v] = false
	}
	
	fmt.Println("Receiving ACKs from:", addrs)
	go receiveACKs(msg, addrs, port, &addrChecked)
	

	for !TrueMap(addrChecked) {
		if err != nil {
			fmt.Println("Error encoding message", err.Error())
		} else {
	    done := true
	    for _, checked := range(addrChecked) {
				done = done && checked		   
			}
			if !done {
				if err != nil {
					log.Fatal(err)
				} else {
					// Encode and send data
					err = encoder.Encode(&msg)  	
				}
			}
			time.Sleep(time.Duration(resendLat) * time.Second)
		}
		encoder = gob.NewEncoder(conn)
	}
	// Close connection
	conn.Close()
}

func receiveACKs(msg MsgI, addrs []string, port string, addrChecked *map[string]bool) {

	addr, err := net.ResolveUDPAddr("udp", MyIp(port))
	if err != nil {
		log.Fatal(err)
	}

	// Open up a connection
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("recibiendo y escuchando acks en ", MyIp(port))

	decoder := gob.NewDecoder(conn)

	i := 0
	time.Now()
	
	for i < len(addrs) {
		// Decode data
		var res MsgI
		err = decoder.Decode(&res)
		if err == nil {
			switch v := res.(type) {
			case *Ack:
				//Exists in expected addresses and is meant to this node
				expectedAddr := stringInArray(v.From, addrs)
				fmt.Println("receiving ACK from:", v.From)

				if v.To == MyIp(port) && expectedAddr && !(*addrChecked)[v.From] {
					(*addrChecked)[v.From] = true
					i++
				} else {
					fmt.Println("Ack from an unexpected host, already checked, or not for me")
				}
			case *Msg:
				fmt.Println(v)
			case *Marker:
				fmt.Println(v)
			default:
				fmt.Printf("I don't know about type %T!\n", v)
			}
			decoder = gob.NewDecoder(conn)
		}
	}

	conn.Close()
}

func ReceiveGroupM(localClock *v.VClock, defaultMulticastAddress string, port string, messageList *[]MsgI) {
	fmt.Printf("Listening on %s\n", defaultMulticastAddress)

	// Parse the string address
	addr, err := net.ResolveUDPAddr("udp", defaultMulticastAddress)
	if err != nil {
		fmt.Println("Listen addr error: ", err.Error())
	}

	// Open up a connection
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	decoder := gob.NewDecoder(conn)

	// Loop forever reading from the socket
	for {
		// Decode data
		var msg MsgI
		err = decoder.Decode(&msg)
		fmt.Println("Received  Msg:", msg, "port:", port)
		if err != nil {
			fmt.Println("Error decoding message: ", err.Error())
			fmt.Println("Message received", msg)
		} else {
			switch val := msg.(type) {
			case *Ack:
				fmt.Println("Received ACK", val)
			case *Msg:
				if val.From != MyIp(port) {
					//Check if message exists in messageList
					if len(*messageList) > 0 {
						for _, m := range(*messageList) {
							if (!m.GetClock().Compare(val.Clock, v.Equal) && m.GetFrom() != val.From ){
								localClock.Merge(val.Clock)
								*messageList = append(*messageList, *val)
								sort.Sort(ByClock(*messageList))
								go sendACK(val.From, MyIp(port))
							}
						}
					} else {
						localClock.Merge(val.Clock)
						*messageList = append(*messageList, *val)
						sort.Sort(ByClock(*messageList))
						go sendACK(val.From, MyIp(port))
					}
				}
			default:
				fmt.Printf("I don't know about type %T!\n", val)
			}
		}
		time.Sleep(time.Duration(2) * time.Second)
		decoder = gob.NewDecoder(conn)
	}

	// Shut down the connection.
	conn.Close()
}

func sendACK(dest string, ori string) {

	raddr, err := net.ResolveUDPAddr("udp", dest)
	if err != nil {
		fmt.Println("Error resolving dest: ", err.Error())
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		fmt.Println("Dial error:", err.Error())
	}

	// Encode and send data
	var inter MsgI
	inter = Ack{nil, ori, dest}
	encoder := gob.NewEncoder(conn)
	fmt.Println("sending ACK to:", dest)
	err = encoder.Encode(&inter)

	conn.Close()
}