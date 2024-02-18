package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	ip := GetOutboundIP()
	fmt.Println(ip)
	fmt.Println("Server Running...")
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("Error listening:", err.Error())

	}
	defer ln.Close()

	fmt.Println("Waiting for client...")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error listening:", err.Error())
		}
		go processClient(conn)
	}
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func processClient(conn net.Conn) {
	fmt.Println("client connected")

	buffer := make([]byte, 1024)
	mLen, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("Received: ", string(buffer[:mLen]))
	_, err = conn.Write([]byte("Thanks! Got your message:" + string(buffer[:mLen])))
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	conn.Close()

}
