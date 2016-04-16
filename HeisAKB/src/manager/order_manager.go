package manager

import (
	"datatypes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	//"network"
	"os"
	"time"
)

var next_floor int
var direction datatypes.Direction
var number_of_floors int
var current_floor int

var shared_orders []datatypes.ExternalOrder
var private_orders []datatypes.ExternalOrder

const ORDER_TIMELIMIT = 20

func InitOrderManager(n_FLOORS int, newInternalOrderChan chan datatypes.InternalOrder,
	newExternalOrderChan chan datatypes.ExternalOrder, currentFloorToOrderManagerChan chan int,
	orderFinishedChan chan int, dirChan chan datatypes.Direction, receivedOrderChan chan datatypes.ExternalOrder,
	receivedCostChan chan datatypes.CostInfo, shareOrderChan chan datatypes.ExternalOrder,
	shareCostChan chan datatypes.CostInfo, nextFloorChan chan int, setInternalLightsChan chan datatypes.InternalOrder,
	setExternalLightsChan chan datatypes.ExternalOrder) {

	shared_orders = make([]datatypes.ExternalOrder, 0)
	private_orders = make([]datatypes.ExternalOrder, 0)

	if _, err := os.Stat("privateOrdersBackup.txt"); os.IsNotExist(err) {
		createOrderBackupFile()
	}

	next_floor = -1
	number_of_floors = n_FLOORS
	direction = datatypes.DOWN
	rand.Seed(time.Now().UTC().UnixNano())

	restoreOrders(&private_orders)
	fmt.Println("Private orders from file are:", private_orders)
	findNextFloorToGoTo()
	fmt.Println("Init direction is: ", direction)

	go orderManager(newInternalOrderChan, newExternalOrderChan, currentFloorToOrderManagerChan, orderFinishedChan, dirChan, receivedOrderChan, receivedCostChan,
		shareOrderChan, shareCostChan, nextFloorChan, setInternalLightsChan, setExternalLightsChan) //mangler setExternalLightsChan
	go handleExpiredOrder(newExternalOrderChan, orderFinishedChan)
}

func orderManager(newInternalOrderChan chan datatypes.InternalOrder, newExternalOrderChan chan datatypes.ExternalOrder, currentFloorToOrderManagerChan chan int, orderFinishedChan chan int, dirChan chan datatypes.Direction, receivedOrderChan chan datatypes.ExternalOrder, receivedCostChan chan datatypes.CostInfo, shareOrderChan chan datatypes.ExternalOrder, shareCostChan chan datatypes.CostInfo, nextFloorChan chan int, setInternalLightsChan chan datatypes.InternalOrder, setExternalLightsChan chan datatypes.ExternalOrder) {
	for {
		select {
		case new_internal_order := <-newInternalOrderChan:
			handleInternalOrder(new_internal_order)
			setInternalLightsChan <- new_internal_order

		case new_external_order := <-newExternalOrderChan:
				received_order := false
				handleSharedOrder(new_external_order, received_order, shareCostChan, receivedCostChan, shareOrderChan)
				setExternalLightsChan <- new_external_order
		
		case received_network_order := <-receivedOrderChan:
			received_order := true
			handleSharedOrder(received_network_order, received_order, shareCostChan, receivedCostChan, shareOrderChan)
			if received_network_order.Executed_order{
				external_order_up := datatypes.ExternalOrder{New_order: false, Executed_order: true, Floor: received_network_order.Floor, Direction: 1}
				external_order_down := datatypes.ExternalOrder{New_order: false, Executed_order: true, Floor: received_network_order.Floor, Direction: -1}
				setExternalLightsChan <- external_order_up
				setExternalLightsChan <-  external_order_down
			} else {
				setExternalLightsChan <- received_network_order
			}

		case finished_order := <-orderFinishedChan:
			external_order_up := datatypes.ExternalOrder{New_order: false, Executed_order: true, Floor: finished_order, Direction: 1}
			external_order_down := datatypes.ExternalOrder{New_order: false, Executed_order: true, Floor: finished_order, Direction: -1}
			internal_order := datatypes.InternalOrder{Executed_order: true, Floor: finished_order}

			handleFinishedOrder(external_order_up, shareOrderChan)
			setExternalLightsChan <- external_order_down
			setInternalLightsChan <- internal_order
			setExternalLightsChan <- external_order_up

		case update_current_floor := <-currentFloorToOrderManagerChan:
			if update_current_floor != -1 {
				current_floor = update_current_floor
			} else {
				if direction == datatypes.UP {
					current_floor++
				} else if direction == datatypes.DOWN {
					current_floor--
				} else {
					fmt.Println("manager.Case: update_current_floor: Invalid direction")
				}
			}

		case update_direction := <-dirChan:
			if update_direction != datatypes.STOP {
				direction = update_direction
			}

		case <-time.After(250 * time.Millisecond):
			nextFloorChan <- next_floor
		}
	}
}

