package types

type Filter struct {
	Name    string
	Operand string
	Value   any
}

type Criteria struct {
	Filters []Filter
}
