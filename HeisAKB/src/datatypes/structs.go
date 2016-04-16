package datatypes

type InternalOrder struct {
	Executed_order bool
	Floor          int
}

type ExternalOrder struct {
	New_order      bool
	Executed_order bool
	Floor          int
	Direction      int
	Timestamp		int64
}

type CostInfo struct {
	Cost      int
	Floor     int
	Direction int
	ID        int
}
