package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func LaunchServer() {
	listen, err := net.Listen("tcp4", ":"+PORT)
	if err != nil {
		log.Fatal(err)
		LaunchError()
	}
	// close listener
	defer listen.Close()

	allClients.total = 0
	MAXCLIENT = 3

	fmt.Println("Listening on port " + PORT)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go FirstConnection(conn)
	}
}
