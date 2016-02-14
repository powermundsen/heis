package main
 
import (
    "fmt"
    "net"
    "time"
    "strconv"
)
 
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
    }
}
 
func main() {
    ServerAddr,err := net.ResolveUDPAddr("udp","127.0.0.1:30000")
    CheckError(err)
 
    LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
    CheckError(err)
 
    Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
    CheckError(err)
 
    defer Conn.Close()
    i := 0
    for {
            b := make([]byte, UDP_PACKET_SIZE)
            n, addr, err := conn.ReadFromUDP(b)
            if err != nil {
                log.Printf("Error: UDP read error: %v", err)
                continue
            }
            Println("Recvd: %s", b)
        }
        time.Sleep(time.Second * 1)

}