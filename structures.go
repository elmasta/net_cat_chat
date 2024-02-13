package main

import (
	"net"
	"sync"
)

var MessagingMutex sync.Mutex

type AllClients struct {
	mu      sync.Mutex
	total   int
	clients []Client
}

type Client struct {
	conn     net.Conn
	nickname []byte
}

var PORT string
var allClients AllClients
var MAXCLIENT int
