package main
import (
	"driver"
	"fmt"
)

//Git test

func main(){

	driver.Elevator_init()


	fmt.Println("Starting up the elevator")
	
	fmt.Println("Elevator init finished")
	fmt.Println("Finding closest floor")


	
	//Arrive at floor (Hent inn hvilken etasje - Send etasje til event manager som setter lys osv)


	//Registrer etasje
	//Skriv ut etasje og lagre i minnet
	driver.Elevator_set_floor_indicator(3)
	
}