func handleSharedOrder(order datatypes.ExternalOrder, received_order bool, shareCostChan chan datatypes.CostInfo, receivedCostChan chan datatypes.CostInfo, shareOrderChan chan datatypes.ExternalOrder) {
	//mangler setExternalLightsChan og newInternalOrderChan
	order.Timestamp = time.Now().Unix()
	if !(received_order) {
		fmt.Println("Manager.handlesharedorder: sharing order")
		shareOrderChan <- order
	}
	updateSharedOrders(order)

	if !(order.Executed_order) {
		cost := calculateCost(order)
		fmt.Println("Cost blaaaah: ", cost.Cost)
		fmt.Println(" ")
		shareCostChan <- cost
		//setExternalLightsChan <- //senddetsommåsendes
		add_order := orderOnAuction(cost, receivedCostChan)
		if add_order {
			updatePrivateOrders(order)
			findNextFloorToGoTo()
		}
	}
}

func handleInternalOrder(order datatypes.InternalOrder) {
	converted_order_1 := datatypes.ExternalOrder{New_order: true, Executed_order: order.Executed_order, Floor: order.Floor, Direction: 1, Timestamp: (time.Now().Unix() + 3600)} //Her setter jeg New_order = true og Direction -1 som dummies fordi vi ikke for lov til å sette nil
	converted_order_2 := datatypes.ExternalOrder{New_order: true, Executed_order: order.Executed_order, Floor: order.Floor, Direction: -1, Timestamp: (time.Now().Unix() + 3600)}
	updatePrivateOrders(converted_order_1)
	updatePrivateOrders(converted_order_2)

	findNextFloorToGoTo()
}

func handleFinishedOrder(order datatypes.ExternalOrder, shareOrderChan chan datatypes.ExternalOrder) {
	updatePrivateOrders(order)
	updateSharedOrders(order)
	shareOrderChan <- order
	findNextFloorToGoTo()
}

func handleExpiredOrder(newExternalOrderChan chan datatypes.ExternalOrder, orderFinishedChan chan int){
	for{
		select{
			case <- time.After(5 * time.Second):
				shared_orders_copy := make([]datatypes.ExternalOrder, 0)
				shared_orders_copy = shared_orders

				for items := range shared_orders_copy {
					if (time.Now().Unix() - shared_orders_copy[items].Timestamp) > ORDER_TIMELIMIT {
						order := shared_orders_copy[items]
						//si til andre at den er håndert, kall handleFinishedOrder
						order.Executed_order = true
						orderFinishedChan <- order.Floor

						//del med andre heiser, newExternalOrderChan
						order.Executed_order = false
						order.New_order = true
						newExternalOrderChan <- order
					}
			}
		}
	}
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
		/*for items := range private_orders_copy {
			if private_orders_copy[items] == new_order {
				already_added = true
			}
		}*/ //LEGG TIL ordre slik som før
		for items := range private_orders_copy {
			if private_orders_copy[items].Floor == new_order.Floor && private_orders_copy[items].Direction == new_order.Direction {
				already_added = true
			}
		}
		if already_added == false {
			private_orders_copy = append(private_orders_copy, new_order)
		}
		private_orders = private_orders_copy
	}

	fmt.Println("Private orders are:", private_orders)
	backupOrders(private_orders)
}

