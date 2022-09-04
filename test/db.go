package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ddliu/go-dbless"
	"github.com/spf13/cast"
)

type setupDB func() *dbless.DB

func doTestDB(t *testing.T, setup setupDB) {
	db := setup()
	// tables, err := db.ListTables("")
	// if err != nil {
	// 	return err
	// }

	// hasTable := false
	// for _, v := range tables {
	// 	if v == "test" {
	// 		hasTable = true
	// 	}
	// }

	// if !hasTable {
	// 	return errors.New("table not exist")
	// }

	columns, err := db.ListColumns("", "test")
	if err != nil {
		t.Fatal("list columns error", err)
	}

	if columns[1].Name() != "name" {
		t.Fatal("list columns error: invalid name")
	}

	doTestUtils(t, db)

	doTestNull(t, db)

	doTestResource(t, setup)
}

func doTestUtils(t *testing.T, db *dbless.DB) {
	for i := 1; i <= 100; i++ {
		id, err := db.InsertGetID("test", map[string]interface{}{
			"name":       fmt.Sprintf("record%d", i),
			"created_at": time.Now(),
			"updated_at": time.Now(),
		})

		if err != nil || cast.ToInt(id) != i {
			t.Fatal("insert error", err.Error(), i)
		}
	}

	result, err := db.GetPaging(10, 2, "select * from test order by id asc")
	if err != nil || result.Pagination.PageTotal != 10 || len(result.List) != 10 || cast.ToString(result.List[0]["name"]) != "record11" {
		t.Fatal("paging error", err)
	}

	// unmarshal
	type Test struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	}

	m := Test{}
	err = result.List[0].Unmarshal(&m)

	if err != nil || m.Name != "record11" {
		t.Fatal("unmarshal error ", err)
	}
}

func doTestNull(t *testing.T, db *dbless.DB) {
	row, err := db.GetRow("select 1 as num, null as name")
	if err != nil || row["name"] != nil {
		t.Fatal("nullable error", err, row["name"])
	}

	id, err := db.InsertGetID("test", map[string]interface{}{
		"name":       "null",
		"updated_at": nil,
	})

	if err != nil {
		t.Fatal(err)
	}

	row, err = db.GetRow("select id, updated_at from test where name = 'null'")
	if err != nil {
		t.Fatal(err)
	}

	if row.ID() != id {
		t.Fatal()
	}
}
