package main

import (
	"fmt"
	"log"
	"net"
	"os"
	// "strconv"
	// "bufio"
)

func main() {
	connectToServer()
}

func getTarget() string {
	var in *os.File
	var err error

	// flag solution
	// arg0 := flag.String("file", "", "enter a file path")
	// flag.Parse()

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

func connectToServer() {
	remote := "192.168.178.70:8080"

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
