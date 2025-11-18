package p2p

// Interface to represent a remote node
type Peer interface {
}

// Handles communication between nodes in the network
// Can be TCP, UDP, Websockets etc
type Transport interface {
	ListenAndAccept() error
}
