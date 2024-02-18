package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	PORT = "8081"
)

func main() {
	// identify hosts in network
	server := listAvailableServers()

	connectToServer(server)
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func listAvailableServers() string {
	// send requests to all hosts with application port
	// assuming last byte octet is host part: 192.168.178.XXX
	ip := "192.168.178."

	fmt.Println("Available Servers: ")

	for i := 0; i < 255; i++ {
		fullIp := ip + strconv.Itoa(i)

		// check if server is acutally there before try to resolve
		timeout := 1 * time.Millisecond
		conn, err := net.DialTimeout("tcp", fullIp+":"+PORT, timeout)
		if err != nil {
			// log.Println("Site unreachable, error: ", err)
			continue
		}
		defer conn.Close()

		fmt.Println("IP: " + fullIp + ":" + PORT)
	}

	systemIp := GetOutboundIP()
	fmt.Println("Your IP: ", systemIp.String())
	fmt.Print("insert last octet of server: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	fmt.Println(input.Text())

	return input.Text()
}

func getTarget() string {
	var in *os.File
	var err error

	arg0 := os.Args[1:][0]

	switch name := arg0; {
	case name == "":
		in = os.Stdin
	default:
		if in, err = os.Open(name); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(arg0)

	stat, err := in.Stat()
	if err != nil {
		fmt.Println(err)
	}

	data := make([]byte, stat.Size())
	in.Read(data)

	fmt.Println(data)
	fmt.Println(string(data))

	return string(data)
}

func connectToServer(server string) {
	remote := "192.168.178." + server + ":" + PORT
	fmt.Println("aiafgnig", remote)
	//establish connection
	connection, err := net.Dial("tcp", remote)
	if err != nil {
		panic(err)
	}

	content := getTarget()

	//send some data
	_, err = connection.Write([]byte(content))

	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("Received: ", string(buffer[:mLen]))
	defer connection.Close()
}
