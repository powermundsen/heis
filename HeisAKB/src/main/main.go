package main

import (
	"controller"
	"datatypes"
	"manager"
	"network"
	"runtime"
	"ioHandling"
	"time"
	"os"
)

const n_FLOORS int = 4

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	nextFloorChan := make(chan int)
	var doorCloseChan <-chan time.Time
	dirChan := make(chan datatypes.Direction)
	shareOrderChan := make(chan datatypes.ExternalOrder)
	receivedOrderChan := make(chan datatypes.ExternalOrder)
	shareCostChan := make(chan datatypes.CostInfo)
	receivedCostChan := make(chan datatypes.CostInfo)
	orderFinishedChan := make(chan int, 10)

	newInternalOrderChan := make(chan datatypes.InternalOrder)
	newExternalOrderChan := make(chan datatypes.ExternalOrder, 2)
	currentFloorToOrderManagerChan := make(chan int)
	currentFloorToElevControllerChan := make(chan int)
	setInternalLightsChan := make(chan datatypes.InternalOrder, 10)
	setExternalLightsChan := make(chan datatypes.ExternalOrder, 10)
	setDoorOpenLightChan := make(chan bool)
	setMotorDirectionChan := make(chan datatypes.Direction)

	if !network.InitNetworkHandler(shareOrderChan, receivedOrderChan, shareCostChan, receivedCostChan){
		os.Exit(1)
	}
	manager.InitOrderManager(n_FLOORS, newInternalOrderChan, newExternalOrderChan, currentFloorToOrderManagerChan, orderFinishedChan, dirChan, receivedOrderChan, receivedCostChan, shareOrderChan, shareCostChan, nextFloorChan, setInternalLightsChan, setExternalLightsChan)
	ioHandling.InitIo(n_FLOORS, newInternalOrderChan, newExternalOrderChan, currentFloorToOrderManagerChan, currentFloorToElevControllerChan, setInternalLightsChan, setExternalLightsChan, setDoorOpenLightChan, setMotorDirectionChan)
	controller.InitElevController(n_FLOORS, nextFloorChan, currentFloorToElevControllerChan, doorCloseChan, dirChan, setDoorOpenLightChan, setMotorDirectionChan, orderFinishedChan) 

	select {}
}
