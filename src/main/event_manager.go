package main

import(
	"runtime"
	"fmt"
	"controller"
	"manager"
	//"network"
	"driver"
	"time"
	)


const N_FLOOR int = 4


type internalOrder struct{
	floor int 
}


func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())

	nextFloorChan 		:= make(chan int)
	currentFloorChan 	:= make(chan int)
	timerChan 			:= make(chan int)
	dirChan				:= make(chan int)
	//shareOrderChan	 	:= make(chan network.ExternalOrder)
	//receivedOrderChan	:= make(chan network.ExternalOrder)
	//shareCostChan		:= make(chan network.CostInfo)
	//receivedCostChan	:= make(chan network.CostInfo)
	//newExternalOrderChan:= make(chan network.ExternalOrder)
	newInternalOrderChan 	:= make(chan int)  //(chan internalOrder)
	initChan				:= make(chan bool) //nødvendig for controller.InitElevController


	
	if (driver.Elevator_init() == 0){
		fmt.Println("Could not connect to IO")
		return
	}
	//network.InitNetworkHandler(shareOrderChan, receivedOrderChan, shareCostChan, receivedCostChan)
	//manager.InitOrderManager(shareOrderChan, receivedOrderChan, shareCostChan, receivedCostChan, newExternalOrderChan, newInternalOrderChan, dirChan)
	controller.InitElevController(N_FLOOR, initChan, nextFloorChan, currentFloorChan, timerChan, dirChan) //initChan lagt til for fun
	//go network.ListenForExternalOrders(receivedOrderChan) //go network.ListenForExternalOrders(externalOrderChan)
	//go network.ListenForInternalOrders(receivedOrderChan) //finner ikke network.ListenForInternalOrder(internalOrderChan)

	/*for{
		select{
			case new_external_order := <- receivedOrderChan: //case new_external_order := <- externalOrderChan:
				shareOrderChan<-new_external_order
				
		}
	}
	*/
	//Registrer knapper
	fmt.Println("Init done")
	manager.InitManager(newInternalOrderChan, nextFloorChan)
	for {
		fmt.Println("Reached for-loop")
		if(driver.Elevator_is_button_pushed(2,1)){
			fmt.Println("Button 2 pushed")
			newInternalOrderChan <- 2
			//manager.RelayOrderToController(newInternalOrderChan, nextFloorChan)
		}
		if(driver.Elevator_is_button_pushed(2,2)){
			newInternalOrderChan <- 3
			//manager.RelayOrderToController(newInternalOrderChan, nextFloorChan)
		}
		time.Sleep(250*time.Millisecond)
	}
	


	/*Send knappetrykk til ordermanager på kanal

	Ordermanager sender disse videre til elevator controller
	så kjører elevator controller heisen til riktig etasje*/


}

