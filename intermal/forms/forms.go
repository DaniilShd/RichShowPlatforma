package forms

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Form struct {
	url.Values
	Errors errors
}

// Valid return true if there are no errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New init a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Cheks for Required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This cannot be blank")
		}
	}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		f.Errors.Add(field, "This cannot be blank")
		return false
	}
	return true
}

// Check for string length
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := r.Form.Get(field)
	if len(x) < length {
		f.Errors.Add(field, "This field must be more")
		return false
	}
	return true
}

func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "This no email")
	}
}

func (f *Form) IsFloat(field string) {
	if govalidator.IsFloat(f.Get(field)) || govalidator.IsInt(f.Get(field)) {
		f.Errors.Add(field, "This no number")
	}
}
