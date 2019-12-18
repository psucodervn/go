package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type StructValidator struct {
	validator *validator.Validate
}

type Error struct {
	Field string      `json:"field,omitempty"`
	Rule  string      `json:"rule,omitempty"`
	Value interface{} `json:"value,omitempty"`
	Text  string      `json:"text,omitempty"`
}

func NewStructValidator() *StructValidator {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &StructValidator{validator: v}
}

func (v *StructValidator) Validate(i interface{}) []Error {
	err := v.validator.Struct(i)
	if err == nil {
		return nil
	}

	valErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	var errs []Error
	for _, e := range valErrors {
		errs = append(errs, Error{
			Field: e.Namespace(),
			Rule:  e.Tag(),
			Value: e.Value(),
			Text:  fmt.Sprintf("%s failed on tag %s", e.Namespace(), e.Tag()),
		})
	}
	return errs
}
