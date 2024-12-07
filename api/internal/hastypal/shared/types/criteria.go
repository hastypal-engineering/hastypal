package types

type Filter struct {
	Name    string
	Operand string
	Value   any
}

type Relation struct {
	Table string
	Field string
}

type Criteria struct {
	Filters []Filter
	Join    []Relation
}
