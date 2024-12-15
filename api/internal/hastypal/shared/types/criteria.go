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
	Table interface{}
	Field string
}

type Criteria struct {
	Filters []Filter
	Join    []Relation
}

func NewCriteria() Criteria {
	return Criteria{}
}

func (c Criteria) Equal(fieldName string, value any) Criteria {
	c.Filters = append(c.Filters, Filter{Name: fieldName, Operand: constants.Equal, Value: value})

	return c
}

func (c Criteria) GreaterThanOrEqual(fieldName string, value any) Criteria {
	c.Filters = append(c.Filters, Filter{Name: fieldName, Operand: constants.GreaterThanOrEqual, Value: value})

	return c
}

func (c Criteria) LessThanOrEqual(fieldName string, value any) Criteria {
	c.Filters = append(c.Filters, Filter{Name: fieldName, Operand: constants.LessThanOrEqual, Value: value})

	return c
}

func (c Criteria) LeftJoin(withTable interface{}, onField string) Criteria {
	c.Join = append(c.Join, Relation{Type: constants.LeftJoin, Table: withTable, Field: onField})

	return c
}
