package main

import (
	"net"
	"strconv"
)

// StartListenerOnIP starts an UDP listener on the given listenIP and listenPort. Received messages or errors will be
// returned in the result channel.
func StartListenerOnIP(listenIP string, listenPort int, result chan ListenerResult) {
	startInternalListener(listenIP, listenPort, result)
}

// StartListener starts an UDP listener on 0.0.0.0 and listenPort. Received messages or errors will be
// returned in the result channel.
func StartListener(listenPort int, result chan ListenerResult) {
	startInternalListener("", listenPort, result)
}

// startInternalListener listens on a UDP "socket" (UDP has no sockets but yeah...) and reacts on errors or message receives.
func startInternalListener(listenIP string, listenPort int, result chan ListenerResult) {
	// Resolv listener address :port
	listenerStr := listenIP + ":" + strconv.Itoa(listenPort)
	ServerAddr, err := net.ResolveUDPAddr("udp", listenerStr)
	CheckErrorListener(err, result)

	// Start UDP listener
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckErrorListener(err, result)
	defer ServerConn.Close()

	buf := make([]byte, 1024)
loop:
	for {
		n, sourceAddress, err := ServerConn.ReadFromUDP(buf)
		bytes := buf[0:n]
		msg := string(bytes)
		if CheckListenerChannelIsClose(result) {
			break loop
		}
		CheckErrorListener(err, result)
		result <- ListenerResult{Result: Result{res: msg}, resultBytes: bytes, sourceIP: net.IP.String((*sourceAddress).IP), sourcePort: sourceAddress.Port}
	}
}

// SendMessage sends a UDP package to the given targetIP and targetPort "socket" and
// the request result will be queued to the result channel.
func SendMessage(targetIP string, targetPort int, msg string, result chan Result) {
	bindstr := targetIP + ":" + strconv.Itoa(targetPort)
	ServerAddr, err := net.ResolveUDPAddr("udp", bindstr)
	CheckError(err, result)

	LocalAddr, err := net.ResolveUDPAddr("udp", ":0")
	CheckError(err, result)

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err, result)

	defer Conn.Close()
	var buf = []byte(msg)
	_, err = Conn.Write(buf)
	CheckError(err, result)
}

// GetLocalIP gets the current LAN IP of this PC.
func GetLocalIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return net.IP.String((*localAddr).IP)
}
