package driver
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "fsm.h"
*/
import "C"

//import "fmt"


func FSM_onInitBetweenFloors(){
	C.fsm_onInitBetweenFloors(void)
}

func FSM_onRequestButtonPress(btn_floor int, btn_type Button){
	C.fsm_onRequestButtonPress(C.int(btn_floor), C.elev_button_type(btn_type))
}


func FSM_onFloorArrival(newFloor int){
	C.fsm_onFloorArrival(C.int(newFloor))
}

func FSM_onDoorTimeout(){
	C.fsm_onDoorTimeout(void)
}
