package manager

import (
	"datatypes"
	"fmt"
	"math"
	//"network"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var next_floor int
var direction datatypes.Direction
var number_of_floors int
var current_floor int

var shared_orders []datatypes.ExternalOrder
var private_orders []datatypes.ExternalOrder

func InitOrderManager(n_FLOORS int, newInternalOrderChan chan datatypes.InternalOrder, newExternalOrderChan chan datatypes.ExternalOrder, currentFloorToOrderManagerChan chan int, orderFinishedChan chan int, dirChan chan datatypes.Direction, receivedOrderChan chan datatypes.ExternalOrder, receivedCostChan chan datatypes.CostInfo, shareOrderChan chan datatypes.ExternalOrder, shareCostChan chan datatypes.CostInfo, nextFloorChan chan int) {
	shared_orders = make([]datatypes.ExternalOrder, 0)
	private_orders = make([]datatypes.ExternalOrder, 0)

	if _, err := os.Stat("privateOrdersBackup.txt"); os.IsNotExist(err) {
		createOrderBackupFile()
	}

	next_floor = -1
	number_of_floors = n_FLOORS
	direction = datatypes.DOWN

	restoreOrders(&private_orders)
	fmt.Println("Private orders from file are:", private_orders)
	findNextFloorToGoTo()
	fmt.Println("Init direction is: ", direction)

	go orderManager(newInternalOrderChan, newExternalOrderChan, currentFloorToOrderManagerChan, orderFinishedChan, dirChan, receivedOrderChan, receivedCostChan,
		shareOrderChan, shareCostChan, nextFloorChan) //mangler setExternalLightsChan
}

func orderManager(newInternalOrderChan chan datatypes.InternalOrder, newExternalOrderChan chan datatypes.ExternalOrder, currentFloorToOrderManagerChan chan int, orderFinishedChan chan int, dirChan chan datatypes.Direction, receivedOrderChan chan datatypes.ExternalOrder, receivedCostChan chan datatypes.CostInfo, shareOrderChan chan datatypes.ExternalOrder, shareCostChan chan datatypes.CostInfo, nextFloorChan chan int) {
	//mangler setExternalLightsChan()

	for {
		select {
		case new_internal_order := <-newInternalOrderChan:
			handleInternalOrder(new_internal_order)

		case new_external_order := <-newExternalOrderChan:
			received_order := false
			handleSharedOrder(new_external_order, received_order, shareCostChan, receivedCostChan)

		case received_network_order := <-receivedOrderChan:
			received_order := true
			handleSharedOrder(received_network_order, received_order, shareCostChan, receivedCostChan)

		case finished_order := <-orderFinishedChan:
			order := datatypes.ExternalOrder{New_order: false, Executed_order: true, Floor: finished_order, Direction: int(direction)}
			handleFinishedOrder(order, shareOrderChan)

		case update_current_floor := <-currentFloorToOrderManagerChan:
			current_floor = update_current_floor

		case update_direction := <-dirChan:
			if update_direction != datatypes.STOP {
				direction = update_direction
			}

		case <-time.After(2 * time.Second): //dette vil skape bug
			nextFloorChan <- next_floor
		}
	}
}

func handleSharedOrder(order datatypes.ExternalOrder, received_order bool, shareCostChan chan datatypes.CostInfo, receivedCostChan chan datatypes.CostInfo) {
	//mangler setExternalLightsChan og newInternalOrderChan
	if !(received_order) {
		//shareOrderChan <- order
	}
	updateSharedOrders(order)

	if !(order.Executed_order) {
		cost := calculateCost(order)
		//shareCostChan <- cost
		//setExternalLightsChan <- //senddetsommåsendes
		add_order := orderOnAuction(cost, receivedCostChan)
		if add_order {
			updatePrivateOrders(order)
			findNextFloorToGoTo()
		}
	}
}

func handleInternalOrder(order datatypes.InternalOrder) {
	converted_order_1 := datatypes.ExternalOrder{New_order: true, Executed_order: order.Executed_order, Floor: order.Floor, Direction: 1} //Her setter jeg New_order = true og Direction -1 som dummies fordi vi ikke for lov til å sette nil
	converted_order_2 := datatypes.ExternalOrder{New_order: true, Executed_order: order.Executed_order, Floor: order.Floor, Direction: -1}
	updatePrivateOrders(converted_order_1)
	updatePrivateOrders(converted_order_2)

	findNextFloorToGoTo()
}

func handleFinishedOrder(order datatypes.ExternalOrder, shareOrderChan chan datatypes.ExternalOrder) {
	updatePrivateOrders(order)
	updateSharedOrders(order)
	//shareOrderChan <- order
	findNextFloorToGoTo()
}

func updatePrivateOrders(new_order datatypes.ExternalOrder) {
	fmt.Println("Entered manager.updatePrivateOrders")
	private_orders_copy := make([]datatypes.ExternalOrder, len(private_orders))
	private_orders_copy = private_orders
	already_added := false

	if len(private_orders_copy) == 0 {
		private_orders_copy = append(private_orders_copy, new_order)
	}

	if new_order.Executed_order == true { //SLETT ordre fra private orders
		remaining_orders := make([]datatypes.ExternalOrder, 0)

		for items := range private_orders_copy {
			if private_orders_copy[items].Floor != new_order.Floor {
				remaining_orders = append(remaining_orders, private_orders_copy[items])
			}
		}

		private_orders = remaining_orders

	} else {
		for items := range private_orders_copy {
			if private_orders_copy[items] == new_order {
				already_added = true
			}
		} //LEGG TIL ordre slik som før
		if already_added == false {
			private_orders_copy = append(private_orders_copy, new_order)
		}
		private_orders = private_orders_copy
	}

	fmt.Println("Private orders are:", private_orders)
	backupOrders(private_orders)
}

func updateSharedOrders(new_order datatypes.ExternalOrder) { //Denne tar ikke hensyn til om en ordre er utført for å gå opp en etasje men sletter hele etasjen
	fmt.Println("Entered manager.updateSharedOrders")
	shared_orders_copy := make([]datatypes.ExternalOrder, 0)
	shared_orders_copy = shared_orders
	already_added := false

	if len(shared_orders_copy) == 0 {
		shared_orders_copy = append(shared_orders_copy, new_order)
	}

	if new_order.Executed_order == true { //SLETT ordre fra private orders
		remaining_orders := make([]datatypes.ExternalOrder, 0)

		for items := range shared_orders_copy {
			if shared_orders_copy[items].Floor != new_order.Floor {
				remaining_orders = append(remaining_orders, shared_orders_copy[items])
			}
		}
		shared_orders = remaining_orders

	} else {
		for items := range shared_orders_copy {
			if shared_orders_copy[items] == new_order {
				already_added = true
			}
		} //LEGG TIL ordre slik som før
		if already_added == false {
			shared_orders_copy = append(shared_orders_copy, new_order)
		}
		shared_orders = shared_orders_copy
	}
	fmt.Println("Shared orders are:", shared_orders)
}

func calculateCost(order datatypes.ExternalOrder) datatypes.CostInfo {
	cost := datatypes.CostInfo{Cost: 0, Floor: order.Floor, Direction: order.Direction}
	cost.Cost = int(math.Abs((float64(current_floor - order.Floor)))) //ha hvis current floor = -1? Bør vel aldri være -1 i orderManager -> filtrere bort -1 når man mottar på currentFloorToOrderManagerChan
	if int(direction) != order.Direction {
		cost.Cost += 2
	}
	return cost
}

func orderOnAuction(my_cost datatypes.CostInfo, receivedCostChan chan datatypes.CostInfo) bool {
	var external_cost datatypes.CostInfo
	for {
		select {
		case <-time.After(time.Millisecond * 500): //Tiden for timeout må finnes ved prøving/feiling
			fmt.Println("Auction timed out")
			return true
		case external_cost = <-receivedCostChan:
			if (external_cost.Floor == my_cost.Floor) && (external_cost.Direction == my_cost.Direction) {
				if external_cost.Cost < my_cost.Cost {
					return false
				} /*else if (external_cost.Cost == my_cost.Cost) && (external_cost.IP < my_cost.IP) {
					return false
				}*/
			}
		}
	}
}

func findNextFloorToGoTo() { //read only
	fmt.Println("Entered manager.findNextFloorToGoTo")
	//current_floor := <-currentFloorToOrderManagerChan
	original_direction := direction

	if len(private_orders) == 0 {
		fmt.Println("Private orders are empty")
		next_floor = -1
	} else {

		if direction == datatypes.UP {
			if checkUpFromCurrent(current_floor) {
				return
			}
			if checkDownFromTop(current_floor) {
				return
			}

			direction = datatypes.DOWN
		}

		if direction == datatypes.DOWN {
			if checkDownFromCurrent(current_floor) {
				return
			}
			if checkUpFromBottom(current_floor) {
				return
			}

			direction = datatypes.UP

			if direction != original_direction {
				if checkUpFromCurrent(current_floor) {
					return
				}
				if checkDownFromTop(current_floor) {
					return
				}
			}
		}
	}
}

func checkUpFromCurrent(current_floor int) bool {

	for checking_floor := current_floor; checking_floor <= number_of_floors; checking_floor++ {
		for _, order := range private_orders {
			if order.Floor == checking_floor && order.Direction == 1 {
				next_floor = checking_floor
				return true
			}
		}
	}
	return false
}

func checkDownFromTop(current_floor int) bool {
	for checking_floor := number_of_floors; checking_floor >= current_floor; checking_floor-- {
		for _, order := range private_orders {
			if order.Floor == checking_floor && order.Direction == -1 {
				next_floor = checking_floor
				return true
			}
		}
	}
	return false
}

func checkDownFromCurrent(current_floor int) bool {
	for checking_floor := current_floor; checking_floor >= 0; checking_floor-- {
		for _, order := range private_orders {
			if order.Floor == checking_floor && order.Direction == -1 {
				next_floor = checking_floor
				return true
			}
		}
	}
	return false
}

func checkUpFromBottom(current_floor int) bool {

	for checking_floor := 0; checking_floor <= current_floor; checking_floor++ {
		for _, order := range private_orders {
			if order.Floor == checking_floor && order.Direction == 1 {
				next_floor = checking_floor
				return true
			}
		}
	}
	return false
}

func restoreOrders(orders *[]datatypes.ExternalOrder) {
	buffer, _ := ioutil.ReadFile("privateOrdersBackup.txt")
	err := json.Unmarshal(buffer, orders)
	if err != nil {
		fmt.Println("error in restore order:", err)
	}
}

func backupOrders(orders []datatypes.ExternalOrder) {
	buffer, _ := json.Marshal(orders)
	err := ioutil.WriteFile("privateOrdersBackup.txt", []byte(buffer), 0777)
	if err != nil {
		log.Fatal(err)
	}
}

func createOrderBackupFile() {
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
