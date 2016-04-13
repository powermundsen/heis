package main

import (
	"controller"
	"datatypes"
	"fmt"
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
	setInternalLightsChan := make(chan []datatypes.InternalOrder)
	setExternalLightsChan := make(chan []datatypes.ExternalOrder)
	setDoorOpenLightChan := make(chan bool)
	setMotorDirectionChan := make(chan datatypes.Direction)
	//network.InitNetworkHandler(shareOrderChan, receivedOrderChan, shareCostChan, receivedCostChan)
	//manager.InitOrderManager(shareOrderChan, receivedOrderChan, shareCostChan, receivedCostChan, newExternalOrderChan, newInternalOrderChan, dirChan)
	ioHandling.InitIo(n_FLOORS, newInternalOrderChan, newExternalOrderChan, currentFloorToOrderManagerChan, currentFloorToElevControllerChan, setInternalLightsChan, setExternalLightsChan, setDoorOpenLightChan, setMotorDirectionChan)
	//feil currentfloorchan i controller
	time.Sleep(50 * time.Millisecond)
	controller.InitElevController(n_FLOORS, nextFloorChan, currentFloorToElevControllerChan, doorCloseChan, dirChan, setDoorOpenLightChan, setMotorDirectionChan) //initChan lagt til for fun
	//go network.ListenForExternalOrders(receivedOrderChan) //go network.ListenForExternalOrders(externalOrderChan)
	//go network.ListenForInternalOrders(receivedOrderChan) //finner ikke network.ListenForInternalOrder(internalOrderChan)

	/*for{
			select{
				case new_external_order := <- receivedOrderChan: //case new_external_order := <- externalOrderChan:
					shareOrderChan<-new_external_order

			}

	}	*/

	//Registrer knapper
	fmt.Println("All init done")
	//manager.InitManager(newInternalOrderChan, nextFloorChan)
	go relayOrderToController(newInternalOrderChan, nextFloorChan)

	/*
		for {
			time.Sleep(500 * time.Millisecond)
			fmt.Println("Running... ")
		}
	*/
	select {}
	/*
		for {
			if driver.Elevator_is_button_pushed(2, 1) {
				fmt.Println("Button 2 pushed")
				newInternalOrderChan <- 2
				//manager.RelayOrderToController(newInternalOrderChan, nextFloorChan)
			}
			if driver.Elevator_is_button_pushed(2, 2) {
				newInternalOrderChan <- 3
				//manager.RelayOrderToController(newInternalOrderChan, nextFloorChan)
			}
			time.Sleep(250 * time.Millisecond)
		}*/

	/*Send knappetrykk til ordermanager på kanal

	Ordermanager sender disse videre til elevator controller
	så kjører elevator controller heisen til riktig etasje*/

}

func relayOrderToController(newInternalOrderChan chan datatypes.InternalOrder, nextFloorChan chan int) {
	for {
		select {
		case newInternalOrder := <-newInternalOrderChan:
			nextFloorChan <- newInternalOrder.Floor
			//fmt.Println("Sent order to controller")
		}
	}

}
