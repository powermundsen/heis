package ioHandling

import (
	"datatypes"
	"driver"
	"fmt"
	"time"
)

var n_FLOORS int

func InitIo(n_FLOORS int, newInternalOrderChan chan datatypes.InternalOrder, 
	newExternalOrderChan chan datatypes.ExternalOrder, currentFloorToOrderManagerChan chan int, 
	currentFloorToElevControllerChan chan int, setInternalLightsChan chan int, setExternalLightsChan chan int){
	if driver.Elevator_init() == 0 {
		fmt.Println("Could not connect to IO")
		return
	}
	n_FLOORS = n_FLOORS
	go ioManager(newInternalOrderChan, newExternalOrderChan, currentFloorToOrderManagerChan, 
		currentFloorToElevControllerChan, setInternalLightsChan, setExternalLightsChan)
}

func ioManager(newInternalOrderChan chan datatypes.InternalOrder, newExternalOrderChan chan datatypes.ExternalOrder,
	currentFloorToOrderManagerChan chan int, currentFloorToElevControllerChan chan int, 
	setInternalLightsChan chan int, setExternalLightsChan chan int){
	for {
		readAllInternalButtons(newInternalOrderChan)
		readAllExternalButtons(newExternalOrderChan)
		readCurrentFloor(currentFloorToElevControllerChan, currentFloorToOrderManagerChan)
		setInternalOrderLights(setInternalLightsChan)
		setExternalOrderLights(setExternalLightsChan)
		time.Sleep(25 * time.Millisecond)
	}
}

func readAllInternalButtons(newInternalOrderChan chan datatypes.InternalOrder) {
	var order InternalOrder
	for floor := 1; floor < n_FLOORS+1; floor++ { //numberOfFloors
		if driver.Elevator_is_button_pushed(2, floor) {
			order = floor
			newInternalOrderChan <- order
		}
	}
}

func readAllExternalButtons(newExternalOrderChan chan datatypes.ExternalOrder){
	var order ExternalOrder
	order.new_order = true
	order.executed_order = false

	for floor := 1; floor < n_FLOORS+1; floor++ {
		for button := 0; button < 2; button++ {
			if driver.Elevator_is_button_pushed(button, floor) {
				order.floor = floor
				order.direction = button
				newExternalOrderChan <- order
			}
		}
	}
}

func readCurrentFloor(currentFloorToOrderManagerChan chan int, currentFloorToElevControllerChan chan int) {
	currentFloor := driver.Elevator_get_floor_sensor_signal()
	driver.Elevator_set_floor_indicator(currentFloor)
	currentFloorToOrderManagerChan <- currentFloor
	currentFloorToElevControllerChan <- currentFloor
}

func setExternalOrderLights(setExternalLightsChan chan []int) {
	for i:= 0; i < n_FLOORS+1; i++ {
		driver.Elevator_set_button_lamp(driver.BUTTON_OUTSIDE_UP, i, setExternalLightsChan[1][i])
		driver.Elevator_set_button_lamp(driver.BUTTON_OUTSIDE_DOWN, i, setExternalLightsChan[2][i]) //[M_FLOORS+i]
	}
}

func setInternalOrderLights(setInternalLightsChan chan []int) {
	//size := len(setInternalLightsChan) //kan fjernes
	for i := 0; i < n_FLOORS+1; i++ {
		driver.Elecator_set_button_lamp(driver.BUTTON_INSIDE_COMMAND, i, setInternalLightsChan[i])
	}

}