func updateSharedOrders(new_order datatypes.ExternalOrder) { //Denne tar ikke hensyn til om en ordre er utført for å gå opp en etasje men sletter hele etasjen
	//fmt.Println("Entered manager.updateSharedOrders")
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
			if shared_orders_copy[items].Floor == new_order.Floor && shared_orders_copy[items].Direction == new_order.Direction {
				already_added = true
			}
		}
		if already_added == false {
			shared_orders_copy = append(shared_orders_copy, new_order)
		}
		shared_orders = shared_orders_copy
	}
	fmt.Println("Shared orders are:", shared_orders)
}

func calculateCost(order datatypes.ExternalOrder) datatypes.CostInfo {
	cost := datatypes.CostInfo{Cost: 0, Floor: order.Floor, Direction: order.Direction, ID: rand.Intn(int(time.Now().UnixNano()))}
	cost.Cost = current_floor - order.Floor 
	if cost.Cost < 0 {
		if int(direction) != order.Direction {
			cost.Cost = int(math.Abs(float64(cost.Cost))) + 2
		} else if direction == datatypes.UP && order.Direction == 1 {
			cost.Cost = int(math.Abs(float64(cost.Cost)))
		} else {
			cost.Cost = int(math.Abs(float64(cost.Cost))) + 4
		}
	} else if cost.Cost > 0 {
		if int(direction) != order.Direction {
			cost.Cost = int(math.Abs(float64(cost.Cost))) + 2
		} else if direction == datatypes.UP && order.Direction == 1 {
			cost.Cost = int(math.Abs(float64(cost.Cost))) + 4
		} else {
			cost.Cost = int(math.Abs(float64(cost.Cost)))
		}
	}
	return cost
}

func orderOnAuction(my_cost datatypes.CostInfo, receivedCostChan chan datatypes.CostInfo) bool {
	var external_cost datatypes.CostInfo
	for {
		select {
		case <-time.After(time.Millisecond * 500): //Tiden for timeout må finnes ved prøving/feiling
			fmt.Println("  ")
			fmt.Println("Order won!")
			fmt.Println("Auction timed out")
			fmt.Println("  ")
			return true
		case external_cost = <-receivedCostChan:
			fmt.Println("Evaluating external cost", external_cost)
			if (external_cost.Floor == my_cost.Floor) && (external_cost.Direction == my_cost.Direction) {
				if external_cost.Cost < my_cost.Cost {
					return false
				} else if (external_cost.Cost == my_cost.Cost) && (external_cost.ID < my_cost.ID) {
					return false
				}
			}
		}
	}
}

func findNextFloorToGoTo() { //read only
	//fmt.Println("Entered manager.findNextFloorToGoTo")
	original_direction := direction

	if len(private_orders) == 0 {
		fmt.Println("Private orders are empty")
		next_floor = -1
	} else {

		if direction == datatypes.UP {

			if checkUpOrdersAboveCurrentFloor(current_floor) {
				return
			}
			if checkDownOrdersFromTopFloor(current_floor) {
				return
			}

			direction = datatypes.DOWN
		}

		if direction == datatypes.DOWN {
			if checkDownOrdersBelowCurrentFloor(current_floor) {
				return
			}
			if checkUpOrdersFromBottomFloor(current_floor) {
				return
			}

			direction = datatypes.UP

			if direction != original_direction {
				if checkUpOrdersAboveCurrentFloor(current_floor) {
					return
				}
				if checkDownOrdersFromTopFloor(current_floor) {
					return
				}
			}
		}
	}
}

func checkUpOrdersAboveCurrentFloor(current_floor int) bool {

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

func checkDownOrdersBelowCurrentFloor(current_floor int) bool {
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

func checkDownOrdersFromTopFloor(current_floor int) bool {
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

func checkUpOrdersFromBottomFloor(current_floor int) bool {

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
