package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		PORT = "8989"
		LaunchServer()
	}

	if len(os.Args) >= 3 {
		LaunchError()
	}

	if len(os.Args) == 2 {
		PORT = os.Args[1]
		LaunchServer()
	}
}

func LaunchError() {
	fmt.Println("[USAGE]: ./TCPChat $port")
	os.Exit(0)
}
