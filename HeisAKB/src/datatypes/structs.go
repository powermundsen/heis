package datatypes

type InternalOrder struct {
	Floor int
}

//Jeg endret parametrene i strukten til å ha stor bokstav, da kan vi deklarere nye instanser slik
//	ex_order := datatypes.ExternalOrder{New_order:false, Executed_order: true, Floor : 2, Direction: 1}
type ExternalOrder struct {
	New_order      bool
	Executed_order bool
	Floor          int
	Direction      int
}

type CostInfo struct {
	Cost      int
	Floor     int
	Direction int
}
