package main

// Result is a combined type for channels to return result or error
type Result struct {
	err error
	res string
}

// ListenerResult contains the results of a listener service event call
type ListenerResult struct {
	Result
	resultBytes []byte // byte array of the result
	sourceIP    string // listener packet source ip
	sourcePort  int    // listener packet source port
}

// Data is a struct which contains all information which are sent
// between the instances of this tool
type Data struct {
	SrcIP    string `json:"srcIP"`
	DestIP   string `json:"destIP"`
	DestPort int    `json:"destPort"`
	Message  []byte `json:"message"`
}
