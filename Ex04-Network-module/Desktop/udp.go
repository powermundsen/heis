/*package main;

import(
	"fmt"
	"net"
)

const(
SERVER_PORT = "33546"
SERVER_IP = "129.241.187.23"
CONN_TYPE = "udp"
)

func main(){
	response := make([]byte,1024);
	udpServerAddr,_ := net.ResolveUDPAddr(CONN_TYPE, SERVER_IP + ":" + SERVER_PORT);
	udpConnection, err := net.DialUDP("udp", nil, udpServerAddr);
	message := []byte("Connect to: 129.241.187.142:20014\x00")
	//a := "halla"
	n,err := udpConnection.Write(message)
	fmt.Println(message)


	fmt.Println("Antall pen som er blitt sendt er:",n)
	fmt.Println("Feilmeldingen:",err)
	udpConnection.Read(response)
	fmt.Println(string(response))
}*/

package main;

import (
	"fmt"
	"net"
	"strconv"
	"time"
);

func listen(conn *net.UDPConn) {
	
	buffer := make([]byte, 1024);

	for {
		messageSize, _, _ := conn.ReadFromUDP(buffer);
		fmt.Println("Received: " + string(buffer[0:messageSize]));
	}
}

func transmit(conn *net.UDPConn) {
	
	for {
		time.Sleep(2000*time.Millisecond);

		message1 := "Hello server";
		message2 := "you are an asshole";
		conn.Write([]byte(message1));
		fmt.Println("Sent: " + message1);
		conn.Write([]byte(message2));

	}
}

func main() {

	serverIP := "129.241.187.255";
	serverPort := 20014;
	serverAddr, _ := net.ResolveUDPAddr("udp", serverIP + ":" + strconv.Itoa(serverPort));

	listenPort := 20014;
	listenAddr, _ := net.ResolveUDPAddr("udp", ":" + strconv.Itoa(listenPort));

	fmt.Println(listenAddr);
	fmt.Println(serverAddr);

	listenConn, _ := net.ListenUDP("udp", listenAddr);
	transmitConn, _ := net.DialUDP("udp", nil, serverAddr);

	go listen(listenConn);
	go transmit(transmitConn);

	d_chan := make(chan bool, 1);
	<- d_chan;
}