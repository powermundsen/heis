package ioHandling

import (
	"datatypes"
	"driver"
	"fmt"
	"time"
)

var n_floors int

func InitIo(number_of_floors int, newInternalOrderChan chan datatypes.InternalOrder,
	newExternalOrderChan chan datatypes.ExternalOrder, currentFloorToOrderManagerChan chan int,
	currentFloorToElevControllerChan chan int, setInternalLightsChan chan []bool, setExternalLightsChan chan []bool,
	setDoorOpenLightChan chan bool, setMotorDirectionChan chan datatypes.Direction) {
	if driver.Elevator_init() == 0 {
		fmt.Println("Could not connect to IO")
		return
	}
	n_floors = number_of_floors
	go ioManager(newInternalOrderChan, newExternalOrderChan, currentFloorToOrderManagerChan,
		currentFloorToElevControllerChan, setInternalLightsChan, setExternalLightsChan,
		setDoorOpenLightChan, setMotorDirectionChan)
}

func ioManager(newInternalOrderChan chan datatypes.InternalOrder, newExternalOrderChan chan datatypes.ExternalOrder,
	currentFloorToOrderManagerChan chan int, currentFloorToElevControllerChan chan int,
	setInternalLightsChan chan []bool, setExternalLightsChan chan []bool,
	setDoorOpenLightChan chan bool, setMotorDirectionChan chan datatypes.Direction) {
	for {
		fmt.Println("InitElevControllerio")

		select {
		case <-time.After(100 * time.Millisecond):
			readAllInternalButtons(newInternalOrderChan)
			readAllExternalButtons(newExternalOrderChan)
			readCurrentFloor(currentFloorToElevControllerChan, currentFloorToOrderManagerChan)
		case set_internal_lights := <-setInternalLightsChan:
			setInternalOrderLights(set_internal_lights)
		case set_external_lights := <-setExternalLightsChan:
			setExternalOrderLights(set_external_lights)
		case set_door_open_light := <-setDoorOpenLightChan:
			setDoorOpenLight(set_door_open_light)
		case motor_direction := <-setMotorDirectionChan:
			fmt.Println("case motordirection")
			setMotorDirection(motor_direction)
		}
	}
}

func readAllInternalButtons(newInternalOrderChan chan datatypes.InternalOrder) {
	var order datatypes.InternalOrder
	for floor := 0; floor < n_floors; floor++ {
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
		for floor := 0; floor < n_floors; floor++ {
			if driver.Elevator_is_button_pushed(button, floor) {
				order.Floor = floor
				if button == driver.BUTTON_OUTSIDE_UP {
					order.Direction = int(datatypes.UP)
				} else {
					order.Direction = int(datatypes.DOWN)
				}
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
	//currentFloorToOrderManagerChan <- currentFloor
	currentFloorToElevControllerChan <- currentFloor
}


func setDoorOpenLight(set_door_open_light bool) {
	driver.Elevator_set_door_open_lamp(set_door_open_light)
}

func setMotorDirection(motor_direction datatypes.Direction) {
	if motor_direction == datatypes.UP {
		driver.Elevator_set_motor_direction(1)
	} else if motor_direction == datatypes.DOWN {
		driver.Elevator_set_motor_direction(-1)
	} else {
		driver.Elevator_set_motor_direction(0)
	}

}

func setExternalOrderLights(setExternalLightsChan chan []ExternalOrder) { //Skal ta inn sharedOrder Slice med external_order IKKE int
	lightSlice := <-setExternalLightsChan
	if len(lightSlice) == 0{
		fmt.Println("No new orders that need lights changed in recieved slice")
	} else {
		for items := range lightSlice {
				select {
					case items.Direction == 1:
						driver.Elevator_set_button_lamp(driver.BUTTON_OUTSIDE_UP, items.Floor , true)
					case items.Direction == -1:
						driver.Elevator_set_button_lamp(driver.BUTTON_OUTSIDE_DOWN, items.Floor , true)
					}
		}
	}
}

func setInternalOrderLights(setInternalLightsChan chan []ExternalOrder) {
	lightSlice := <-setInternalLightsChan
	if len(lightSlice) == 0{
		fmt.Println("No new orders that need lights changed in recieved slice")
	} else {
		for items := range lightSlice {
				driver.Elevator_set_button_lamp(driver.BUTTON_INSIDE_COMMAND, items.Floor , true)
			}
		}

}

