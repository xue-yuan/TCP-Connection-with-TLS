package main

import (
	"fmt"
	"net"
)

var currentClientNum int8 = 0

// StartServer To start server
func StartServer() {
	listen, err := net.Listen("tcp", "127.0.0.1:5278")
	clientID := 0
	if err != nil {
		// error
	}

	defer closeServer(listen)
	for {
		connection, err := listen.Accept()
		if err != nil {
			// error
			continue
		}
		go handleConnection(connection, clientID)
		clientID++
		currentClientNum++
	}
}

func handleConnection(connection net.Conn, clientID int) {
	defer closeConnection(connection, clientID)
	fmt.Println("--- Accept Connection ---")
	showCurrentClientNum()

	readChannel := make(chan []byte, 1024)
	writeChannel := make(chan []byte, 1024)

	go readConnection(connection, readChannel)
	go writeConnection(connection, writeChannel)

	for {
		select {
		case msg := <-readChannel:
			if string(msg) == "exit\n" {
				return
			}
			writeChannel <- msg
		}
	}
}

func readConnection(connection net.Conn, readChan chan []byte) {
	buffer := make([]byte, 2048)
	for {
		n, err := connection.Read(buffer)
		if n > 0 {
			fmt.Println("> Receive:", string(buffer[:n]))
		}
		if err != nil {
			// error
		}
		readChan <- buffer[:n]
	}
}

func writeConnection(connection net.Conn, writeChan chan []byte) {
	for {
		select {
		case msg := <-writeChan:
			_, err := connection.Write(msg)
			fmt.Println("> Write:", string(msg))
			if err != nil {
				// error
			}
		}
	}
}

func closeConnection(connection net.Conn, clientID int) {
	fmt.Printf("--- Close ConnectionID: %d ---\n", clientID)
	currentClientNum--
	showCurrentClientNum()
	connection.Close()
}

func closeServer(listen net.Listener) {
	fmt.Println("--- Close Server ---")
	listen.Close()
}

func showCurrentClientNum() {
	fmt.Printf("--- Current Online: %d ---\n", currentClientNum)
}
