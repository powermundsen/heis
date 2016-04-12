package manager

import(
	"time"
	"fmt"
	"datatypes"
)

var next_floor 	int
var direction 	datatypes.Direction
var number_of_floors	int

shared_orders := make([]datatypes.ExternalOrder, 0)
private_orders := make([]datatypes.ExternalOrder, 0)

func InitOrderManager(n_FLOORS int, newInternalOrderChan chan datatypes.InternalOrder, newExternalOrderChan chan datatypes.ExternalOrder, currentFloorToOrderManagerChan chan int, orderFinnishedChan chan int, dirChan chan datatypes.Direction, receivedOrderChan chan datatypes.ExternalOrder, receivedCostChan chan datatypes.CostInfo,shareOrderChan chan datatypes.ExternalOrder, shareCostChan chan datatypes.CostInfo, nextFloorChan chan int){

	if _, err := os.Stat("privateOrdersBackup.txt"); os.IsNotExist(err) {
		createOrderBackupFile()
	}
	restoreOrders(private_orders)

	number_of_floors = n_FLOORS
	next_floor = -1 //-1 betyr at det ikke er noen etasje heisen skal gå til

	go orderManager(newInternalOrderChan, newExternalOrderChan, currentFloorToOrderManagerChan, orderFinnishedChan, dirChan, receivedOrderChan, receivedCostChan, 
					shareOrderChan, shareCostChan, nextFloorChan)//mangler setExternalLightsChan
}

func orderManager(newInternalOrderChan chan datatypes.InternalOrder, newExternalOrderChan chan datatypes.ExternalOrder, currentFloorToOrderManagerChan int,	orderFinnishedChan chan int, dirChan chan datatypes.Direction, receivedOrderChan chan datatypes.ExternalOrder, receivedCostChan chan datatypes.CostInfo,				shareOrderChan chan datatypes.ExternalOrder, shareCostChan chan datatypes.CostInfor, nextFloorChan chan int){
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
	if !(received_order){
		//del på shareOrderChan slik at den blir broadcastet til andre heiser
	}
	updateSharedOrders(order)

	if !(order.Executed_order){
		cost := calculateCost(order)
		shareCostChan <- cost
		setExternalLightsChan <- //senddetsommåsendes
		add_order := orderOnAuction(cost)
		if(add_order){
			new_internal_order := order.Floor
			newInternalOrderchan <- new_internal_order
		}
	}

}

func handlePrivateOrder(order datatypes.InternalOrder, finnished bool){
	updatePrivateOrders(order)
	findNextFloorToGoTo()
}

func updatePrivateOrders(order datatypes.InternalOrder){
	private_orders_copy := make([]datatypes.ExternalOrder, 0)
	private_orders_copy = private_orders 

	if len(private_orders_copy) == 0{
		private_orders_copy= append(private_orders_copy, new_order)

	} else {
		for items := range private_orders_copy {
			if private_orders_copy[items] != new_order  {
				private_orders_copy = append(private_orders_copy, new_order)
			} 
		}
	}
	private_orders = private_orders_copy
	fmt.Println("Mottatt meldling er:")
	fmt.Println(new_order)

	backupOrders(private_orders)
}

func updateSharedOrders(order datatypes.ExternalOrder){
	shared_orders_copy := make([]datatypes.InternalOrder, 0)
	shared_orders_copy = private_orders 

	if len(shared_orders_copy) == 0{
		shared_orders_copy= append(shared_orders_copy, new_order)

	} else {
		for items := range shared_orders_copy {
			if shared_orders_copy[items] != new_order  {
				shared_orders_copy = append(shared_orders_copy, new_order)
			} 
		}
	}
	private_orders = shared_orders_copy
	fmt.Println("Mottatt meldling er:")
	fmt.Println(new_order)
}

func calculateCost(order datatypes.ExternalOrder) datatypes.CostInfo{
	//tar inn en bestilling, sjekker opp mot privateOrderMatrix(slice), ser på delta(etasjer), hvor mange stopp som er lagt inn på veien, samt current direction
	//	må bare bestemme oss for vektingen på etasje, stopp og endring i dir
	return cost
}

func orderOnAuction(my_cost datatypes.CostInfo, receivedCostChan chan datatypes.CostInfo) bool{
	//start en timer
	var external_cost datatypes.CostInfo 
	for{
		select{
			case <- time.After(time.Millisecond * 500): //Tiden for timeout må finnes ved prøving/feiling
        		fmt.Println("Auction timed out")
        		return true
			case external_cost <- receivedCostChan:
				if((external_cost.Floor == my_cost.Floor) && (external_cost.Direction == my_cost.Direction)){
					if(external_cost.Cost < my_cost.Cost){
						return false
					}else if((external_cost.Cost == my_cost.Cost) && (external_cost.IP < my_cost.IP)){
						return false
					}
				}
		}
	}
}

func findNextFloorToGoTo(){
	//samme logikk som i datastyring
	//if UP
		//noen etasjer over som har bestilling (opp eller intern)? Hvis ja, og den etasjen er lavere enn current next_floor, sett next_floor til den etasjen
		//hvordan er slice bygget opp? Det bestemmer rekkefølgen på søket
	//if DOWN

	//må se på hvordan slices fungerer, og bygget opp. Ut i fra det kan algoritmen (if-setningene) skrives
}

func restoreOrders(orders *[]datatypes.ExternalOrder){
    buffer, _ := ioutil.ReadFile("privateOrdersBackup.txt")
    err := json.Unmarshal(buffer, orders)
    if err != nil {
		fmt.Println("error:", err)
	}
}

func backupOrders(orders []datatypes.ExternalOrder){
    fmt.Println(orders)
    buffer, _ := json.Marshal(orders)
    err := ioutil.WriteFile("privateOrdersBackup.txt", []byte(buffer), 0777)
    if err != nil {
        log.Fatal(err)
    }
}

func createOrderBackupFile(){
	var (
    	newFile *os.File
   		err     error
	)

	newFile, err = os.Create("privateOrdersBackup.txt")
    if err != nil {
        log.Fatal(err)
    }

    log.Println(newFile)
    newFile.Close()
}
