package main

import (
	"fmt"

	"local.com/MyClient"
)

func main() {
	systemIp := MyClient.GetOutboundIP()

	// identify hosts in network
	// server := listAvailableServers(systemIp)

	// create buffered channel to push all remote server to
	var serversChannel = make(chan string)

	go MyClient.FindRemoteServers(systemIp, serversChannel)

	// push all online remote server into channel, then read it out here and display on cli
	for _, server := range <-serversChannel {
		fmt.Println(server)
	}

	// option := MyClient.ChooseServer()

	// MyClient.ConnectToServer(option)
}
