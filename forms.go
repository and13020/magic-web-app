package main

import "net/url"

type formErrors map[string][]string

// Form struct to hold form data and validation errors
type Form struct {
	url.Values
	Errors formErrors
}

func (e formErrors) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e formErrors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}

func NewForm(form url.Values) *Form {
	return &Form{
		form,
		formErrors(map[string][]string{}),
	}
}

func (f Form) Valid() bool {
	return len(f.Errors) == 0
}
