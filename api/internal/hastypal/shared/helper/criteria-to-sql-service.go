package helper

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"strings"
)

type CriteriaToSqlService struct {
	Reflection *ReflectionHelper
	Table      string
	Fields     []interface{}
}

func NewCriteriaToSqlService(entity interface{}) (*CriteriaToSqlService, error) {
	reflection := NewReflectionHelper()

	fields, _ := reflection.ExtractDatabaseFields(entity)
	table, _ := reflection.ExtractTableName(entity)

	return &CriteriaToSqlService{
		Table:      table,
		Fields:     fields,
		Reflection: reflection,
	}, nil
}

func (c *CriteriaToSqlService) Transform(criteria types.Criteria) (string, error) {
	if len(criteria.Filters) == 0 {
		return "SELECT * FROM" + " " + c.Table, nil
	}

	if len(criteria.Join) > 0 {
		joinClause := c.constructJoinClause(criteria)

		sql := "SELECT * FROM " + c.Table + " " + joinClause + " WHERE "

		completeSQL := sql + c.constructWhereClause(criteria)

		return completeSQL, nil
	}

	sql := "SELECT * FROM " + c.Table + " WHERE "

	completeSQL := sql + c.constructWhereClause(criteria)

	return completeSQL, nil
}

func (c *CriteriaToSqlService) constructWhereClause(criteria types.Criteria) string {
	var where []string

	for _, filter := range criteria.Filters {
		_, isStringValue := filter.Value.(string)

		if isStringValue {
			sqlStringValue := fmt.Sprintf("'%s'", filter.Value)

			clause := filter.Name + " " + filter.Operand + " " + sqlStringValue

			where = append(where, clause)

			continue
		}

		_, isIntValue := filter.Value.(int)

		if isIntValue {
			sqlIntValue := fmt.Sprintf("%d", filter.Value)

			clause := filter.Name + " " + filter.Operand + " " + sqlIntValue

			where = append(where, clause)

			continue
		}

		_, isFloatValue := filter.Value.(float64)

		if isFloatValue {
			sqlIntValue := fmt.Sprintf("%d", filter.Value)

			clause := filter.Name + " " + filter.Operand + " " + sqlIntValue

			where = append(where, clause)

			continue
		}

		sqlBooleanValue := fmt.Sprintf("%t", filter.Value)

		clause := filter.Name + " " + filter.Operand + " " + sqlBooleanValue

		where = append(where, clause)
	}

	return strings.Join(where, " AND ")
}

func (c *CriteriaToSqlService) constructJoinClause(criteria types.Criteria) string {
	var join []string

	for _, relation := range criteria.Join {
		relationTableName, _ := c.Reflection.ExtractTableName(relation.Table)
		relationTableAlias := relationTableName[:3]

		joinTableFields, _ := c.Reflection.ExtractDatabaseFields(relation.Table)

		c.Fields = append(c.Fields, joinTableFields...)

		join = append(join, fmt.Sprintf(
			"%s JOIN %s %s ON %s.id = %s.%s",
			relation.Type,
			relationTableName,
			relationTableAlias,
			c.Table,
			relationTableAlias,
			relation.Field,
		))
	}

	return strings.Join(join, " ")
}
