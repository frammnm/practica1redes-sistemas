package p1

import (
	"encoding/gob"
	"fmt"
	"net"
	"sort"
	"strconv"
	"time"
	v "vclock"
	"strings"
)

func ReceiveGroup(localClock *v.VClock, addrs []string, port string, messageList *[]MsgI, messageMap *map[string][]MsgI) {

	//  Preparing to receive conncection
	ln, err := net.Listen("tcp", ":" + port)

	if err != nil {
		fmt.Println("Listening on port error: ", err.Error())
	}

	snapReq := false
	addrCheckedMarker := make(map[string]bool)
	for _, v := range addrs {
		addrCheckedMarker[v] = false
	}

	for {
		// Accept incoming connection
		conn, err := ln.Accept()
		fmt.Printf("\nAccepted connection from: [%v] g:%v  port:%v\n", conn.RemoteAddr().String(), addrs, port)
		expected := false
		for _, v := range(addrs) {
			if strings.Split(v, ":")[0] == strings.Split(conn.RemoteAddr().String(), ":")[0] {
				expected = true
				break
			}
		}

		if expected {
			fmt.Printf("Messanger belongs to group. Messanger:[%s] Group: [%v]\n", strings.Split(conn.RemoteAddr().String(), ":")[0], addrs)
		} else {
			fmt.Printf("Messanger does NOT belong to this group. Messanger:[%s] Group: [%v]\n", strings.Split(conn.RemoteAddr().String(), ":")[0], addrs)
		}

		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
		}
		// Decode data
		var msg MsgI
		decoder := gob.NewDecoder(conn)
		err = decoder.Decode(&msg)

		if err != nil {
			fmt.Printf("Error decoding message from: [%s]\n", conn.RemoteAddr().String())
			fmt.Println(err.Error())
		} else {

			switch val := msg.(type) {
			case *Msg:
				if !snapReq {
					fmt.Printf("Message received: **%v**\n", val)
					localClock.Merge(val.Clock)
					*messageList = append(*messageList, *val)
					sort.Sort(ByClock(*messageList))
				} else {
					fmt.Printf("Message received AFTER snapshot request: **%v**\n", val)
					localClock.Merge(val.Clock)
					(*messageMap)[val.From] = append((*messageMap)[val.From], *val)
					sort.Sort(ByClock((*messageMap)[val.From]))
				}
			case *Marker:
				if !snapReq {
					fmt.Println("RECEIVED FIRST MARKER  ", addrCheckedMarker, val.From)
					addrCheckedMarker[val.From] = true
					var copyList []MsgI
					copy(copyList, *messageList)
					val.State = copyList
					val.From = MyIp(port)
					lats := []string{"1", "1"}
					go SendGroup(val, addrs, lats)
					snapReq = true
				} else {
					if uncheck, ok := addrCheckedMarker[val.From]; ok && !uncheck{
						fmt.Println("RECEIVED SECOND MARKER  ", addrCheckedMarker, val.From)
						if addrCheckedMarker[val.From] == false {
							addrCheckedMarker[val.From] = true
						}
						if TrueMap(addrCheckedMarker){
							fmt.Printf("All markers received in port: [%v]\n", port)
							fmt.Printf("Snapshot of the system for [%v]:\n MyState: %v \n AfterMarker: %v\n", port, messageList, messageMap)
						}
					}
				}
			default:
				fmt.Printf("I don't know about type %T!\n", val)
			}
		}
		// Shut down the connection.
		conn.Close()
	}
}

func SendToNode(msg MsgI, addr string, lat string) {

	// Simulate latency
	intLat, _ := strconv.Atoi(lat)
	time.Sleep(time.Duration(intLat) * time.Second)

	// Connect to the socket described by addr
	trying := 0
	for trying < 3 {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Println("Dial error, retrying..:", err.Error())
		} else {
			// Encode and send data
			encoder := gob.NewEncoder(conn)
			err = encoder.Encode(&msg)
			trying = 3
			// Close connection
			conn.Close()
		}
		trying++
	}

}

func SendGroup(msg MsgI, addrs []string, lats []string) {
	for i := range addrs {
		go SendToNode(msg, addrs[i], lats[i])
	}
}


func ChandyLamport(msg MsgI, addrs []string, lats []string) {
	fmt.Println("******************************************")
	fmt.Println("Starting Snapshot")
	fmt.Println("******************************************")
	SendGroup(msg, addrs, lats)
}