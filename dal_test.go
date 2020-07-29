package dal

import (
	"reflect"
	"testing"
)

func Test_DBInfo(t *testing.T) {
	model := Model{
		DriverName:     "postgres",
		DataSourceName: "invalid",
	}
	println(model.DriverName)
}

func TestModel(t *testing.T) {
	model := Model{
		DriverName:     "mysql",
		DataSourceName: "test@tcp(localhost)/test",
	}
	defer func() {
		if err := model.SQL("drop table `user`;"); err != nil {
			t.Error(err)
		}
	}()

	if err := model.SQL("create table user (id int primary key, name varchar(64));"); err != nil {
		t.Error(err)
	}
	type T struct {
		ID   int    `field:"id"`
		Name string `field:"name"`
	}
	values := []T{
		{1, "a"},
		{2, "b"},
	}

	// version
	if info := model.DBInfo(); len(info) != 1 {
		t.Error("cannot get database version info")
	}

	// write
	if _, err := model.Update("user", values); err != nil {
		t.Error(err)
	}

	// read
	checkRead := func() {
		if err := model.Read("user", []string{"id", "name"}, "", T{}); err != nil {
			t.Error(err)
		}
		if len(model.Records) != len(values) {
			t.Error("length of query results and records are not the same")
		}
		for i := 0; i < len(values); i++ {
			record := model.Records[i].(T)
			if record.ID != values[i].ID || record.Name != values[i].Name {
				t.Error("query results differs from origin values")
			}
		}
	}
	checkRead()

	// re-read
	checkRead()
}

func Test_parseValue(t *testing.T) {
	type being struct {
		Name string `field:"name"`
	}
	type person struct {
		being
		Age int `field:"age"`
	}
	p := person{
		being: being{
			Name: "John",
		},
		Age: 12,
	}
	fields, query := parseValue(reflect.ValueOf(p), "staff", "Update")

	tFields := []string{"Name", "Age"}
	if len(fields) != len(tFields) || fields[0] != tFields[0] || fields[1] != tFields[1] {
		t.Errorf("output fields %+v is different from %+v", fields, tFields)
	}
	if query != "insert into staff(name,age) values(?,?) on duplicate key update name=?,age=?;" {
		t.Error("output query string is wrong")
	}
}
