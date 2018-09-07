package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
	"sync"
)

// main is the entry point of this tool
func main() {
	//host := flag.Bool("host", false, "true if you are hosts")
	listenerIP := flag.String("lip", "0.0.0.0", "listener ip")
	listenPort := flag.Int("lport", 5000, "the port on which will be listened for another FIFAFucker instance")
	targetIP := flag.String("tip", "127.0.0.1", "target ip of another FIFAFucker instance")
	targetPort := flag.Int("tport", 5000, "the target port of another FIFAFucker instance")

	flag.Parse()

	// start service listener (and decodes the message to re-send the packets which were forwarded)
	go startServiceListener(*listenerIP, *listenPort)
	// starts service listener on local FIFA instances
	go startFifaListener(*targetIP, *targetPort)

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

// startFifaListener listens on the UDP "socket" :9999 with broadcast as target IP (255.255.255.255)
// This method collects the broadcast packages of FIFA and fucks them while changing the destination IP to this local IP with listener port
// 255.255.255.255:9999 will be fucked to localIP:4987
func startFifaListener(targetIP string, targetPort int) {
	// fucks the BroadcastIP:9999 to localIP:4987
	listenPort := 4987
	localIP := GetLocalIP()
	fmt.Println("Fucked the FIFA Broadcast packets to local service", localIP, ":", listenPort, "...")
	go FuckSourceSocketOfPacket("udp.DstPort == 9999 and ip.DstAddr == 255.255.255.255", "", 0, localIP, 4987)
	fmt.Println("Starting listening for redirected FIFA packets (", listenPort, ":", listenPort, ")")
	resultChannel := make(chan ListenerResult)
	go StartListener(listenPort, resultChannel)
	// for every channel entry send a package to the target ip
	for elem := range resultChannel {
		jsonData, err := json.Marshal(Data{SrcIP: localIP, DestIP: "1.1.1.1", DestPort: 9999, Message: elem.resultBytes})
		CheckErrorInternal(err)
		msgJSONStr := string(jsonData)
		sendResult := make(chan Result)
		go SendMessage(targetIP, targetPort, msgJSONStr, sendResult)
		close(sendResult)
	}
}

// startServiceListener listens on the UDP "socket" :XXXXX for incoming administration UDP packets from another instance
// of this tool
func startServiceListener(listenerIP string, listenPort int) {
	fmt.Println("Starting listening for UDP packets ( :", listenPort, ")")
	resultChannel := make(chan ListenerResult)
	go StartListenerOnIP(listenerIP, listenPort, resultChannel)
	// for every channel entry sent the package to FIFA (and fucks it for a valid recognized packet)
	for elem := range resultChannel {
		bytes := elem.resultBytes
		var res Data
		err := json.Unmarshal(bytes, &res)
		CheckErrorInternal(err)
		newSrcIP := res.SrcIP
		newTargetIP := res.DestIP
		newTargetPort := res.DestPort
		sendResult := make(chan Result)
		localIP := GetLocalIP()
		fmt.Println("Fucked the service packets to FIFA", localIP, ":", newTargetPort, "...")
		go FuckSourceSocketOfPacket("udp.DstPort == "+strconv.Itoa(newTargetPort)+" and ip.DstAddr == "+newTargetIP, newSrcIP, 0, localIP, newTargetPort)
		go SendMessage(newTargetIP, newTargetPort, string(res.Message), sendResult)
		close(sendResult)
	}
}
