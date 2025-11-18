package main

import (
	"filesystem/p2p"
	"log"
)

func main() {
	tr := p2p.NewTCPTransport(":3000")

	// Short variable declaration inside if statement
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
