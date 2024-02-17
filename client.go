package main

import (
	"fmt"
	"net"
)

func main() {
	connectToServer()
}

func connectToServer() {
	remote := "192.168.178.70:8080"

	//establish connection
	connection, err := net.Dial("tcp", remote)
	if err != nil {
		panic(err)
	}

	///send some data
	_, err = connection.Write([]byte("Hello Server! Greetings."))
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("Received: ", string(buffer[:mLen]))
	defer connection.Close()
}
