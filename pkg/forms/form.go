package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "this field cannot be blank")
		}

	}
}

func (f *Form) MaxLenth(field string, d int) {

	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This is too ling (maximun is %d characters)", d))
	}

}

func (f *Form) PermittiedValues(field string, opts ...string) {
	value := f.Get(field)

	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "this field is invalid")
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
