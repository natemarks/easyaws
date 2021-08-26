package internal

import "testing"

func TestDdd(t *testing.T) {
	type args struct {
		is string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "valid", args: args{is: "dfgdgd"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
