package datatypes

//const N_ELEVATORS = 1
const M_FLOORS int = 4

//ExternalOrderArray bool := [2][4]bool{} // Vi vil ha noe sånt men ikke som en global variabel! 
//InternalOrderArray := [4]bool{} // Vi vil ha noe sånt men ikke som en global variabel! 

var ExternalOrderArray [2][M_FLOORS] bool
var InternalOrderArray [M_FLOORS] bool

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
