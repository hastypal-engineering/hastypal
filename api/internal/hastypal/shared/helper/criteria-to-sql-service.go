package helper

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"strings"
)

type CriteriaToSqlService struct {
	table string
}

func NewCriteriaToSqlService(table string) *CriteriaToSqlService {
	return &CriteriaToSqlService{
		table: table,
	}
}

func (c *CriteriaToSqlService) Transform(criteria types.Criteria) (string, error) {
	if len(criteria.Filters) == 0 {
		return "SELECT * FROM" + " " + c.table, nil
	}

	var where []string

	sql := "SELECT * FROM" + " " + c.table + " WHERE "

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

	completeSQL := sql + strings.Join(where, " AND ")

	return completeSQL, nil
}
