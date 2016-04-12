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
	currentFloorChan := make(chan int)
	timerChan := make(chan int)
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
	setInternalLightsChan := make(chan []bool)
	setExternalLightsChan := make(chan []bool)
	setDoorOpenLightChan := make(chan bool)
	setMotorDirectionChan := make(chan datatypes.Direction)
	//network.InitNetworkHandler(shareOrderChan, receivedOrderChan, shareCostChan, receivedCostChan)
	//manager.InitOrderManager(shareOrderChan, receivedOrderChan, shareCostChan, receivedCostChan, newExternalOrderChan, newInternalOrderChan, dirChan)
	ioHandling.InitIo(n_FLOORS, newInternalOrderChan, newExternalOrderChan, currentFloorToOrderManagerChan, currentFloorToElevControllerChan, setInternalLightsChan, setExternalLightsChan, setDoorOpenLightChan, setMotorDirectionChan)
	controller.InitElevController(n_FLOORS, nextFloorChan, currentFloorChan, timerChan, dirChan, setDoorOpenLightChan, setMotorDirectionChan) //initChan lagt til for fun
	//go network.ListenForExternalOrders(receivedOrderChan) //go network.ListenForExternalOrders(externalOrderChan)
	//go network.ListenForInternalOrders(receivedOrderChan) //finner ikke network.ListenForInternalOrder(internalOrderChan)

	/*for{
			select{
				case new_external_order := <- receivedOrderChan: //case new_external_order := <- externalOrderChan:
					shareOrderChan<-new_external_order

			}

	}	*/

	//Registrer knapper
	fmt.Println("Init done")
	//manager.InitManager(newInternalOrderChan, nextFloorChan)
	go RelayOrderToController(newInternalOrderChan, nextFloorChan)

	for {
		time.Sleep(250 * time.Millisecond)
		fmt.Println("Running... ")
	}
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

func RelayOrderToController(newInternalOrderChan chan datatypes.InternalOrder, nextFloorChan chan int) {
	for {
		select {
		case newInternalOrder := <-newInternalOrderChan:
			nextFloorChan <- newInternalOrder.Floor
			fmt.Println("Sent order to controller")
		}
	}

}
