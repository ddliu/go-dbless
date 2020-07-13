package dbless

import (
	"strings"
)

type Filter map[string]interface{}

func (f Filter) GetWhere(db *DB) (string, []interface{}) {
	var where []string
	var params []interface{}
	for k, v := range f {
		where = append(where, k)
		params = append(params, v)
	}

	placeholders := db.Driver.Placeholder(params)
	for i, placeholder := range placeholders {
		where[i] = db.Driver.QuoteIdentifier(where[i]) + " = " + placeholder
	}

	return strings.Join(where, " AND "), params
}
