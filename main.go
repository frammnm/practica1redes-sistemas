package main

import (
	"fmt"
	"os"
	"time"
	p1 "CodigoPractica/lib"
	"github.com/urfave/cli"
	v "vclock"
)


func main() {

	app := cli.NewApp()

	var mode string
	var role string
	var localClock v.VClock
	localClock = v.VClock(map[string]uint64{"192.168.1.70:8081": 0, "192.168.1.70:8082": 0, "192.168.1.70:8083": 0})
	// localClock = v.VClock(map[string]uint64{"155.210.154.200:17431": 0, "155.210.154.199:17432": 0, "155.210.154.197:17433": 0})
	defaultAddresses := []string{"192.168.1.70:8081", "192.168.1.70:8082", "192.168.1.70:8083"}
	// defaultAddresses := []string{"155.210.154.200:17431", "155.210.154.199:17432", "155.210.154.197:17433"}
	defaultMulticastAddress := "239.0.0.0:9999"
	// defaultMulticastAddress := "239.0.074.003:9999"

	portA := "8081"
	portB := "8082"
	portC := "8083"

	// portA := "17431"
	// portB := "17432"
	// portC := "17433"
	
	var messageListA []p1.MsgI
	var messageListB []p1.MsgI
	var messageListC []p1.MsgI
	
	messageMapA := make(map[string][]p1.MsgI)
	messageMapB := make(map[string][]p1.MsgI)
	messageMapC := make(map[string][]p1.MsgI)


	app.Action = func(c *cli.Context) error {

		role = c.Args()[0]
		mode = c.Args()[1]

		return nil
	}

	app.Run(os.Args)

	switch role {
			case "A":
				//Comportamiento del nodo iniciador
				if mode == "single" || mode == "single/snap" {
					go p1.ReceiveGroup(&localClock, defaultAddresses[1:], portA, &messageListA, &messageMapA)
		    	
		    	lats := []string{"5", "10"}
		    	localClock.Tick(defaultAddresses[0])
					msg := p1.New(localClock, "prueba", "192.168.1.70:8081", "192.168.1.70:8082-192.168.1.70:8083")
					// msg := p1.New(localClock, "prueba", "155.210.154.200:17431", "155.210.154.200:17432-155.210.154.200:17433")
					go p1.SendGroup(msg, defaultAddresses[1:], lats)

		    	for {
		    		fmt.Println("*********************************")
		    		fmt.Println("\nOrdered message list for node A: ", messageListA)
		    		fmt.Printf("Port for [%v]: \n AfterMarker: %v\n", portA, messageMapA)
		    		fmt.Println("Clock for node A: ", localClock)
		    		fmt.Println("*********************************")
		    		time.Sleep(time.Duration(5) * time.Second)
		    	}

				} else {
					go p1.ReceiveGroupM(&localClock, defaultMulticastAddress, portA, &messageListA)

					time.Sleep(time.Duration(7) * time.Second)
					localClock.Tick(defaultAddresses[0])
					msg := p1.New(localClock, "prueba", "192.168.1.70:8081", "192.168.1.70:8082-192.168.1.70:8083")
					// msg := p1.New(localClock, "prueba", "155.210.154.200:17431", "155.210.154.200:17432-155.210.154.200:17433")
					go p1.SendGroupM(msg, defaultMulticastAddress, defaultAddresses[1:], 1, portA)

					for {
		    		fmt.Println("*********************************")
		    		fmt.Println("\nOrdered message list for node A: ", messageListA)
		    		fmt.Println("Clock for node A: ", localClock)
		    		fmt.Println("*********************************")
		    		time.Sleep(time.Duration(6) * time.Second)
		    	}
				}
			case "B":
				//Comportamiento de receptor que reenvia
				if mode == "single" || mode == "single/snap" {
					tmp := make([]string, len(defaultAddresses))
					copy(tmp, defaultAddresses)
					go p1.ReceiveGroup(&localClock, append(tmp[:1], defaultAddresses[2]), portB, &messageListB, &messageMapB)

					time.Sleep(time.Duration(6) * time.Second)
					lats := []string{"6", "6"}
		    	localClock.Tick(defaultAddresses[1])
					msg := p1.New(localClock, "prueba", "192.168.1.70:8082", "192.168.1.70:8083-192.168.1.70:8081")
					// msg := p1.New(localClock, "prueba", "155.210.154.200:17432", "155.210.154.200:17433-155.210.154.200:17431")
					tmp = make([]string, len(defaultAddresses))
					copy(tmp, defaultAddresses)
					go p1.SendGroup(msg, append(tmp[:1], defaultAddresses[2]), lats)

					for {
						fmt.Println("*********************************")
		    		fmt.Println("\nOrdered message list for node B: ", messageListB)
		    		fmt.Printf("Port for [%v]: \n AfterMarker: %v\n", portB, messageMapB)
		    		fmt.Println("Clock for node B: ", localClock)
		    		fmt.Println("*********************************")
		    		time.Sleep(time.Duration(5) * time.Second)
		    	}

				} else {

					go p1.ReceiveGroupM(&localClock, defaultMulticastAddress, portB, &messageListB)
					
					time.Sleep(time.Duration(10) * time.Second)
					localClock.Tick(defaultAddresses[1])
					msg := p1.New(localClock, "prueba", "192.168.1.70:8082", "92.168.1.70:8083-192.168.1.70:8081")
					// msg := p1.New(localClock, "prueba", "155.210.154.200:17432", "155.210.154.200:17433-155.210.154.200:17431")
					tmp := make([]string, len(defaultAddresses))
					copy(tmp, defaultAddresses)
					go p1.SendGroupM(msg, defaultMulticastAddress, append(tmp[:1], defaultAddresses[2]), 1, portB)


					for {
		    		fmt.Println("*********************************")
		    		fmt.Println("\nOrdered message list for node B: ", messageListB)
		    		fmt.Println("Clock for node B: ", localClock)
		    		fmt.Println("*********************************")
		    		time.Sleep(time.Duration(6) * time.Second)
		    	}
					
				}
			case "C":
				//Comportamiento de receptor que solo escucha
				if mode == "single" || mode == "single/snap" {
					go p1.ReceiveGroup(&localClock, defaultAddresses[:2], portC, &messageListC, &messageMapC)

					time.Sleep(time.Duration(20) * time.Second)
					lats := []string{"2", "2"}


					start := time.Now()
					now := time.Now()
					once := false

					for {

						if mode == "single/snap" && (now.Sub(start) > (8 * time.Second)) && !once {
							msg := p1.Marker{localClock, defaultAddresses[2], "192.168.1.70:8081-192.168.1.70:8082", messageListC}
							// msg := p1.Marker{localClock, defaultAddresses[2], "155.210.154.200:17431-155.210.154.200:17432", messageListC}
							p1.ChandyLamport(msg, defaultAddresses[:2], lats)
							once = true
						}

						//recalculate state after 20 secs
						// if now.Sub(start) > (20 * time.Second) {
						// 	start = time.Now()
						// 	now = time.Now()
						// 	once = false
						// }

						fmt.Println("*********************************")
		    		fmt.Println("\nOrdered message list for node C: ", messageListC)
		    		fmt.Println("Clock for node C: ", localClock)
		    		fmt.Println("*********************************")
		    		time.Sleep(time.Duration(3) * time.Second)
		    		localClock.Tick(defaultAddresses[2])
						msg := p1.New(localClock, "prueba", defaultAddresses[2], "192.168.1.70:8081-192.168.1.70:8082")
						// msg := p1.New(localClock, "prueba", defaultAddresses[2], "155.210.154.200:17431-155.210.154.200:17432")
		    		go p1.SendGroup(msg, defaultAddresses[:2], lats)
		    		now = time.Now()
		    	}

				} else {
					go p1.ReceiveGroupM(&localClock, defaultMulticastAddress, portC, &messageListC)
					time.Sleep(time.Duration(20) * time.Second)

					for {
		    		fmt.Println("*********************************")
		    		fmt.Println("\nOrdered message list for node C: ", messageListC)
		    		fmt.Println("Clock for node C: ", localClock)
		    		fmt.Println("*********************************")
		    		time.Sleep(time.Duration(6) * time.Second)
		    		localClock.Tick(defaultAddresses[2])
						msg := p1.New(localClock,  "prueba", defaultAddresses[2], "192.168.1.70:8081-192.168.1.70:8082")
						// msg := p1.New(localClock,  "prueba", defaultAddresses[2], "155.210.154.200:17431-155.210.154.200:17432")
						go p1.SendGroupM(msg, defaultMulticastAddress, defaultAddresses[:2], 3, portC)
		    	}
				}
			default:
				fmt.Printf("What role is this? -> [%s]\n", role)
	}

}