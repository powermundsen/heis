package main

import(
	"runtime"
	"fmt"
	"controller"
	"network/network_handler"
	"driver"
	)


const N_FLOOR int = 4

/*
type externalOrder struct{	
	new_order 		bool 
	executed_order 	bool
	floor 			int
	direction 		int 
}
*/

type internalOrder struct{
	floor int 
}
/*
type costInfo  struct{
	cost 	int
	floor 	int 
	dir 	bool
}
*/

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())

	nextFloorChan 		:= make(chan int)
	currentFloorChan 	:= make(chan int)
	timerChan 			:= make(chan int)
	dirChan				:= make(chan int)
	shareOrderChan	 	:= make(chan network.ExternalOrder)
	receivedOrderChan	:= make(chan network.ExternalOrder)
	shareCostChan		:= make(chan network.CostInfo)
	receivedCostChan	:= make(chan network.CostInfo)
	newExternalOrderChan:= make(chan network.ExternalOrder)
	newInternalOrderChan 	:= make(chan internalOrder)

	
	if (Elevator_init() == 0){
		fmt.Println("Could not connect to IO")
		return
	}
	InitNetworkHandler(shareOrderChan, receivedOrderChan, shareCostChan, receivedCostChan)
	InitOrderManager(shareOrderChan, receivedOrderChan, shareCostChan, receivedCostChan, newExternalOrderChan, newInternalOrderChan, dirChan)
	InitElevController(N_FLOOR, nextFloorChan, currentFloorChan, timerChan, dirChan)
	go listenForExternalOrders(externalOrderChan)
	go listenForInternalOrders(internalOrderChan)

	for{
		select{
			case new_external_order := <- externalOrderChan:
				shareOrderChan<-new_external_order
				
		}
	}

}

