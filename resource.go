// Copyright 2019 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package dbless

import (
	"fmt"
	"time"

	"database/sql"

	"github.com/spf13/cast"
)

type Record map[string]interface{}

type Field struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (r Record) ID() uint64 {
	id, ok := r["id"]
	if !ok {
		return 0
	}

	return cast.ToUint64(id)
}

type Resource struct {
	Name      string  `json:"name"`
	Title     string  `json:"title"`
	Fields    []Field `json:"fields"`
	Timestamp bool    `json:"timestamp"`
	DB        *sql.DB `json:"-"`
}

func (r Resource) Get(id uint64) (Record, error) {
	row, err := DBGetRow(r.DB, "select * from "+quote(r.Name)+" where id = ?", id)

	if err != nil {
		return nil, err
	}

	return Record(row), nil
}

func (r Resource) Save(record Record) (uint64, error) {
	id := record.ID()

	delete(record, "id")
	delete(record, "created_at")
	delete(record, "updated_at")
	// update
	if id > 0 {
		if r.Timestamp {
			record["updated_at"] = time.Now()
		}

		_, err := DBUpdate(r.DB, r.Name, record, "id = ?", id)
		return id, err
	} else {
		if r.Timestamp {
			record["created_at"] = time.Now()
			record["updated_at"] = time.Now()
		}

		id, err := DBInsert(r.DB, r.Name, record)

		return id, err
	}
}

type ListInput struct {
	Pagination Pagination `json: "pagination"`
	// Filter map[string]interface{}
}

type ListOutput struct {
	Pagination
	List []Record `json:"list"`
}

func (r Resource) List(input ListInput) (*ListOutput, error) {
	sqlTotal := "SELECT COUNT(*) FROM " + quote(r.Name)
	total, err := DBScalar(r.DB, sqlTotal)
	if err != nil {
		return nil, err
	}

	input.Pagination.SetTotal(uint(total))
	input.Pagination.Valid()

	var list []Record
	if total != 0 {
		offset, limit := input.Pagination.GetOffsetLimit()
		sqlRows := fmt.Sprintf("SELECT * FROM %s ORDER BY id DESC LIMIT %d, %d", quote(r.Name), offset, limit)
		rows, err := DBGetRows(r.DB, sqlRows)
		if err != nil {
			return nil, err
		}

		for _, row := range rows {
			list = append(list, Record(row))
		}
	}
	return &ListOutput{
		Pagination: input.Pagination,
		List:       list,
	}, nil
}

func (r Resource) Delete(id uint64) error {
	_, err := DBDelete(r.DB, r.Name, "id = ?", id)
	return err
}
