package ioHandling

import (
	"datatypes"
	"driver"
	"fmt"
	"time"
)

var n_floors int
var previous_floor int

func InitIo(number_of_floors int, newInternalOrderChan chan datatypes.InternalOrder,
	newExternalOrderChan chan datatypes.ExternalOrder, currentFloorToOrderManagerChan chan int,
	currentFloorToElevControllerChan chan int, setInternalLightsChan chan datatypes.InternalOrder, 
	setExternalLightsChan chan datatypes.ExternalOrder, setDoorOpenLightChan chan bool, 
	setMotorDirectionChan chan datatypes.Direction) {
	if driver.Elevator_init() == 0 {
		fmt.Println("Could not connect to IO")
		return
	}
	fmt.Println("IO init done")
	n_floors = number_of_floors
	previous_floor = -1
	go ioManager(newInternalOrderChan, newExternalOrderChan, currentFloorToOrderManagerChan,
		currentFloorToElevControllerChan, setInternalLightsChan, setExternalLightsChan,
		setDoorOpenLightChan, setMotorDirectionChan)
}

func ioManager(newInternalOrderChan chan datatypes.InternalOrder, newExternalOrderChan chan datatypes.ExternalOrder,
	currentFloorToOrderManagerChan chan int, currentFloorToElevControllerChan chan int,
	setInternalLightsChan chan datatypes.InternalOrder, setExternalLightsChan chan datatypes.ExternalOrder,
	setDoorOpenLightChan chan bool, setMotorDirectionChan chan datatypes.Direction) {
	for {
		select {
		case <-time.After(100 * time.Millisecond):
			readAllInternalButtons(newInternalOrderChan)
			readAllExternalButtons(newExternalOrderChan)
			readCurrentFloor(currentFloorToElevControllerChan, currentFloorToOrderManagerChan)

		case internal_orders := <-setInternalLightsChan:
			setInternalOrderLights(internal_orders)

		case shared_orders := <-setExternalLightsChan:
			setExternalOrderLights(shared_orders)

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

func readCurrentFloor(currentFloorToElevControllerChan chan int, currentFloorToOrderManagerChan chan int) {
	current_floor := driver.Elevator_get_floor_sensor_signal()
	if current_floor != previous_floor {
		if current_floor != -1 {
			driver.Elevator_set_floor_indicator(current_floor)
		}
		currentFloorToElevControllerChan <- current_floor
		previous_floor = current_floor
	}

}

func setDoorOpenLight(set_door_open_light bool) {
	driver.Elevator_set_door_open_lamp(set_door_open_light)
}

func setMotorDirection(motor_direction datatypes.Direction) {
	if motor_direction == datatypes.UP {
		driver.Elevator_set_motor_direction(driver.MOTOR_DIRECTION_UP)
	} else if motor_direction == datatypes.DOWN {
		driver.Elevator_set_motor_direction(driver.MOTOR_DIRECTION_DOWN)
	} else {
		driver.Elevator_set_motor_direction(driver.MOTOR_DIRECTION_STOP)
	}
}
/*
func setInternalOrderLights(set_internal_lights []bool) {
	for i := 0; i < n_floors+1; i++ {
		driver.Elevator_set_button_lamp(driver.BUTTON_INSIDE_COMMAND, i, set_internal_lights[i])
	}
}


func setExternalOrderLights(set_external_lights []bool) {
	for i := 0; i < n_floors+1; i++ {
		driver.Elevator_set_button_lamp(driver.BUTTON_OUTSIDE_UP, i, set_external_lights[i])
		driver.Elevator_set_button_lamp(driver.BUTTON_OUTSIDE_DOWN, i, set_external_lights[n_floors+i])
	}
}
*/
//LEGGES INN IGJEN ETTER TESTING MED BOOL-VERSJON

func setExternalOrderLights(external_orders datatypes.ExternalOrder) { //Skal ta inn sharedOrder Slice med external_order IKKE int
	if len(external_orders) == 0 {
		fmt.Println("No new orders that need lights changed in recieved slice")
	} else {
		if items.Direction == 1 {
			driver.Elevator_set_button_lamp(driver.BUTTON_OUTSIDE_UP, items.Floor, true)
		}
		if items.Direction == -1 {
			driver.Elevator_set_button_lamp(driver.BUTTON_OUTSIDE_DOWN, items.Floor, true)
		}

		
	}
}


func setInternalOrderLights(internal_orders datatypes.ExternalOrder) {
	if len(internal_orders) == 0 {
		fmt.Println("No new orders that need lights changed in recieved slice")
	} else {
		driver.Elevator_set_button_lamp(driver.BUTTON_INSIDE_COMMAND, items.Floor, true)
		
	}

}

