package datatypes

type InternalOrder struct {
	floor int
}

type ExternalOrder struct {
	new_order      bool
	executed_order bool
	floor          int
	direction      int
}

type CostInfo struct {
	cost      int
	floor     int
	direction int
}
