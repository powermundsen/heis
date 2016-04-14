package manager

import (
	//"time"
	"fmt"
	//."driver"
	"datatypes"
)

func InitManager(newInternalOrderChan chan datatypes.InternalOrder, nextFloorChan chan int) {
	go RelayOrderToController(newInternalOrderChan, nextFloorChan)
}

func RelayOrderToController(newInternalOrderChan chan datatypes.InternalOrder, nextFloorChan chan int) {
	for {
		select {
		case newInternalOrder := <-newInternalOrderChan:
			nextFloorChan <- newInternalOrder
			fmt.Println("Sent order to controller")
		}
	}

}
