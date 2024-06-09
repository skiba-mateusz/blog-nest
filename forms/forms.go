package forms

import (
	"fmt"
	"net/mail"
	"net/url"
	"strings"
)

type Form struct {
	Values url.Values
	Errors errors
}

func New(values url.Values) *Form {
	return &Form{
		Values: values,
		Errors: errors{},
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Values.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field is required")
		}
	}
}

func (f *Form) MinLength(field string, length int) {
	value := f.Values.Get(field)
	if len(value) < length {
		f.Errors.Add(field, fmt.Sprintf("This field needs to be at least %d characters long", length))
	}
}

func (f *Form) Email(field string) {
	value := f.Values.Get(field)
	_, err := mail.ParseAddress(value)
	if err != nil {
		f.Errors.Add(field, "Invalid email address")
	}
}

func (f *Form) PasswordsMatch(field, repeatedField string) {
	value := f.Values.Get(field)
	repeatedValue := f.Values.Get(repeatedField)
	if value != repeatedValue {
		f.Errors.Add(field, "Passwords do not match")
		f.Errors.Add(repeatedField, "Passwords do no match")
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}