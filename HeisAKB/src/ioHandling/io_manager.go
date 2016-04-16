package ioHandling

import (
	"datatypes"
	"driver"
	"fmt"
	"time"
)

//var n_floors int
var previous_floor int
var external_button_control_variable datatypes.ExternalOrder

func InitIo(number_of_floors int, newInternalOrderChan chan datatypes.InternalOrder,
	newExternalOrderChan chan datatypes.ExternalOrder, currentFloorToOrderManagerChan chan int,
	currentFloorToElevControllerChan chan int, setInternalLightsChan chan datatypes.InternalOrder, setExternalLightsChan chan datatypes.ExternalOrder,
	setDoorOpenLightChan chan bool, setMotorDirectionChan chan datatypes.Direction) {
	if driver.Elevator_init() == 0 {
		fmt.Println("Could not connect to IO")
		return
	}
	fmt.Println("IO init done")
	n_floors := number_of_floors
	previous_floor = -1
	external_button_control_variable = datatypes.ExternalOrder{New_order: true, Executed_order: false, Floor: -1, Direction: 0, Timestamp: time.Now().Unix()}
	fmt.Println("::::::::::::TImestamp er::::::::::::",external_button_control_variable.Timestamp)

	go ioManager(newInternalOrderChan, newExternalOrderChan, currentFloorToOrderManagerChan,
		currentFloorToElevControllerChan, setInternalLightsChan, setExternalLightsChan,
		setDoorOpenLightChan, setMotorDirectionChan, n_floors)
}

func ioManager(newInternalOrderChan chan datatypes.InternalOrder, newExternalOrderChan chan datatypes.ExternalOrder,
	currentFloorToOrderManagerChan chan int, currentFloorToElevControllerChan chan int,
	setInternalLightsChan chan datatypes.InternalOrder, setExternalLightsChan chan datatypes.ExternalOrder,
	setDoorOpenLightChan chan bool, setMotorDirectionChan chan datatypes.Direction, n_floors int) {
	

	for {
		select {
		case <-time.After(100 * time.Millisecond):
			detectAndSendInternalButtonCall(newInternalOrderChan, n_floors)
			detectAndSendExternalButtonCall(newExternalOrderChan, n_floors)
			readCurrentFloor(currentFloorToElevControllerChan, currentFloorToOrderManagerChan)

		case internal_order := <-setInternalLightsChan:
			fmt.Println("ioManager.Case: setInternalLights")
			setInternalOrderLights(internal_order)

		case external_order := <-setExternalLightsChan:
			fmt.Println("ioManager.Case: setExternalLights")
			setExternalOrderLights(external_order)

		case set_door_open_light := <-setDoorOpenLightChan:
			setDoorOpenLight(set_door_open_light)

		case motor_direction := <-setMotorDirectionChan:
			fmt.Println("ioManager.Case: setMotordirection")
			setMotorDirection(motor_direction)
		}
	}
}

func detectAndSendInternalButtonCall(newInternalOrderChan chan datatypes.InternalOrder, n_floors int) {
	var order datatypes.InternalOrder
	for floor := 0; floor < n_floors; floor++ {
		if driver.Elevator_is_button_pushed(driver.BUTTON_INSIDE_COMMAND, floor) {
			order.Floor = floor
			newInternalOrderChan <- order
		}
	}
}

func detectAndSendExternalButtonCall(newExternalOrderChan chan datatypes.ExternalOrder, n_floors int) {
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
				if ((order.Floor != external_button_control_variable.Floor) && (order.Direction != external_button_control_variable.Direction)){
					newExternalOrderChan <- order
					external_button_control_variable = order
					external_button_control_variable.Timestamp = time.Now().Unix()
					fmt.Println("::::::::::::ORDER TIMESTAMP::::::::::::",external_button_control_variable.Timestamp)

				} else if ((time.Now().Unix() - external_button_control_variable.Timestamp) > 2) {
					fmt.Println("::::::::::::REALTIME TIMESTAMP::::::::::::", time.Now().Unix())
					//fmt.Println(external_button_control_variable.Timestamp)
					//fmt.Println("------------------------------------------------------------")
					external_button_control_variable.Timestamp = time.Now().Unix()
					newExternalOrderChan <- order
				}
			}
		}
	}
}

func readCurrentFloor(currentFloorToElevControllerChan chan int, currentFloorToOrderManagerChan chan int){
	current_floor := driver.Elevator_get_floor_sensor_signal()
	if current_floor != previous_floor {
		if current_floor != -1 {
			driver.Elevator_set_floor_indicator(current_floor)
		}
		currentFloorToElevControllerChan <- current_floor
		currentFloorToOrderManagerChan <- current_floor
		previous_floor = current_floor
	}
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

func setInternalOrderLights(internal_order datatypes.InternalOrder) {
	if internal_order.Executed_order {
		driver.Elevator_set_button_lamp(driver.BUTTON_INSIDE_COMMAND, internal_order.Floor, false)
	} else {
		driver.Elevator_set_button_lamp(driver.BUTTON_INSIDE_COMMAND, internal_order.Floor, true)
	}
}

func setExternalOrderLights(external_order datatypes.ExternalOrder) {

	if external_order.Direction == 1 {
		if external_order.Executed_order {
			driver.Elevator_set_button_lamp(driver.BUTTON_OUTSIDE_UP, external_order.Floor, false)
		} else {
			driver.Elevator_set_button_lamp(driver.BUTTON_OUTSIDE_UP, external_order.Floor, true)
		}
	} else if external_order.Direction == -1 {
		if external_order.Executed_order {
			driver.Elevator_set_button_lamp(driver.BUTTON_OUTSIDE_DOWN, external_order.Floor, false)
		} else {
			driver.Elevator_set_button_lamp(driver.BUTTON_OUTSIDE_DOWN, external_order.Floor, true)
		}

	}
}
