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
	currentFloorToElevControllerChan chan int, setInternalLightsChan chan []int, setExternalLightsChan chan []int) {
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
	setInternalLightsChan chan []int, setExternalLightsChan chan []int) {
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
	var order datatypes.InternalOrder
	for floor := 1; floor < n_FLOORS+1; floor++ {
		if driver.Elevator_is_button_pushed(driver.BUTTON_INSIDE_COMMAND, floor) {
			order.Floor = floor
			newInternalOrderChan <- order
		}
	}
}

func readAllExternalButtons(newExternalOrderChan chan datatypes.ExternalOrder) {
	var order datatypes.ExternalOrder
	order.New_order = true
	order.Executed_order = false

	for button := driver.BUTTON_OUTSIDE_UP; button <= driver.BUTTON_OUTSIDE_DOWN; button++ {
		for floor := 1; floor < n_FLOORS+1; floor++ {
			if driver.Elevator_is_button_pushed(button, floor) { //dobbelsjekk denne if'en
				order.Floor = floor
				order.Direction = int(button) //endre her
				newExternalOrderChan <- order
			}
		}
	}
}

func readCurrentFloor(currentFloorToOrderManagerChan chan int, currentFloorToElevControllerChan chan int) {
	currentFloor := driver.Elevator_get_floor_sensor_signal()
	if currentFloor != -1 {
		driver.Elevator_set_floor_indicator(currentFloor)
	}
	currentFloorToOrderManagerChan <- currentFloor
	currentFloorToElevControllerChan <- currentFloor
}

func setExternalOrderLights(setExternalLightsChan chan []int) {
	lightSlice := <-setExternalLightsChan
	for i := 0; i < n_FLOORS+1; i++ {
		driver.Elevator_set_button_lamp(driver.BUTTON_OUTSIDE_UP, i, lightSlice[n_FLOORS])
		driver.Elevator_set_button_lamp(driver.BUTTON_OUTSIDE_DOWN, i, lightSlice[n_FLOORS+i])
	}
}

func setInternalOrderLights(setInternalLightsChan chan []int) {
	lightSlice := <-setInternalLightsChan
	for i := 0; i < n_FLOORS+1; i++ {
		driver.Elevator_set_button_lamp(driver.BUTTON_INSIDE_COMMAND, i, lightSlice[i])
	}

}
