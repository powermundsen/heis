package main

import (
	"controller"
	"datatypes"
	//"fmt"
	//"manager"
	"runtime"
	//"network"
	//"driver"
	"ioHandling"
	"time"
)

const n_FLOORS int = 4

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	nextFloorChan := make(chan int)
	//currentFloorChan := make(chan int)
	var doorCloseChan <-chan time.Time
	dirChan := make(chan int)
	//shareOrderChan	 	:= make(chan network.ExternalOrder)
	//receivedOrderChan	:= make(chan network.ExternalOrder)
	//shareCostChan		:= make(chan network.CostInfo)
	//receivedCostChan	:= make(chan network.CostInfo)
	//newExternalOrderChan:= make(chan network.ExternalOrder)
	//newInternalOrderChan := make(chan int) //(chan internalOrder)
	//initChan := make(chan bool) //nødvendig for controller.InitElevController

	// InitIO
	newInternalOrderChan := make(chan datatypes.InternalOrder)
	newExternalOrderChan := make(chan datatypes.ExternalOrder)
	currentFloorToOrderManagerChan := make(chan int)
	currentFloorToElevControllerChan := make(chan int)
	setInternalLightsChan := make(chan datatypes.InternalOrder)
	setExternalLightsChan := make(chan datatypes.ExternalOrder)
	setDoorOpenLightChan := make(chan bool)
	setMotorDirectionChan := make(chan datatypes.Direction)
	//network.InitNetworkHandler(shareOrderChan, receivedOrderChan, shareCostChan, receivedCostChan)
	//manager.InitOrderManager(shareOrderChan, receivedOrderChan, shareCostChan, receivedCostChan, newExternalOrderChan, newInternalOrderChan, dirChan)
	ioHandling.InitIo(n_FLOORS, newInternalOrderChan, newExternalOrderChan, currentFloorToOrderManagerChan, currentFloorToElevControllerChan, setInternalLightsChan, setExternalLightsChan, setDoorOpenLightChan, setMotorDirectionChan)
	//time.Sleep(50 * time.Millisecond) 		//mulig denne trengs!!
	controller.InitElevController(n_FLOORS, nextFloorChan, currentFloorToElevControllerChan, doorCloseChan, dirChan, setDoorOpenLightChan, setMotorDirectionChan) //initChan lagt til for fun
	//go network.ListenForExternalOrders(receivedOrderChan) //go network.ListenForExternalOrders(externalOrderChan)
	//go network.ListenForInternalOrders(receivedOrderChan) //finner ikke network.ListenForInternalOrder(internalOrderChan)

	//manager.InitManager(newInternalOrderChan, nextFloorChan)
	go relayOrderToController(newInternalOrderChan, nextFloorChan, setInternalLightsChan)

	select {}
}

func relayOrderToController(newInternalOrderChan chan datatypes.InternalOrder, nextFloorChan chan int, setInternalLightsChan chan datatypes.InternalOrder) {
	for {
		select {
		case newInternalOrder := <-newInternalOrderChan:
			setInternalLightsChan <- newInternalOrder

			nextFloorChan <- newInternalOrder.Floor

			newInternalOrder.Executed_order = true
			setInternalLightsChan <- newInternalOrder
			//fmt.Println("Sent order to controller")
		}
	}

}
