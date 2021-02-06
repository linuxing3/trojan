package core

import (
	"reflect"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSqliteUser_Init(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(defaultPath), &gorm.Config{})
	var user SqliteUser
	tests := []struct {
		name string
		user *SqliteUser
		want *gorm.DB
	}{
		// TODO: Add test cases.
		{
			name: "Test sqlite Gorm method",
			user: &user,
			want:  db,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.user.Init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SqliteUser.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSqliteUser_Crud(t *testing.T) {
	query := make(map[string]interface{})
	query["UserName"] = "xray"
	var user SqliteUser
	type args struct {
		flag  string
		query map[string]interface{}
	}
	tests := []struct {
		name string
		user *SqliteUser
		args args
		want *SqliteUser
	}{
		// TODO: Add test cases.
		{
			name: "test crud for sqlite user",
			user: &user,
			args: args{
				flag: "create",
				query: query,
			},
			want: &user,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.user.Crud(tt.args.flag, tt.args.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SqliteUser.Crud() = %v, want %v", got, tt.want)
			}
		})
	}
}
