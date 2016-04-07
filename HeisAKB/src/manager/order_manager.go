package manager

import(
	//"time"
	"fmt"
	//."driver"
)


func InitManager(newInternalOrderChan chan int, nextFloorChan chan int){
	go RelayOrderToController(newInternalOrderChan, nextFloorChan)
}


func RelayOrderToController(newInternalOrderChan chan int, nextFloorChan chan int){
	for{
		select{
			case newInternalOrder := <- newInternalOrderChan:
				nextFloorChan <- newInternalOrder
				fmt.Println("Sent order to controller")
			default:

		}
	}

}


