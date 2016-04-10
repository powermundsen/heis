package main

import (
	"fmt"
	"datatypes"
	"ioHandling"
	"driver"
	"time"
)

const n_FLOORS int = 4

func main(){

//Arraytesting
datatypes.OrdersSlicesInit(n_FLOORS)
fmt.Println(datatypes.ExternalOrdersSlice)
fmt.Println(datatypes.InternalOrdersSlice)

driver.Elevator_set_motor_direction(-1)
// Påfølgende 10 linjer er for å teste io_manager. Status nå: greier ikke kjøre funksjoner fra io_manager i main, feilmelding sier at funksjonene er undefined..
newInternalOrderChan := make(chan datatypes.InternalOrder)
newExternalOrderChan := make(chan datatypes.ExternalOrder)
currentFloorToOrderManagerChan := make(chan int)
currentFloorToElevControllerChan := make(chan int)
setInternalLightsChan := make(chan [] int)
setExternalLightsChan := make(chan [] int)

ioHandling.InitIo(n_FLOORS, newInternalOrderChan, newExternalOrderChan, currentFloorToOrderManagerChan, currentFloorToElevControllerChan, setInternalLightsChan, setExternalLightsChan)
//io.ioManager(newInternalOrderChan, newExternalOrderChan, currentFloorToOrderManagerChan, currentFloorToElevControllerChan, setInternalLightsChan, setExternalLightsChan)
//ReadAllInternalButtons()

time.Sleep(10*time.Second)
}