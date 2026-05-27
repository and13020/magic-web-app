package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

type formErrors map[string][]string

// Form struct to hold form data and validation errors
type Form struct {
	url.Values
	Errors formErrors
}

// Add(field, message) adds the given field/message as key/value to the form errors map
func (e formErrors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get(field) retrieves the field error from form errors map is present
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

// Required(fields...) accepts any number of strings.
// Each field with only white space characters or is empty will add an error to form
func (f *Form) Required(fields ...string) *Form {
	for _, field := range fields {

		v := f.Get(field)

		if strings.TrimSpace(v) == "" {
			f.Errors.Add(field, fmt.Sprintf("Field: %v is required\n", field))
		}
	}
	return f
}

// MinLength will add an error to the form if min length not met
func (f *Form) MinLength(field string, minLen int) *Form {
	v := strings.TrimSpace(f.Get(field))
	if len(v) < minLen {
		f.Errors.Add(field, fmt.Sprintf("Field: %v must be at least: %d characters long\n", field, minLen))
	}
	return f
}

// MaxLength will add an error to the form if max length not met
func (f *Form) MaxLength(field string, maxLen int) *Form {
	v := strings.TrimSpace(f.Get(field))
	if len(v) >= maxLen {
		f.Errors.Add(field, fmt.Sprintf("Field: %v is longer than given min length: %d\n", field, maxLen))
	}
	return f
}

func (f *Form) MatchPass(p1, p2 string) *Form {

	v1 := strings.TrimSpace(f.Get(p1))
	v2 := strings.TrimSpace(f.Get(p2))

	if v1 != v2 {
		f.Errors.Add(p1, "Passwords must match\n")
	}

	return f
}

func (f *Form) ValidEmail(field string) *Form {
	v := strings.TrimSpace(f.Get(field))
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(v) {
		f.Errors.Add(field, fmt.Sprintf("Field: %v must be a valid email address\n", field))
	}
	return f
}

func (f *Form) PlainText(fields ...string) *Form {
	for _, field := range fields {
		v := f.Get(field)
		for _, r := range v {
			if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
				f.Errors.Add(field, fmt.Sprintf("Field: %v must contain only valid plaintext characters\n", field))
				break
			}
		}
	}
	return f
}
