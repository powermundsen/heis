package main

import (
	"fmt"
	"datatypes"
	"ioHandling"
)

const n_FLOORS int = 4

func main(){

//Arraytesting
datatypes.OrdersSlicesInit(n_FLOORS)
fmt.Println(datatypes.ExternalOrdersSlice)
fmt.Println(datatypes.InternalOrdersSlice)


// Påfølgende 10 linjer er for å teste io_manager. Status nå: greier ikke kjøre funksjoner fra io_manager i main, feilmelding sier at funksjonene er undefined..
//newInternalOrderChan := make(chan datatypes.InternalOrder)
//newExternalOrderChan := make(chan datatypes.ExternalOrder)
//currentFloorToOrderManagerChan := make(chan int)
//currentFloorToElevControllerChan := make(chan int)
//setInternalLightsChan := make(chan int)
//setExternalLightsChan := make(chan int)

//InitIo(newInternalOrderChan, newExternalOrderChan, currentFloorToOrderManagerChan, currentFloorToElevControllerChan, setInternalLightsChan, setExternalLightsChan)
//io.ioManager(newInternalOrderChan, newExternalOrderChan, currentFloorToOrderManagerChan, currentFloorToElevControllerChan, setInternalLightsChan, setExternalLightsChan)
//ReadAllInternalButtons()

}