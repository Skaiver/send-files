package ClientServerCommunication

import (
	"fmt"
	"log"
	"net"
)

func IsServerValidRemote(ip string, port string) bool {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write([]byte("ARE_U_A_SERVER?"))
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	buffer := make([]byte, 1024)
	mLen, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	responseMsg := string(buffer[:mLen])

	if responseMsg == "YES!" {
		defer conn.Close()
		return true
	} else {

		fmt.Println("Received: ", string(buffer[:mLen]))
		defer conn.Close()
	}
	return false
}
