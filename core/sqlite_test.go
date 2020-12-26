package core

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestSqlite_GetDB(t *testing.T) {
	type fields struct {
		Enabled  bool
		Path     string
		Password string
		Table    string
	}

	f := fields{
		Enabled:  true,
		Path:     "./xray.db",
		Password: "",
		Table:    "",
	}

	db, _ := sql.Open("sqlite3", f.Path)
	tests := []struct {
		name   string
		fields fields
		want   *sql.DB
	}{
		// TODO: Add test cases.
		{
			name:   "Test sqlite db get method",
			fields: f,
			want:   db,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlite := &Sqlite{
				Enabled:  tt.fields.Enabled,
				Path:     tt.fields.Path,
				Password: tt.fields.Password,
				Table:    tt.fields.Table,
			}
			if got := sqlite.GetDB(); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("Sqlite.GetDB() = %v, want %v", reflect.TypeOf(got).String(), reflect.TypeOf(got).String())
			} else {
				t.Logf("Sqlite.GetDB() = %v, want %v", reflect.TypeOf(got).String(), reflect.TypeOf(got).String())
			}
		})
	}
}
