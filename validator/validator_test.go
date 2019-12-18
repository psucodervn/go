package validator

import (
	"encoding/json"
	"os"
	"testing"
)

func TestStructValidator_Validate(t *testing.T) {
	type args struct {
		req interface{}
	}
	tests := []struct {
		name      string
		args      args
		numErrors int
	}{
		{
			name: "validation should failed",
			args: args{
				req: struct {
					Name string `json:"name" validate:"required"`
					Age  int    `json:"age" validate:"number,gt=5"`
				}{
					Age: 3,
				},
			},
			numErrors: 2,
		},
	}

	v := NewStructValidator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := v.Validate(tt.args.req)
			if len(errs) > 0 {
				_ = json.NewEncoder(os.Stderr).Encode(errs)
			}
			if len(errs) != tt.numErrors {
				t.Errorf("Validate() return %v errors, want %v errors", len(errs), tt.numErrors)
			}
		})
	}
}
