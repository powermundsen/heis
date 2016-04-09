package datatypes

//const N_ELEVATORS = 1
//ExternalOrderArray bool := [2][4]bool{} // Vi vil ha noe sånt men ikke som en global variabel! 
//InternalOrderArray := [4]bool{} // Vi vil ha noe sånt men ikke som en global variabel! 


// Arrays må ha bestemt størrelse, må skrives inn med tall
// Slices kan bestemmes til å ha n_FLOORS størrelse

var ExternalOrdersArray [2][4] bool
var InternalOrdersArray [4] bool

//func ordersArrayConfig(n_FLOORS int) {
//	var ExternalOrdersArray [2][n_FLOORS] bool
//	var InternalOrdersArray [n_FLOORS] bool
//}


/*
type ExternalOrder struct {
	floor1button0 bool
	floor1button1 bool
	floor2button0 bool
	floor2button1 bool
	//mangler resten av etasjene

	// Vi må finne en løsning på hvordan en struct kan ha et array inni seg

}
*/
