package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"filesystem/crypto"
	"filesystem/p2p"
	"filesystem/store"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)

	fileServerOpts := FileServerOpts{
		EncKey:            crypto.NewEncryptionKey(),
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: store.CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main() {
	// Servers look like this) :3000 <----> :5000 <----> :7000
	s1 := makeServer(":3000", "")
	s2 := makeServer(":7000", "")
	s3 := makeServer(":5000", ":3000", ":7000")

	// Start servers and let it sleep to bind to port
	go func() { log.Fatal(s1.Start()) }()
	time.Sleep(2 * time.Second)

	go func() { log.Fatal(s2.Start()) }()
	time.Sleep(2 * time.Second)

	go func() { log.Fatal(s3.Start()) }()
	time.Sleep(2 * time.Second)

	// Here, s3 sender
	// s1, s2 receiver
	for i := range 20 {
		// Make filename & file content
		key := fmt.Sprintf("picture_%d.png", i)
		data := bytes.NewReader([]byte("my big data file here!"))

		// Store file in the distributed system through hashing & P2P
		s3.Store(key, data)

		// Delete local copy on s3
		if err := s3.store.Delete(s3.ID, key); err != nil {
			log.Fatal(err)
		}

		// Retrieve the file from peers since the local copy was deleted
		r, err := s3.Get(key)
		if err != nil {
			log.Fatal(err)
		}

		// Read the file into memory
		b, err := io.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(b))
	}
}
