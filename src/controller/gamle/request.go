package driver
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "request.h"
*/
import( 
	"C"
	."controller"
)
//import "fmt"

func Requests_chooseDirection(Elevator e){
	C.requests_chooseDirection(C.Elevator(e)) __attribute__((pure))
}

func Requests_shouldStop(Elevator e){
	C.requests_shouldStop(C.Elevator(e)) __attribute__((pure))
}

func requests_clearAtCurrentFloor(Elevator e){
	C.requests_clearAtCurrentFloor(C.Elevator(e)) __attribute__((pure))
}