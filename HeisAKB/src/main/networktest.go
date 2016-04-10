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

	
	/*
	1. lag tråd som lytter til receivedOrderChan for alltid
			- Når noe kommer på kanalen
			- Se om det allerede ligger inne
				- Om det ligger inne, ikke gjør noe
				- Om den ikke ligger inne
					- Legg til i orderslice
	
	*/

	//Deklarerer slice til bruk med ordre
	orders := make([]datatypes.ExternalOrder, 0) // Denne har null som argument fordi den er tom

	//Skriv ut alle meldingene som kommer inn på kanalen recievedOrderChan mens programmet lever
	go func () {
		for {
			//Lagre melding til en variabel
			neworder := <-receivedOrderChan

			//Siden orders er tom må vi ha en sjekk for å se om vi skal lege til første element
			//Hvis den ikke er tom søker vi iterativt gjennom den for å se ettter lignende ordre
			if len(orders) == 0{
				orders = append(orders, neworder)
			} else {
				for items := range orders {
					if orders[items] != neworder  {
						orders = append(orders, neworder)
					} 
				}
			}
		fmt.Println("Mottatt meldling er:")
		fmt.Println(neworder)
		}
	}()


	//Send en testordre
	ex_order := datatypes.ExternalOrder{New_order:false, Executed_order: true, Floor : 2, Direction: 1}
	network.BroadcastExternalOrder(ex_order)


	time.Sleep(time.Second * 2)

	//Send en testordre til
	ex_order1 := datatypes.ExternalOrder{New_order:true, Executed_order: true, Floor : 2, Direction: 1}
	network.BroadcastExternalOrder(ex_order1)

	time.Sleep(time.Second * 2)

	fmt.Println(orders)





	fmt.Println("Lager/oppdaterer fil til heisordre")

	var ( //Kan flyttes øverst tror jeg
    	newFile *os.File
   		err     error
	)
	//Lager ny fil
	newFile, err = os.Create("ordrebackup.txt")
    if err != nil {
        log.Fatal(err)
    }
    log.Println(newFile)
    newFile.Close()

    //Skriver med Quick write to file 
    fmt.Println(orders)
    //Konverterer til json
    buffer, _ := json.Marshal(orders)
    //Skriver til fil
    err = ioutil.WriteFile("ordrebackup.txt", []byte(buffer), 0777)
    if err != nil {
        log.Fatal(err)
    }
    //Tester å gjennopprette fra fil så jeg sletter slice
    fmt.Println("Tester å hente inn igjen data fra fil")
    fmt.Println("Setter orders til 0: ")
    orders = make([]datatypes.ExternalOrder, 0)
    fmt.Println(orders)

    //Decoder fil og legger inn igjen i slices
    buffer, _ = ioutil.ReadFile("ordrebackup.txt")
    err = json.Unmarshal(buffer, &orders)
    if err != nil {
		fmt.Println("error:", err)
	}
    fmt.Println("Skriver ut info hentet fra fil")
    fmt.Println(orders)

	

}