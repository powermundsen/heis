package datatypes

//ExternalOrderArray int := [2][4]int{} // Vi vil ha noe sånt men ikke som en global variabel! 
//InternalOrderArray := [4]int{} // Vi vil ha noe sånt men ikke som en global variabel! 


var ExternalOrdersSlice [] int
var InternalOrdersSlice [] int

func OrdersSlicesInit(n_FLOORS int) {
	InternalOrdersSlice = make([]int, n_FLOORS)
	ExternalOrdersSlice = make([]int, 2 * n_FLOORS)

	for i := range(InternalOrdersSlice){
		ExternalOrdersSlice[i] = 0
		ExternalOrdersSlice[n_FLOORS+i] = 0
		InternalOrdersSlice[i] = 1
		
	}
}


/*
type ExternalOrder struct {
	floor1button0 int
	floor1button1 int
	floor2button0 int
	floor2button1 int
	//mangler resten av etasjene

	// Vi må finne en løsning på hvordan en struct kan ha et array inni seg

}
*/
