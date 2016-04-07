package driver
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "timer.h"
*/
import "C"
//import "fmt"

func Timer_start(duration double){
	C.timer_start(C.double(duration))
}

func Timer_stop(){
	C.timer_stop(void)
}

func Timer_timedOut(){
	C.timer_timedOut(void)
}