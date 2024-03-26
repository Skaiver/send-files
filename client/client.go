package MyClient

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"local.com/ClientServerCommunication"
)

const (
	PORT = "8081"
)

func ChooseServer() string {
	fmt.Print("choose server to connect to: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	return input.Text()
}

func FindRemoteServers(systemIp net.IP, c chan<- string) {
	fmt.Println("should be going over hosts:")
	ip := systemIp

	// iterate over IP addresses in the subnet
	for i := 0; i < 255; i++ {
		if i == 0 {
			// get first actual device ip in subnet e.g. 192.168.17.1/36 -> only last octet is host part
			ip = getStartingIP(ip)
		} else {
			// increment ip to e.g. 192.168.17.2/36
			ip = incrementIP(ip)
		}
		// fmt.Println("testing now: ", ip)

		// send ping
		connectAttempt := raw_connect(ip.String(), PORT)

		if connectAttempt {
			if ClientServerCommunication.IsServerValidRemote(ip.String(), PORT) {
				//fmt.Println("Host found:", ip.String())

				// send available remote server ip to buffered channel
				c <- ip.String()
				// fmt.Println(ip.String())
			}
		}
	}
}

func testIfServerIsAvailableHost(ip string) bool {
	conn, err := net.Dial("tcp", ip+":"+PORT)
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

func getCIDRAdress() string {
	interfaces, err := net.Interfaces()
	// fmt.Println("interfaces: ", interfaces)
	if err != nil {
		fmt.Println("Fehler beim Abrufen der Netzwerkschnittstellen:", err)
	}

	// Durch die Netzwerkschnittstellen iterieren und die erste nicht-Loopback-Schnittstelle finden
	var cidrAddr string
	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback == 0 && iface.Flags&net.FlagUp != 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				fmt.Println("Fehler beim Abrufen der Adressen der Schnittstelle", iface.Name, ":", err)
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					_, bits := ipnet.Mask.Size()
					cidrAddr = fmt.Sprintf("%s/%d", ipnet.IP.String(), bits)
					break
				}
			}
		}
		if cidrAddr != "" {
			break
		}
	}

	// CIDR-Adresse ausgeben
	if cidrAddr == "" {
		fmt.Println("Keine aktive Netzwerkschnittstelle gefunden")
	} else {
		fmt.Println("CIDR-Adresse des Systems:", cidrAddr)
	}

	return cidrAddr
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

func listAvailableServers(systemIp net.IP) string {
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

func ConnectToServer(server string) {
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
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("Received: ", string(buffer[:mLen]))
	defer connection.Close()
}

// ParseCIDR parses a CIDR string and returns the *net.IPNet.
func parseCIDR(cidr string) (*net.IPNet, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	return ipNet, err
}

// IncrementIP increments an IP address.
func incrementIP(ip net.IP) net.IP {
	ipTo4 := ip.To4()
	ipTo4[3]++

	return net.ParseIP(ipTo4.String())
}

func getStartingIP(ip net.IP) net.IP {
	ipTo4 := ip.To4()
	ipTo4[3] = 0

	return net.ParseIP(ipTo4.String())
}

func raw_connect(host string, port string) bool {
	timeout := time.Millisecond * 10
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		// fmt.Println("Connecting error:", err)
	}
	if conn != nil {
		defer conn.Close()
		fmt.Println("Opened", net.JoinHostPort(host, port))
		return true
	}
	return false
}
