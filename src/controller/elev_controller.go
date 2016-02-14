package controller

import(
	"time"
	"fmt"
	."driver"
)

var TOTAL_FLOORS int  
var current_floor int
var next_floor int
var direction int 
var door_open = false
var init_status  = false
var reset_timer = false

func InitElevController(N_FLOOR int, initChan chan bool, nextFloorChan chan int, currentFloorChan chan int, timerChan chan int, dirChan chan int ){
	init_status = true
	TOTAL_FLOORS = N_FLOOR
	direction = -1					
	dirChan <- direction 			
	Elevator_set_motor_direction(-1)

	go elevController(nextFloorChan , currentFloorChan, timerChan, dirChan)
}

func elevController(nextFloorChan chan int, currentFloorChan chan int, timerChan chan int, dirChan chan int ){

	for{
		select{
			case door_timer_is_done := <- timerChan:
				if door_timer_is_done == 1{
					closeDoor()
				}

			case floor_reached := <- currentFloorChan:
				current_floor = floor_reached
				if (current_floor != -1){
					floorReached(timerChan)
					dirChan <- direction
				}

			case go_to_next_floor := <- nextFloorChan:
				next_floor = go_to_next_floor
				goToFloor(timerChan)
				dirChan <- direction
			
			default:
				time.Sleep(5*time.Millisecond)	
		}
	}
}

func goToFloor(timerChan chan int){
	if next_floor == -1{
		direction = 0
		Elevator_set_motor_direction(0) //Elevator_set_motor_direction(direction)
	
	}else if next_floor == current_floor{
		floorReached(timerChan)
		// må muligens ha funksjon eller channel som sier i fra til brain at en bestilling er fullført
	}else if !door_open{
		if next_floor > current_floor{
			direction = 1
			Elevator_set_motor_direction(1) //Elevator_set_motor_direction(direction)

		}else if next_floor < current_floor{
			direction = -1
			Elevator_set_motor_direction(-1) //Elevator_set_motor_direction(direction)
		}
	}
}

func floorReached(timerChan chan int){
	Elevator_set_floor_indicator(current_floor)
	if init_status{
		direction = 0
		Elevator_set_motor_direction(0) //Elevator_set_motor_direction(direction)
		init_status = false 
	}
	if current_floor == next_floor{
		fmt.Println("Stopping at floor #",current_floor)
		Elevator_set_motor_direction(0)
		door_open = true
		go startDoorTimer(timerChan)					
		Elevator_set_door_open_lamp(1)
		if current_floor == TOTAL_FLOORS{
			direction = -1				
		}else if current_floor == 1{
			direction = 1				
		}
	}
}

func closeDoor(){
	if door_open{
		door_open = false
		Elevator_set_door_open_lamp(-1)
	}
}


func startDoorTimer(timerChan chan int)  {
	select {
		case <- time.After(3 * time.Second):
			timerChan <- 1
	}
}