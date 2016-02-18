package main

import(
	"runtime"
	"fmt"
	."controller"
	."network"
	."driver"
	)


const N_FLOOR int = 4


type externalOrder struct{	
	new_order 		bool 
	executed_order 	bool
	floor 			int
	direction 		int 
}


type internalOrder struct{
	floor int 
}

type costInfo  struct{
	cost 	int
	floor 	int 
	dir 	bool
}

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())

	nextFloorChan 		:= make(chan int)
	currentFloorChan 	:= make(chan int)
	timerChan 			:= make(chan int)
	dirChan				:= make(chan int)
	shareOrderChan	 	:= make(chan externalOrder)
	receivedOrderChan	:= make(chan externalOrder)
	shareCostChan		:= make(chan costInfo)
	receivedCostChan	:= make(chan costInfo)
	newExternalOrderChan:= make(chan externalOrder)
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

