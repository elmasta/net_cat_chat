package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

// Gère la première connexion du client.
// Initialise dans la structure AllClients
// -> renvoie le pingouin, récupère le pseudo, met à jour les données
// /!\ mutex ?

func FirstConnection(conn net.Conn) {
	conn.Write(EntryMessage())

	buffer := make([]byte, 64)

	n, err := conn.Read(buffer)
	if err != nil {
		return
	}

	if n == 1 {
		conn.Write([]byte("Please choose a name.\n"))
		conn.Close()
		return
	}

	if n == 64 {
		conn.Write([]byte("Please choose a shorter name.\n"))
		conn.Close()
		return
	}

	name := buffer[:n-1]

	allClients.mu.Lock()
	nameUsed := false
	for _, x := range allClients.clients {
		if string(x.nickname) == string(name) {
			nameUsed = true
			break
		}
	}
	if nameUsed {
		conn.Write([]byte("This name is already used. You have been disconnected"))
		conn.Close()
		allClients.mu.Unlock()
		return
	}

	if allClients.total >= MAXCLIENT {
		conn.Write([]byte("Too many users. You have been disconnected."))
		conn.Close()
		allClients.mu.Unlock()
		return
	}

	client := Client{nickname: name, conn: conn}
	allClients.clients = append(allClients.clients, client)
	allClients.total++
	allClients.mu.Unlock()

	SendToAll([]byte(string(name)+" has connected.\n"), client, false)

	log := GetLog()
	if log != nil {
		conn.Write(log)
	}

	go HandleMessage(client)
}

func FakeHandleMessage(client Client) {
	for {
		buffer := make([]byte, 1024)
		_, err := client.conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
		}

	}
}

func EntryMessage() []byte {
	f, _ := os.ReadFile("entryMessage.txt")

	return f
}

func GetLog() []byte {
	fileDate := time.Now()
	now := fileDate.Format("2006-01-02")
	f, err := os.ReadFile("./log - " + now)
	if err != nil {
		return nil
	}

	return f
}
