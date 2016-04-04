package datatypes

ExternalOrderArray bool := [2][4]bool{} // Vi vil ha noe sånt men ikke som en global variabel! 
InternalOrderArrat := [4]bool{} // Vi vil ha noe sånt men ikke som en global variabel! 

type ExternalOrder struct {
	floor1button0 bool
	floor1button1 bool
	floor2button0 bool
	floor2button1 bool

	// Vi må finne en løsning på hvordan en struct kan ha et array inni seg

}
