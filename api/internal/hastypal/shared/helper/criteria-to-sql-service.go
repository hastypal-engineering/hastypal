package helper

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"reflect"
	"strings"
)

type CriteriaToSqlService struct {
	table string
}

func NewCriteriaToSqlService(entity interface{}) (*CriteriaToSqlService, error) {
	v := reflect.ValueOf(entity)

	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil, exception.
			New("Entity provided is not a struct pointer").
			Trace("reflect.ValueOf", "criteria-to-sql-service.go")
	}

	structValue := v.Elem()
	structType := structValue.Type()

	var scanTargets []interface{}
	var columnNames []string

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)

		columnName := fieldType.Tag.Get("db")

		if columnName == "" {
			continue
		}

		if field.CanAddr() {
			scanTargets = append(scanTargets, field.Addr().Interface())
			columnNames = append(columnNames, columnName)
		}
	}

	return &CriteriaToSqlService{
		table: CamelToSnake(structType.Name()),
	}, nil
}

func (c *CriteriaToSqlService) Transform(criteria types.Criteria) (string, error) {
	if len(criteria.Filters) == 0 {
		return "SELECT * FROM" + " " + c.table, nil
	}

	if len(criteria.Join) > 0 {
		joinClause := c.constructJoinClause(criteria)

		sql := "SELECT * FROM " + c.table + " " + joinClause + " WHERE "

		completeSQL := sql + c.constructWhereClause(criteria)

		return completeSQL, nil
	}

	sql := "SELECT * FROM " + c.table + " WHERE "

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
		relationTableAlias := relation.Table[:3]

		join = append(join, fmt.Sprintf(
			"%s JOIN %s %s ON %s.id = %s.%s",
			relation.Type,
			relation.Table,
			relationTableAlias,
			c.table,
			relationTableAlias,
			relation.Field,
		))
	}

	return strings.Join(join, " ")
}
