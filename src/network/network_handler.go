package network

import(
	_ "fmt"
	"time"
	"net"
	"encoding/json"
	"log"
	"strconv"

)



type ExternalOrder struct{	
	new_order 		bool 
	executed_order 	bool
	floor 			int
	direction 		int 
}

type CostInfo  struct{
	cost 		int
	floor 		int 
	direction 	int
}



var broadcast_order_conn	*net.UDPConn
var receive_order_conn 		*net.UDPConn
var broadcast_cost_conn		*net.UDPConn
var receive_cost_conn		*net.UDPConn


const BROADCAST_IP		= "255.255.255.255"
const PORTNUM_ORDER 	= ":7100"	
const PORTNUM_COST		= ":8100"	
const NUMBER_OF_BROADCAST = 5

func InitNetworkHandler(shareOrderChan chan ExternalOrder,receivedOrderChan chan ExternalOrder,shareCostChan chan CostInfo,receivedCostChan chan CostInfo){
	if (initSockets(BROADCAST_IP,PORTNUM_ORDER,PORTNUM_COST) == false){ 
		return
	}else{
		go networkHandler(shareOrderChan,receivedOrderChan,shareCostChan,receivedCostChan)
	}
}

func networkHandler(shareOrderChan chan ExternalOrder,receivedOrderChan chan ExternalOrder,shareCostChan chan CostInfo,receivedCostChan chan CostInfo){
	

	go listenForExternalOrder(receivedOrderChan)
	go listenForCostUpdate(receivedCostChan)

	for {
		select{
			case order:= <- shareOrderChan:	
				broadcastExternalOrder(order)

			case cost := <- shareCostChan:
				broadcastCostUpdate(cost)										
			default:
				time.Sleep(5*time.Millisecond)
		}
	}
	defer broadcast_order_conn.Close()
	defer broadcast_cost_conn.Close()
}


func initSockets(BROADCAST_IP string, PORTNUM_COST string, PORTNUM_ORDER string) bool {

	//Create broadcast socket for orders
	broadcast_udp_addr,err:=net.ResolveUDPAddr("udp",BROADCAST_IP + PORTNUM_ORDER)
	if err != nil{
		log.Println(" ResolveUDPAddr failed",err) 
	}

	broadcast_order_conn,err:=net.DialUDP("udp",nil,broadcast_udp_addr)
	if err != nil{
		log.Println("Could not establish UDP connection. Enter single eleveator mode. \n",err)
		return false 
	}
	broadcast_order_conn = broadcast_order_conn


	//Create broadcast socket for cost updates
	broadcast_udp_addr,err = net.ResolveUDPAddr("udp",BROADCAST_IP + PORTNUM_COST)
	if err != nil{
		log.Println("ResolveUDPAddr failed",err)
	}

	broadcast_cost_conn,err =net.DialUDP("udp",nil,broadcast_udp_addr)
	if err != nil{
		log.Println("Could not establish UDP connection. Enter single eleveator mode. \n",err)
		return false 
	}
	broadcast_cost_conn = broadcast_cost_conn


	//Create receiver socket for external orders
	listen_addr,err :=net.ResolveUDPAddr( "udp",PORTNUM_ORDER)
	if err != nil{
		log.Println("ResolveUDPAddr failed ",err)
	}

	receive_order_conn,err := net.ListenUDP("udp", listen_addr)
	if err != nil{
		log.Println("Could not establish UDP connection. Enter single eleveator mode. \n",err)
		return false 
	}
	receive_order_conn = receive_order_conn


	//Create receiver socket for cost updates
	listen_addr,err =net.ResolveUDPAddr( "udp",PORTNUM_COST)
	if err != nil{
		log.Println("ResolveUDPAddr failed ",err)
	}

	receive_cost_conn,err := net.ListenUDP("udp", listen_addr)
	if err != nil{
		log.Println("Could not establish UDP connection. Enter single eleveator mode.\n",err)
		return false
	}
	receive_cost_conn = receive_cost_conn

	return true 
}

func listenForExternalOrder(receivedOrderChan chan ExternalOrder){
	buffer := make([] byte, 1024)
	var external_order ExternalOrder

	for {
		len,received_ip,err := receive_order_conn.ReadFromUDP(buffer)
		if err != nil{
			log.Println("Not able to receive external order from ",received_ip)
		}
		err = json.Unmarshal(buffer[0:len], &external_order)
		if(err != nil) {log.Println("Error with Unmarshal \t",err)}
		receivedOrderChan <- external_order 	
		buffer = clearBuffer(buffer,len)	
	}
	defer receive_order_conn.Close() 
}

func broadcastExternalOrder(external_order ExternalOrder){
	buffer, err := json.Marshal(external_order)
	if err != nil{log.Println("Error with Marshal in broadcastExternalOrder() \t",err)}
	for i:=0; i < NUMBER_OF_BROADCAST; i++{
		_,err := broadcast_order_conn.Write(buffer) // broadcast_conn.Write(buffer)
		if err != nil {log.Println("Cannot send order over UDP \t",err)}
	}
}

func listenForCostUpdate(receivedCostChan chan CostInfo){
	buffer := make([] byte, 1024)
	var cost_info CostInfo

	for {
		len,received_ip,err := receive_cost_conn.ReadFromUDP(buffer)
		if err != nil{
			log.Println("Not able to receive cost update from ",received_ip)
		}
		err = json.Unmarshal(buffer[0:len], &cost_info)
		if(err != nil) {log.Println("Error with Unmarshal \t",err)}
		receivedCostChan <- cost_info 	
		buffer = clearBuffer(buffer,len)	
	}
	defer receive_cost_conn.Close()
}
func broadcastCostUpdate(cost_info CostInfo){
	buffer, err := json.Marshal(cost_info)
	if err != nil{log.Println("Error with Marshal in broadcastCostUpdate() \t",err)}
	for i:=0; i < NUMBER_OF_BROADCAST; i++{
		_,err := broadcast_cost_conn.Write(buffer) // broadcast_conn.Write(buffer)
		if err != nil {log.Println("Cannot send cost update over UDP \t",err)}
	}
}

func clearBuffer(buffer []byte, len int) []byte {
	var clear uint8
	clear = 0
	for i:=0;i < len; i++{
		buffer[i] = clear
	}
	return buffer	
}

func getLocalIp() net.IP {	
	local_listen_port 	:= 7600
	addr, _ := net.ResolveUDPAddr("udp4", "255.255.255.255:"+ strconv.Itoa(local_listen_port))

	temp_conn, _ := net.DialUDP("udp4", nil, addr)
	defer temp_conn.Close()
	temp_addr := temp_conn.LocalAddr()
	local_addr, _ := net.ResolveUDPAddr("udp4", temp_addr.String())
	
	return local_addr.IP
}