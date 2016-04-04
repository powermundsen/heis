package io

import (
	"datatypes"
	"driver"
	"fmt"
	"time"
)

func InitIo(newInternalOrderChan chan structs.InternalOrder, newExternalOrderChan chan structs.ExternalOrder, currentFloorToOrderManagerChan chan int, currentFloorToElevControllerChan chan int, setInternalLightsChan chan int, setExternalLightsChan chan int) {
	if driver.Elevator_init() == 0 {
		fmt.Println("Could not connect to IO")
		return
	}
	go ioManager(newExternalOrderChan, newInternalOrderChan, currentFloorToOrderManagerChan, currentFloorToElevControllerChan, setInternalLightsChan, setExternalLightsChan)
}

func ioManager(newInternalOrderChan chan structs.InternalOrder, newExternalOrderChan chan structs.ExternalOrder, currentFloorToOrderManagerChan chan int, currentFloorToElevControllerChan chan int, setInternalLightsChan chan int, setExternalLightsChan chan int) {
	for {
		readAllInternalButtons(newInternalOrderChan)
		readAllExternalButtons(newExternalOrderChan)
		readCurrentFloor(currentFloorToElevControllerChan, currentFloorToOrderManagerChan)
		setInternalOrderLights(setInternalLightsChan)
		setExternalOrderLights(setExternalLightsChan)
		time.Sleep(25 * time.Millisecond)
	}
}

func readAllInternalButtons() {
	var order InternalOrder
	for floor := 1; floor < 5; floor++ { //numberOfFloors
		if driver.Elevator_is_button_pushed(2, floor) {
			order = floor
			newInternalOrderChan <- order
	}
}

func readAllExternalButtons() {
	var order ExternalOrder
	order.new_order = true
	order.executed_order = false

	for floor := 1; floor < 5; floor++ {
		for button := 0; button < 2; button++ {
			if driver.Elevator_is_button_pushed(button, floor) {
				order.floor = floor
				order.direction = button
				newExternalOrderChan <- order
		}
	}
	}
}

func readCurrentFloor() {
	currentFloor := driver.Elevator_get_floor_sensor_signal()
	driver.Elevator_set_floor_indicator(currentFloor)
	currentFloorToOrderManagerChan <- currentFloor
	currentFloorToElevControllerChan <- currentFloor
}

func setExternalOrderLights() {

}

func setInternalOrderLights() {

}
