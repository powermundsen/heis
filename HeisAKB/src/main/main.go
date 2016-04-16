package main

import (
	"controller"
	"datatypes"
	//"fmt"
	"manager"
	"network"
	"runtime"
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
	dirChan := make(chan datatypes.Direction)
	shareOrderChan := make(chan datatypes.ExternalOrder)
	receivedOrderChan := make(chan datatypes.ExternalOrder)
	shareCostChan := make(chan datatypes.CostInfo)
	receivedCostChan := make(chan datatypes.CostInfo)
	orderFinishedChan := make(chan int, 10)

	// InitIO
	newInternalOrderChan := make(chan datatypes.InternalOrder)
	newExternalOrderChan := make(chan datatypes.ExternalOrder, 10)
	currentFloorToOrderManagerChan := make(chan int)
	currentFloorToElevControllerChan := make(chan int)
	setInternalLightsChan := make(chan datatypes.InternalOrder, 10)
	setExternalLightsChan := make(chan datatypes.ExternalOrder, 10)
	setDoorOpenLightChan := make(chan bool)
	setMotorDirectionChan := make(chan datatypes.Direction)

	network.InitNetworkHandler(shareOrderChan, receivedOrderChan, shareCostChan, receivedCostChan)
	manager.InitOrderManager(n_FLOORS, newInternalOrderChan, newExternalOrderChan, currentFloorToOrderManagerChan, orderFinishedChan, dirChan, receivedOrderChan, receivedCostChan, shareOrderChan, shareCostChan, nextFloorChan, setInternalLightsChan, setExternalLightsChan)
	ioHandling.InitIo(n_FLOORS, newInternalOrderChan, newExternalOrderChan, currentFloorToOrderManagerChan, currentFloorToElevControllerChan, setInternalLightsChan, setExternalLightsChan, setDoorOpenLightChan, setMotorDirectionChan)
	//time.Sleep(50 * time.Millisecond) 		//mulig denne trengs!!
	controller.InitElevController(n_FLOORS, nextFloorChan, currentFloorToElevControllerChan, doorCloseChan, dirChan, setDoorOpenLightChan, setMotorDirectionChan, orderFinishedChan) //initChan lagt til for fun

	//go network.ListenForExternalOrders(receivedOrderChan) //go network.ListenForExternalOrders(externalOrderChan)
	//go network.ListenForInternalOrders(receivedOrderChan) //finner ikke network.ListenForInternalOrder(internalOrderChan)

	select {}
}
