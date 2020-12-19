package core

import "testing"

func TestWriteInbloudClient(t *testing.T) {
	type args struct {
		ids  []string
		flag string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WriteInbloudClient(tt.args.ids, tt.args.flag); got != tt.want {
				t.Errorf("WriteInbloudClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
