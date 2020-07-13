package dbless

import (
	"fmt"
	"strings"
)

type Filter map[string]interface{}

func (f Filter) GetWhere(db *DB) (string, []interface{}) {
	var where []string
	var params []interface{}
	for k, v := range f {
		where = append(where, fmt.Sprintf("%s = ?", db.Driver.QuoteIdentifier(k)))
		params = append(params, v)
	}

	return strings.Join(where, " AND "), params
}
