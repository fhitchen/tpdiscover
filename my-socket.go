package main

import (
	"encoding/xml"
	"io/ioutil"
	"net"
	"fmt"
	"time"
	"strings"
	"strconv"
	"os"
	"os/exec"
        log "github.com/sirupsen/logrus"	
)

var portToUse = 23000
var srvid = 160

func getDefaultBroadcastAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if v, ok := address.(*net.IPNet); ok && !v.IP.IsLoopback() {
			if v.IP.To4() != nil {
				ip := v.IP
				mask := v.Mask
				ip = ip.To4()
				fmt.Printf("ip, mask: %v, %v\n", ip, mask)
				return fmt.Sprintf("%d.%d.%d.%d", ip[0] | ^mask[0], ip[1] | ^mask[1], ip[2] | ^mask[2], ip[3] | ^mask[3])
			}
		}
	}
	return "255.255.255.255"
}



func write(pc net.PacketConn, bcast string) {

	//addr,err := net.ResolveUDPAddr("udp4", "192.168.86.255:8829")
	addr,err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:8829", bcast))
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

func addTPBridge(conf Endurox, node, ip, mode string, port int) (c Endurox) {

	appopt := fmt.Sprintf("-f -r -n %s -i %s -p %d -t%s -z30",
		node, ip, port, mode)
	
	v:= Server{
		Comment: "Added by TPDISCOVER service.",
		Name: "tpbridge",
		Min: "1",
		Max: "1",
		Srvid: fmt.Sprintf("%d", srvid),
		Sysopt: fmt.Sprintf("-e ${NDRX_APPHOME}/log/tpbridge_%s.log -r", node),
		Appopt: appopt,
	}

	fmt.Printf("V := %#v\n",v)

	conf.Servers.Server = append(conf.Servers.Server, v)

	file, _ := xml.MarshalIndent(conf, "", " ")
 
	_ = ioutil.WriteFile("ndrxconfig.xml", file, 0644)

	/*enc := xml.NewEncoder(os.Stdout)
	enc.Indent("  ", "    ")
	if err := enc.Encode(conf); err != nil {
		fmt.Printf("error: %v\n", err)
	}*/

	cmd := exec.Command("/bin/sh", "-c", "xadmin reload")

	_, err := cmd.Output()
        if err != nil {
                log.Fatal(fmt.Errorf("cmd.Run(1) failed with %w\n", err))
        }
	
	xadminCmd := fmt.Sprintf("xadmin start -i %d", srvid)

	cmd = exec.Command("/bin/sh", "-c", xadminCmd)

	_, err = cmd.Output()
        if err != nil {
		log.Fatal(fmt.Errorf("cmd.Run(2) failed with %w\n", err))
        }

	srvid += 10
	
	return conf	

}

func updateConfig (addr net.Addr, b []byte, conf Endurox) (c Endurox) {

	// check to see if tpbridge connection exists

	// else add new one.
	m := strings.Split(myAddress, ":")
	a := strings.Split(addr.String(), ":")
	t := strings.Split(string(b), ":")

	if os.Getenv("THOST") < t[1] {
		node, _ := strconv.Atoi(t[1])
		lnode, _ := strconv.Atoi(os.Getenv("THOST"))
		fmt.Printf("%s < %s : %s\n", os.Getenv("THOST"), t[1], b)
		fmt.Printf("-f -n %s -i %s -p %d -tA -z30\n", t[1],
			a[0], portToUse + (lnode * 100) + node)
		conf = addTPBridge(conf, t[1], a[0], "A", portToUse + (lnode * 100) + node)
	} else {
		node, _ := strconv.Atoi(os.Getenv("THOST"))
		lnode, _ := strconv.Atoi(t[1])
		fmt.Printf("%s > %s: %s\n", os.Getenv("THOST"), t[1], b)
		fmt.Printf("-f -n %s -i %s -p %d -tP -z30\n", t[1],
			m[0], portToUse + (lnode * 100) + node)
		conf = addTPBridge(conf, t[1], m[0], "P", portToUse + (lnode * 100) + node)
		
	}
		
	return conf

}

func main() {
	pc,err := net.ListenPacket("udp4", ":8829")
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	bcast := getDefaultBroadcastAddress()

	fmt.Printf("broadcast addrss is: %s\n", bcast)

	nxconf := ReadNdrxconfig()

	fmt.Println("%#v", nxconf)
	
	go write(pc, bcast)

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
				nxconf = updateConfig(addr, buf[:n], nxconf)
				nodes[s[1]] = 0
			}
		}
	}
}
