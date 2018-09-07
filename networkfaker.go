package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/williamfhe/godivert"
)

var currentListenerConditions []string

// FuckSourceSocketOfPacket fucks all network packets which are filtered by the listenerCondition with the given parameters.
// Not needed parameters can be left blank (if so they will not change the corresponding value)
func FuckSourceSocketOfPacket(listenerCondition string, newSourceIP string, newSourcePort int, newTargetIP string, newTargetPort int) {
	if !contains(currentListenerConditions, listenerCondition) {
		currentListenerConditions = append(currentListenerConditions, listenerCondition)
	} else {
		return
	}
	fmt.Println("start fucking:", currentListenerConditions)
	winDivert, err := godivert.NewWinDivertHandle(listenerCondition)
	if err != nil {
		panic(err)
	}
	defer winDivert.Close()

	packetChan, err := winDivert.Packets()
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go fuckPacket(winDivert, packetChan, newSourceIP, newSourcePort, newTargetIP, newTargetPort)
	wg.Wait()
}

func fuckPacket(wd *godivert.WinDivertHandle, packetChan <-chan *godivert.Packet, newSourceIP string, newSourcePort int, newTargetIP string, newTargetPort int) {
	for packet := range packetChan {
		if newSourceIP != "" {
			packet.SetSrcIP(net.ParseIP(newSourceIP))
		}
		if newSourcePort > 0 {
			packet.SetSrcPort(uint16(newSourcePort))
		}
		if newTargetIP != "" {
			packet.SetDstIP(net.ParseIP(newTargetIP))
		}
		if newTargetPort > 0 {
			packet.SetDstPort(uint16(newTargetPort))
		}
		packet.Send(wd)
	}
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
