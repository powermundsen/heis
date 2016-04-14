package main

import (
	//"driver"
	"network"
	"fmt"
	"datatypes"
	//"manager"
	"runtime"
	//"controller"
	"time"
	"os" //For å lage filer (http://www.devdungeon.com/content/working-files-go#create_empty_file)
	"log" //For feilmeldinger
	"io/ioutil" //For å raskt skrive til fil
	"encoding/json"
)


func main(){


	runtime.GOMAXPROCS(runtime.NumCPU())

	shareOrderChan	 	:= make(chan datatypes.ExternalOrder)
	receivedOrderChan	:= make(chan datatypes.ExternalOrder)
	shareCostChan		:= make(chan network.CostInfo)
	receivedCostChan	:= make(chan network.CostInfo)

	network.InitNetworkHandler(shareOrderChan, receivedOrderChan, shareCostChan, receivedCostChan)


	//Deklarerer slice til bruk med ordre
	orders := make([]datatypes.ExternalOrder, 0) // Denne har null som argument fordi den er tom

	//Skriv ut alle meldingene som kommer inn på kanalen recievedOrderChan mens programmet lever
	go func () {
		for {
			addNewExternalOrdersToSlice(receivedOrderChan, &orders)
		}
	}()


	//Send en testordre
	for i := 0; i < 5; i++ {
		ex_order := datatypes.ExternalOrder{New_order:false, Executed_order: true, Floor : 2, Direction: 1}
		network.BroadcastExternalOrder(ex_order)
		time.Sleep(time.Millisecond * 1)
	}
	//Send en testordre til
	ex_order1 := datatypes.ExternalOrder{New_order:true, Executed_order: true, Floor : 2, Direction: 1}
	network.BroadcastExternalOrder(ex_order1)

	time.Sleep(time.Second * 2)

	fmt.Println(orders)

	fmt.Println("Lager/oppdaterer fil til heisordre")

	//createOrderBackupFile()	//Den er muligvis unødvendig. WriteFile i backupOrders lager også en fil
	backupOrders(orders)
	fmt.Println("Tester å hente inn igjen data fra fil")
    fmt.Println("Setter orders til 0: ")
    orders = make([]datatypes.ExternalOrder, 0)
    fmt.Println(orders)
	restoreOrders(&orders)
	fmt.Println("Skriver ut info hentet fra fil")
    fmt.Println(orders)
    

}

func restoreOrders(orders *[]datatypes.ExternalOrder){
    buffer, _ := ioutil.ReadFile("orderbackup.txt")
    err := json.Unmarshal(buffer, orders)
    if err != nil {
		fmt.Println("error:", err)
	}
}

func backupOrders(orders []datatypes.ExternalOrder){
    fmt.Println(orders)
    buffer, _ := json.Marshal(orders)
    err := ioutil.WriteFile("orderbackup.txt", []byte(buffer), 0777)
    if err != nil {
        log.Fatal(err)
    }
}

func createOrderBackupFile(){
	var (
    	newFile *os.File
   		err     error
	)

	newFile, err = os.Create("orderbackup.txt")
    if err != nil {
        log.Fatal(err)
    }

    log.Println(newFile)
    newFile.Close()
}

func addNewExternalOrdersToSlice(receivedOrderChan chan datatypes.ExternalOrder, orders *[]datatypes.ExternalOrder){
	neworder := <-receivedOrderChan

	ordersCopy := make([]datatypes.ExternalOrder, 0)		//Kopierer referert slice for å arbeide med den
	ordersCopy = *orders

	if len(ordersCopy) == 0{
		ordersCopy= append(ordersCopy, neworder)

	} else {
		for items := range ordersCopy {
			if ordersCopy[items] != neworder  {
				ordersCopy = append(ordersCopy, neworder)
			} 
		}
	}
	*orders = ordersCopy
	fmt.Println("Mottatt meldling er:")
	fmt.Println(neworder)
}