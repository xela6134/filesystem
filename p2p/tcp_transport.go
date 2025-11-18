package p2p

import (
	"fmt"
	"net"
	"sync"
)

// Represents the remote node over a TCP established connection
type TCPPeer struct {
	// Underlying connection of the peer
	conn net.Conn

	// If we dial and retrieve a conn   -> outbound -> outbound == true
	// If we accept and retrieve a conn -> inbound  -> outbound == false
	outbound bool
}

type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	shakeHands    HandshakeFunc
	decoder       Decoder

	mutex sync.RWMutex
	peers map[net.Addr]Peer
}

// Public functions, start with uppercase not lowercase

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		shakeHands:    NOPHandshakeFunc,
		listenAddress: listenAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

// Private functions, start with lowercase not uppercase

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %v\n", err)
		}

		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	// %+v includes field names
	fmt.Printf("new incoming connection %+v\n", peer)

	if err := t.shakeHands(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %s\n", err)
		return
	}

	// Read loop
	msg := &Temp{}

	for {
		if err := t.decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP error: %s\n", err)
			continue
		}
	}

}
