package main

import (
	"fmt"

	"local.com/MyClient"
)

func main() {
	systemIp := MyClient.GetOutboundIP()

	fmt.Println(systemIp)

	// identify hosts in network
	// server := listAvailableServers(systemIp)

	// create buffered channel to push all remote server to
	var serversChannel = make(chan string)

	go MyClient.FindRemoteServers(systemIp, serversChannel)

	// push all online remote server into channel, then read it out here and display on cli
	fmt.Println("Following server found: ")
	for server := range serversChannel {
		fmt.Println(server)
	}

	option := MyClient.ChooseServer()

	go MyClient.ConnectToServer(option)
}
