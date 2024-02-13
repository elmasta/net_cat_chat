package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func File2lines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LinesFromReader(f)
}

func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func SendToAll(message []byte, client Client, b bool) {
	//list client remplace client conn
	formTime := time.Now().Format("2006-01-02 15:04:05")
	MessagingMutex.Lock()
	//détecter si buffer est vide sinon mettre en pause le print
	for _, v := range allClients.clients {
		if string(v.nickname) != string(client.nickname) {
			if b {
				v.conn.Write([]byte("[" + formTime + "][" + string(client.nickname) + "]: " + string(message) + "\n"))
			} else {
				v.conn.Write([]byte(string(message)))
			}
		} else {
			v.conn.Write([]byte("\033[1A\033[K[" + formTime + "][" + string(client.nickname) + "]: " + string(message) + "\n"))
		}
	}

	//insert string to file
	fileDate := time.Now()
	now := fileDate.Format("2006-01-02")
	lines, err := File2lines("./log - " + now)
	if err != nil {
		fmt.Println(err)
	}
	fileContent := ""
	for _, line := range lines {
		fileContent += line
		fileContent += "\n"
	}
	if b {
		fileContent += "[" + formTime + "][" + string(client.nickname) + "]: " + string(message) + "\n"
	} else {
		fileContent += string(message)
	}
	ioutil.WriteFile("./log - "+now, []byte(fileContent), 0644)

	MessagingMutex.Unlock()
}

// Gère la réception et la diffusion des messages du client
func HandleMessage(client Client) {

	for {
		message := ReadMessage(client)
		if message == nil {
			HandleDisconnect(client)
			return
		}
		if len(message) > 0 {
			SendToAll(message, client, true)
		}
		// client.conn.Write([]byte("[" + string(client.nickname) + "]: "))
	}
}

func ReadMessage(client Client) []byte {
	bufferSize := 64
	buffer := make([]byte, bufferSize)
	res := ""
	n := bufferSize
	var err error

	for n == bufferSize && buffer[bufferSize-1] != '\n' {
		n, err = client.conn.Read(buffer)
		if err != nil {
			return nil
		}
		res += strings.Split(string(buffer), "\n")[0]
	}

	return []byte(res)
}

// Gère la déconnexion du client
func HandleDisconnect(client Client) {
	SendToAll([]byte(string(client.nickname)+" has left our chat...\n"), client, false)
	client.conn.Close()
	client.conn = nil
	allClients.mu.Lock()
	allClients.total--
	temp := []Client{}
	for _, v := range allClients.clients {
		if string(v.nickname) != string(client.nickname) {
			temp = append(temp, v)
		}
	}
	allClients.clients = temp
	allClients.mu.Unlock()
}
