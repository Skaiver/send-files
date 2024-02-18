package main

import (
	"local.com/MyClient"
)

func main() {
	systemIp := MyClient.GetOutboundIP()

	// identify hosts in network
	// server := listAvailableServers(systemIp)

	MyClient.CallBroadcast(systemIp)

	// connectToServer(server)
}
