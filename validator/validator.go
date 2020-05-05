package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type StructValidator struct {
	validate *validator.Validate
}

type Error struct {
	Field string      `json:"field,omitempty"`
	Rule  string      `json:"rule,omitempty"`
	Value interface{} `json:"value,omitempty"`
	Text  string      `json:"text,omitempty"`
}

type Errors []Error

func (errs Errors) Error() string {
	var es []string
	for _, e := range errs {
		es = append(es, e.Text)
	}
	return strings.Join(es, ", ")
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
	return &StructValidator{validate: v}
}

func (v *StructValidator) Validate(i interface{}) error {
	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}

	valErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	var errs Errors
	typ := reflect.TypeOf(i)
	stripPrefix := len(typ.Name()) > 0
	for _, e := range valErrors {
		ns := e.Namespace()
		if stripPrefix {
			i := strings.IndexRune(ns, '.')
			if i >= 0 {
				ns = ns[i+1:]
			}
		}
		errs = append(errs, Error{
			Field: ns,
			Rule:  e.Tag(),
			Value: e.Value(),
			Text:  fmt.Sprintf("%s failed on tag %s", ns, e.Tag()),
		})
	}
	return errs
}

func (v *StructValidator) GetValidate() *validator.Validate {
	return v.validate
}

func (v *StructValidator) Register(fn func(vl *validator.Validate) error) error {
	return fn(v.validate)
}
