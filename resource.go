// Copyright 2019 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package dbless

import (
	"fmt"
	"time"

	"github.com/spf13/cast"
)

type Field struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Meta        string `json:"meta"`
}

func (r Record) ID() string {
	id, ok := r["id"]
	if !ok {
		return ""
	}

	return cast.ToString(id)
}

type Resource struct {
	Name      string  `json:"name"`
	Title     string  `json:"title"`
	Fields    []Field `json:"fields"`
	Timestamp bool    `json:"timestamp"`
	Meta      string  `json:"meta"`
	DB        *DB     `json:"-"`
}

func (r Resource) Get(filter interface{}) (Record, error) {
	f, ok := filter.(Filter)
	if !ok {
		f = Filter{
			"id": filter,
		}
	}

	where, params := f.GetWhere(r.DB)

	row, err := r.DB.GetRow("select * from "+r.DB.Driver.QuoteIdentifier(r.Name)+" where "+where, params...)

	if err != nil {
		return nil, err
	}

	return row, nil
}

func (r Resource) Save(record Record, filter ...Filter) (string, error) {
	activeFilter := Filter{}

	if len(filter) > 0 {
		for _, f := range filter {
			for k, v := range f {
				activeFilter[k] = v
			}
		}
	}

	if len(activeFilter) > 0 {
		existing, err := r.Get(activeFilter)
		if IsRecordNotFound(err) {

		} else if err != nil {
			return "", err
		} else {
			record["id"] = existing["id"]
		}
	}

	id := record.ID()

	delete(record, "id")
	delete(record, "created_at")
	delete(record, "updated_at")
	// update
	if id != "" {
		if r.Timestamp {
			record["updated_at"] = time.Now()
		}

		placeholders := r.DB.Driver.Placeholder([]interface{}{id})
		_, err := r.DB.Update(r.Name, record, "id = "+placeholders[0], id)
		return id, err
	} else {
		if r.Timestamp {
			record["created_at"] = time.Now()
			record["updated_at"] = time.Now()
		}

		id, err := r.DB.InsertGetID(r.Name, record)

		return id, err
	}
}

type ListInput struct {
	Pagination Pagination `json: "pagination"`
	Filter     Filter     `json:"filter"`
}

type ListOutput struct {
	Pagination
	List []Record `json:"list"`
}

func (r Resource) List(input ListInput) (*ListOutput, error) {
	var params []interface{}
	var whereStr string
	if input.Filter != nil {
		whereStr, params = input.Filter.GetWhere(r.DB)
	}

	if whereStr != "" {
		whereStr = " WHERE " + whereStr
	}

	sqlTotal := "SELECT COUNT(*) FROM " + r.DB.Driver.QuoteIdentifier(r.Name) + whereStr
	total, err := r.DB.Scalar(sqlTotal, params...)
	if err != nil {
		return nil, err
	}

	input.Pagination.SetTotal(uint(total))
	input.Pagination.Valid()

	var list []Record
	if total != 0 {
		offset, limit := input.Pagination.GetOffsetLimit()
		sqlRows := fmt.Sprintf("SELECT * FROM %s %s ORDER BY id DESC LIMIT %d, %d", r.DB.Driver.QuoteIdentifier(r.Name), whereStr, offset, limit)
		rows, err := r.DB.GetRows(sqlRows, params...)
		if err != nil {
			return nil, err
		}

		for _, row := range rows {
			list = append(list, row)
		}
	}
	return &ListOutput{
		Pagination: input.Pagination,
		List:       list,
	}, nil
}

func (r Resource) Delete(filter interface{}) error {
	f, ok := filter.(Filter)
	if !ok {
		f = Filter{
			"id": filter,
		}
	}

	where, params := f.GetWhere(r.DB)
	_, err := r.DB.Delete(r.Name, where, params...)
	return err
}
