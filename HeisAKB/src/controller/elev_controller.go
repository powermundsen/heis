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
var reset_timer = false

func InitElevController(n_FLOORS int, nextFloorChan chan int, currentFloorChan chan int,
	timerChan chan int, dirChan chan int, setDoorOpenLightChan chan bool, setMotorDirectionChan chan datatypes.Direction) {
	fmt.Println("InitElevController1")
	init_status = true
	fmt.Println("InitElevController2")
	n_FLOORS = n_FLOORS
	fmt.Println("InitElevController3")
	direction = -1
	fmt.Println("InitElevController4")
	//dirChan <- direction
	fmt.Println("InitElevController5")
	//----------------------------------------------------------------
	//const dir datatypes.Direction
	setMotorDirectionChan <- datatypes.DOWN
	//mt.Println(setMotorDirectionChan)

	fmt.Println("InitElevController6")
	//----------------------------------------------------------------
	//Elevator_set_motor_direction(-1)

	go elevController(nextFloorChan, currentFloorChan, timerChan, dirChan, setDoorOpenLightChan, setMotorDirectionChan)
}

func elevController(nextFloorChan chan int, currentFloorChan chan int, timerChan chan int, dirChan chan int,
	setDoorOpenLightChan chan bool, setMotorDirectionChan chan datatypes.Direction) {

	for {
		select {
		case door_timer_is_done := <-timerChan:
			if door_timer_is_done == 1 {
				closeDoor(setDoorOpenLightChan)
			}

		case floor_reached := <-currentFloorChan:
			current_floor = floor_reached
			if current_floor != -1 {
				floorReached(timerChan, setMotorDirectionChan)
				//dirChan <- direction
			}

		case go_to_next_floor := <-nextFloorChan:
			fmt.Println("Reached elevController gotonextfloor-case")
			next_floor = go_to_next_floor
			fmt.Println("Next floor is: %d", next_floor)
			goToFloor(timerChan, setMotorDirectionChan)
			//dirChan <- direction

		default:
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func goToFloor(timerChan chan int, setMotorDirectionChan chan datatypes.Direction) {
	if next_floor == -1 {
		direction = 0
		setMotorDirectionChan <- datatypes.STOP
		//Elevator_set_motor_direction(0) //Elevator_set_motor_direction(direction)

	} else if next_floor == current_floor {
		floorReached(timerChan, setMotorDirectionChan)
		// må muligens ha funksjon eller channel som sier i fra til brain at en bestilling er fullført
	} else if !door_open {
		if next_floor > current_floor {
			direction = 1
			setMotorDirectionChan <- datatypes.UP
			//Elevator_set_motor_direction(1) //Elevator_set_motor_direction(direction)

		} else if next_floor < current_floor {
			direction = -1
			setMotorDirectionChan <- datatypes.DOWN
			//Elevator_set_motor_direction(-1) //Elevator_set_motor_direction(direction)
		}
	}
}

func floorReached(timerChan chan int, setMotorDirectionChan chan datatypes.Direction) {
	//Elevator_set_floor_indicator(current_floor)
	if init_status {
		direction = 0
		setMotorDirectionChan <- datatypes.STOP
		//Elevator_set_motor_direction(0) //Elevator_set_motor_direction(direction)
		init_status = false
	}
	if current_floor == next_floor {
		fmt.Println("Stopping at floor #", current_floor)
		setMotorDirectionChan <- datatypes.STOP
		//Elevator_set_motor_direction(0)
		door_open = true
		//go startDoorTimer(timerChan)
		setMotorDirectionChan <- datatypes.UP
		//Elevator_set_door_open_lamp(1)
		if current_floor == n_FLOORS {
			direction = -1
		} else if current_floor == 1 {
			direction = 1
		}
	}
}

func closeDoor(setDoorOpenLightChan chan bool) {
	if door_open {
		door_open = false
		setDoorOpenLightChan <- false
		//Elevator_set_door_open_lamp(false)
	}
}

func startDoorTimer(timerChan chan int) {
	select {
	case <-time.After(3 * time.Second):
		timerChan <- 1
	}
}
