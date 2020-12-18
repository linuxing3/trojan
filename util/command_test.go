package util

import "testing"

func TestCheckCommandExists(t *testing.T) {
	type args struct {
		command string
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
			if got := CheckCommandExists(tt.args.command); got != tt.want {
				t.Errorf("CheckCommandExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
