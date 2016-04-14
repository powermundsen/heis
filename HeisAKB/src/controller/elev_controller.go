package controller

import (
	"datatypes"
	//."driver"
	"fmt"
	"time"
)

var n_FLOORS int
var current_floor int
var next_floor int
var direction int
var door_open = false
var init_status = false

func InitElevController(number_of_floors int, nextFloorChan chan int,
	currentFloorToElevControllerChan chan int, doorCloseChan <-chan time.Time, dirChan chan int,
	setDoorOpenLightChan chan bool, setMotorDirectionChan chan datatypes.Direction) {

	init_status = true
	n_FLOORS = number_of_floors
	next_floor = -1
	direction = -1
	//dirChan <- direction
	setMotorDirectionChan <- datatypes.DOWN

	fmt.Println("elev_controller init done")
	go elevController(nextFloorChan, currentFloorToElevControllerChan, doorCloseChan, dirChan, setDoorOpenLightChan, setMotorDirectionChan)
}

func elevController(nextFloorChan chan int, currentFloorToElevControllerChan chan int, doorCloseChan <-chan time.Time, dirChan chan int,
	setDoorOpenLightChan chan bool, setMotorDirectionChan chan datatypes.Direction) {

	for {
		//fmt.Println("doorCloseChan:", doorCloseChan)
		select {
		case <-doorCloseChan:
			//fmt.Println("elevController.Case: doorCloseChan")
			door_open = false
			setDoorOpenLightChan <- false
			goToFloor(&doorCloseChan, setMotorDirectionChan, setDoorOpenLightChan)

		case floor_reached := <-currentFloorToElevControllerChan:
			fmt.Println("elevController.Case: floor_reached")
			current_floor = floor_reached
			if current_floor != -1 {
				floorReached(&doorCloseChan, setMotorDirectionChan, setDoorOpenLightChan)
				//dirChan <- direction
			}

		case go_to_next_floor := <-nextFloorChan:
			fmt.Println("elevController.Case: go_to_next_floor")
			next_floor = go_to_next_floor
			fmt.Println("Next floor is: ", next_floor)
			goToFloor(&doorCloseChan, setMotorDirectionChan, setDoorOpenLightChan)
			//dirChan <- direction
		}
	}
}

func goToFloor(doorCloseChan *<-chan time.Time, setMotorDirectionChan chan datatypes.Direction, setDoorOpenLightChan chan bool) {
	if next_floor == -1 {
		direction = 0
		setMotorDirectionChan <- datatypes.STOP

	} else if next_floor == current_floor {
		floorReached(doorCloseChan, setMotorDirectionChan, setDoorOpenLightChan)

	} else if !door_open {
		if current_floor == -1 {

		} else if next_floor > current_floor {
			direction = 1
			setMotorDirectionChan <- datatypes.UP

		} else if next_floor < current_floor {
			direction = -1
			setMotorDirectionChan <- datatypes.DOWN
		}
	}
}

func floorReached(doorCloseChan *<-chan time.Time, setMotorDirectionChan chan datatypes.Direction, setDoorOpenLightChan chan bool) {
	//Elevator_set_floor_indicator(current_floor)
	fmt.Println("Reached elevController.floorReached")
	if init_status {
		direction = 0
		setMotorDirectionChan <- datatypes.STOP
		init_status = false
	} else if current_floor == next_floor {
		fmt.Println("Stopping at floor #", current_floor)
		setMotorDirectionChan <- datatypes.STOP
		door_open = true
		*doorCloseChan = time.After(3 * time.Second)
		setDoorOpenLightChan <- true
		if current_floor == n_FLOORS {
			direction = -1
		} else if current_floor == 1 {
			direction = 1
		}
	}

}
