package main

import (
	"net"
	"fmt"
	"time"
	"strings"
	"os"
)

var portToUse = 23001
var listnerUsed = false

func write(pc net.PacketConn) {

  addr,err := net.ResolveUDPAddr("udp4", "192.168.86.255:8829")
  if err != nil {
    panic(err)
  }

	for {
		msg := fmt.Sprintf("From:%s:%d", os.Getenv("THOST"), portToUse)
		_,err = pc.WriteTo([]byte(msg), addr)
		if err != nil {
			panic(err)
		}
		time.Sleep(10 * time.Second)
	}
}

var myAddress string = ""

func updateConfig (addr net.Addr, b []byte) {

	// check to see if tpbridge connection exists

	// else add new one.
	m := strings.Split(myAddress, ":")
	a := strings.Split(addr.String(), ":")
	t := strings.Split(string(b), ":")
	if a[0] < m[0] {
		fmt.Printf("%s < %s : %s\n", a[0], m[0], b)
		fmt.Printf("-f -n %s -i %s -p %s -tA -z30\n", t[1],
			a[0], t[2])
	} else {
		fmt.Printf("%s > %s: %s\n", a[0], m[0], b)
		fmt.Printf("-f -n %s -i %s -p %d -tP -z30\n", t[1],
			m[0], portToUse)
		portToUse++
		
	}
		


}

func main() {
	pc,err := net.ListenPacket("udp4", ":8829")
  if err != nil {
    panic(err)
  }
  defer pc.Close()

  go write(pc)

	buf := make([]byte, 1024)
	nodes := make(map[string]int)
	
	for {
		n,addr,err := pc.ReadFrom(buf)
		if err != nil {
			panic(err)
		}
		s := strings.Split(string(buf[:n]), ":")
		//fmt.Printf("%s sent this: %s\n", addr, buf[:n])
		if s[1] == os.Getenv("THOST") {
			//fmt.Println("Skipping self.")
			myAddress = addr.String()
		} else {
			_, found := nodes[s[1]]
			if ! found {
				updateConfig(addr, buf[:n])
				nodes[s[1]] = 0
			}
		}
	}
}
