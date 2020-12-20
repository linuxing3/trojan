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
	tests := []struct {
		name   string
		fields fields
		want   *sql.DB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlite := &Sqlite{
				Enabled:  tt.fields.Enabled,
				Path:     tt.fields.Path,
				Password: tt.fields.Password,
				Table:    tt.fields.Table,
			}
			if got := sqlite.GetDB(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sqlite.GetDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSqlite_CreateDefaultTable(t *testing.T) {
	type fields struct {
		Enabled  bool
		Path     string
		Password string
		Table    string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlite := &Sqlite{
				Enabled:  tt.fields.Enabled,
				Path:     tt.fields.Path,
				Password: tt.fields.Password,
				Table:    tt.fields.Table,
			}
			if got := sqlite.CreateDefaultTable(); got != tt.want {
				t.Errorf("Sqlite.CreateDefaultTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSqlite_CreateTable(t *testing.T) {
	type fields struct {
		Enabled  bool
		Path     string
		Password string
		Table    string
	}
	type args struct {
		dbName string
		fields []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlite := &Sqlite{
				Enabled:  tt.fields.Enabled,
				Path:     tt.fields.Path,
				Password: tt.fields.Password,
				Table:    tt.fields.Table,
			}
			if got := sqlite.CreateTable(tt.args.dbName, tt.args.fields); got != tt.want {
				t.Errorf("Sqlite.CreateTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSqlite_CreateUser(t *testing.T) {
	type fields struct {
		Enabled  bool
		Path     string
		Password string
		Table    string
	}
	type args struct {
		id         string
		username   string
		base64Pass string
		originPass string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlite := &Sqlite{
				Enabled:  tt.fields.Enabled,
				Path:     tt.fields.Path,
				Password: tt.fields.Password,
				Table:    tt.fields.Table,
			}
			if err := sqlite.CreateUser(tt.args.id, tt.args.username, tt.args.base64Pass, tt.args.originPass); (err != nil) != tt.wantErr {
				t.Errorf("Sqlite.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
