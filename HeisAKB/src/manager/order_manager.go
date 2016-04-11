package manager

import(
	//"time"
	"fmt"
	"datatypes"
)


func InitOrderManager(n_FLOORS int, newInternalOrderChan chan datatypes.InternalOrder, newExternalOrderChan chan datatypes.ExternalOrder, currentFloorToOrderManagerChan chan int,
				orderFinnishedChan chan int, dirChan chan int, receivedOrderChan chan datatypes.ExternalOrder, receivedCostChan chan datatypes.CostInfo,
				shareOrderChan chan datatypes.ExternalOrder, shareCostChan chan datatypes.CostInfo, nextFloorChan chan int){
	
	//createSharedOrdersMatrix, slice?
	//createPrivateOrdersMatrix, slice?
	//getOldOrdersFromFile, Andre, det finnes en funksjon vi allerede kan kalle?
	//sette n_FLOORS

	go orderManager(newInternalOrderChan, newExternalOrderChan, currentFloorToOrderManagerChan, orderFinnishedChan, dirChan, receivedOrderChan, receivedCostChan, 
					shareOrderChan, shareCostChan, nextFloorChan)
}

func orderManager(newInternalOrderChan chan datatypes.InternalOrder, newExternalOrderChan chan datatypes.ExternalOrder, currentFloorToOrderManagerChan int,
				orderFinnishedChan chan int, dirChan chan int, receivedOrderChan chan datatypes.ExternalOrder, receivedCostChan chan datatypes.CostInfo,
				shareOrderChan chan datatypes.ExternalOrder, shareCostChan chan datatypes.CostInfor, nextFloorChan chan int){
				//mangler setExternalLightsChan()

	for{
		select{
			case new_internal_order := <- newInternalOrderChan:
				finnished := false
				handlePrivateOrder(new_internal_order)

			case new_external_order := <- newExternalOrderChan:
				received_order := false
				handleSharedOrder(new_external_order, received_order := false)

			case received_network_order := <- receivedOrderChan:
				received_order := true
				handleSharedOrder(received_network_order, received_order := true)

			case order_finnished := <- orderFinnishedChan:
				finnished := true
				handlePrivateOrder()
		}
	}

}

func handleSharedOrder(order datatypes.ExternalOrder, received_order bool, shareCostChan chan datatypes.CostInfo, receivedCostChan chan datatypes.CostInfo){ 
						//mangler setExternalLightsChan og newInternalOrderChan
	updateSharedOrders(order)
	if !received_order{
		//del på shareOrderChan slik at den blir broadcastet til andre heiser
	}
	cost := calculateCost(order)
	shareCostChan <- cost
	setExternalLightsChan <- senddetsommåsendes
	add_order := orderOnAuction(cost)
	if(add_order){
		new_internal_order := order.Floor
		newInternalOrderchan <- new_internal_order
	}

}

func handlePrivateOrder(order datatypes.InternalOrder, finnished bool){
	updatePrivateOrders(order) //trenger vi finnished(bool)? Bør kanskje InternalOrder inneholde en slik bool, slik som ExternalOrder har?
	findNextFloorToGoTo()
}

func updatePrivateOrders(order datatypes.InternalOrder){
	//updatePrivateOrdersMatrix()
	//writePrivateOrdersToFile()
}

func updateSharedOrders(order datatypes.ExternalOrder){
	//updateSharedOrdersMatrix()
	//writeSharedOrdersToFile()

}

func calculateCost(order datatypes.ExternalOrder) datatypes.CostInfo{
	//tar inn en bestilling, sjekker opp mot privateOrderMatrix(slice), ser på delta(etasjer), hvor mange stopp som er lagt inn på veien, samt current direction
	//	må bare bestemme oss for vektingen på etasje, stopp og endring i dir
	return cost
}

func orderOnAuction(my_cost datatypes.CostInfo, receiveCostChan chan datatypes.CostInfo) bool{
	//start en timer
	var external_cost datatypes.CostInfo 
	for(ikke timer gått ut){
		external_cost <- receiveCostChan
		if((external_cost.Floor == my_cost.Floor) && (external_cost.Direction == my_cost.Direction)){
			if(external_cost.Cost < my_cost.Cost){
				return false
			}
			else if((external_cost.Cost == my_cost.Cost) && (external_cost.IP < my_cost.IP)){
				return false
			}
		}
	}
	return true
}

func findNextFloorToGoTo(){

}
