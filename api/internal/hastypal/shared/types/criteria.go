package types

import (
	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
)

type Filter struct {
	Name    string
	Operand string
	Value   any
}

type Relation struct {
	Type  string
	Table string
	Field string
}

type Criteria struct {
	Filters []Filter
	Join    []Relation
}

func NewCriteria() *Criteria {
	return &Criteria{}
}

func (c *Criteria) Equal(fieldName string, value any) *Criteria {
	c.Filters = append(c.Filters, Filter{Name: fieldName, Operand: constants.Equal, Value: value})

	return c
}

func (c *Criteria) LeftJoin(leftTable interface{}, rightTable interface{}) *Criteria {
	c.Join = append(c.Join, Relation{Type: constants.LeftJoin, Table: "withTable", Field: "onField"})

	return c
}
