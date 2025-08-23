package main

import (
	"net/url"
)

type errors map[string][]string

func (e errors) Get(field string) string {
	var errorsSliece []string = e[field]

	if len(errorsSliece) == 0 {
		return ""
	}
	return errorsSliece[0]
}

func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

type Form struct {
	Data   url.Values
	Errors errors
}

func NewForm(data url.Values) *Form {
	return &Form{

		Data:   data,
		Errors: map[string][]string{},
	}
}

func (f *Form) Has(field string) bool {
	x := f.Data.Get(field)
	if x == "" {
		return false
	}
	return true
}
