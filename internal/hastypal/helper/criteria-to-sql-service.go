package helper

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/types"
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
		_, ok := filter.Value.(string)
		if ok {
			sqlStringValue := fmt.Sprintf("'%s'", filter.Value)

			clause := filter.Name + " " + filter.Operand + " " + sqlStringValue

			where = append(where, clause)

			continue
		}

		sqlStringValue := fmt.Sprintf("%s", filter.Value)

		clause := filter.Name + " " + filter.Operand + " " + sqlStringValue

		where = append(where, clause)
	}

	completeSQL := sql + strings.Join(where, " AND ")

	return completeSQL, nil
}
