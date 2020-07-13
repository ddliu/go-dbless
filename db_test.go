package dbless

import (
	"testing"
)

func TestDB(t *testing.T) {
	db := newDB()
	columns, err := db.ListColumns("", "test")
	if err != nil {
		t.Error(err)
	}

	if columns[1].Name() != "name" {
		t.Error("list columns error")
	}
}
