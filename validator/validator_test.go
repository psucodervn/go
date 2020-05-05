package validator

import (
	"encoding/json"
	"os"
	"regexp"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestStructValidator_Validate(t *testing.T) {
	type Info struct {
		Email string `json:"email" validate:"required"`
	}
	type request struct {
		Name     string `json:"name" validate:"required"`
		Age      int    `json:"age" validate:"number,gt=5"`
		Info     Info   `json:"info"`
		Duration string `json:"duration" validate:"duration,required"`
	}
	type args struct {
		Req request `json:"req"`
	}
	tests := []struct {
		name      string
		args      args
		numErrors int
	}{
		{
			name: "validation should pass",
			args: args{
				Req: request{
					Name:     "psucodervn",
					Age:      15,
					Info:     Info{Email: "psucodervn@example.com"},
					Duration: "1d",
				},
			},
			numErrors: 0,
		},
		{
			name: "validation should fail",
			args: args{
				Req: request{
					Age:      3,
					Duration: "2d);--",
				},
			},
			numErrors: 4,
		},
	}

	v := NewStructValidator()
	v.Register(func(vl *validator.Validate) error {
		reDuration := regexp.MustCompile(`^[1-9]+[0-9]*([smhdwy]|mo)$`)
		return vl.RegisterValidation("duration", func(fl validator.FieldLevel) bool {
			return reDuration.MatchString(fl.Field().String())
		})
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.args.Req)
			if err == nil {
				if tt.numErrors == 0 {
					return
				}
				t.Errorf("Validate() return nil, want %v errors", tt.numErrors)
			}

			errs := err.(Errors)
			if len(errs) > 0 {
				_ = json.NewEncoder(os.Stderr).Encode(errs)
			}
			if len(errs) != tt.numErrors {
				t.Errorf("Validate() return %v errors, want %v errors", len(errs), tt.numErrors)
			}
		})
	}
}
