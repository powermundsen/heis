
package main
 
import (
    "fmt"
    "net"
    "time"
    "strconv"
    "os/exec"
    //"os"
    
)

//Lag en som lytter
//Hvis ingenting på 7 sekunder (Master dør)
	//Bli master
	//Lag klon deg selv
	//Start klone
	//Tell videre fra siste mottatt

var counter = 0
var server_ip = "localhost:20011"
var server_port = 20911
var local_ip = "localhost:30000"
var local_port = 30090



 /////////////////CLIENT//////////////////////////

func main() {

	fmt.Println("Starting main")
    primary()
    go secondary()
    startNew()

    neverExit := make(chan int)
    <-neverExit
}

 /////////////////SERVER//////////////////////////

 /* A Simple function to verify error */

/*
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}
*/

func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
    }
}


func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
    _,err := conn.WriteToUDP([]byte("abcFrom server: Hello I got your message "), addr)
    if err != nil {
        fmt.Printf("Couldn't send response %v", err)
    }
}
 

func primary(){
			///////////LYTTER////////////
	ServerAddr,err := net.ResolveUDPAddr("udp", server_ip)
	CheckError(err)


	Conn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)

	defer Conn.Close() 
	Conn.SetDeadline(time.Now().Add(time.Second*4))
	//i := 0
	buf  := make([]byte, 1024)
	for {
		fmt.Println("testa")
		n, _ , err := Conn.ReadFromUDP(buf)
		fmt.Println("test2")
		
		if err != nil {
			fmt.Println("Error: " , err)
			break
		}
	    
		nerr, _ := err.(net.Error);
		 if nerr.Timeout() {
		 	fmt.Println("timeout")
		 }
		
	    fmt.Println("test3")

	    buffer := string(buf[0:n])
	    fmt.Println(buffer)
	    fmt.Println("test4")
	    
		counter, err = strconv.Atoi(buffer)
	   	Conn.SetDeadline(time.Now().Add(time.Second*4))

	    fmt.Println("Received ", counter)

	}

	//Hvis connection timed out
		//Drep funksjonen os.Exit? 


}

func secondary() {
	
	//////////BROADCASTER///////////////////////

	    /* Lets prepare a address at any address at port 10001*/   
	ServerAddr,err := net.ResolveUDPAddr("udp", server_ip)
	CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", local_ip)
    CheckError(err)

	/* Now listen at selected port */
	ServerConn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	//buf := make([]byte, 1024)

	for {
		counter++
		fmt.Println(counter)
		counter_sending := strconv.Itoa(counter)
		ServerConn.WriteToUDP([]byte(counter_sending) , LocalAddr)

	    time.Sleep(time.Second * 1)
	    
	}


}

func startNew() {
	fmt.Println("startNew test")
	//os.StartProcess("gnome-terminal -x [\".\\oving6\"]", )
	cmd := exec.Command("gnome-terminal", "-e", "./oving6")
	cmd.Start()

	select {}
}