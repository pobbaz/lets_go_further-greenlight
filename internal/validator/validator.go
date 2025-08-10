package validator

import "regexp"

// regex for checking email format
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\\\/?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	Errors map[string]string
}

// New() function create a validator with empty errors map
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid() function return true if there are no errors
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError() function adds an error for a key if it doesn’t already exist
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check() function adds an error if the validation check is NOT ok
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// returns true if value is in the given list
func In(value string, list ...string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// returns true if value matches the regex.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// returns true if all values in the slice are unique
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)
	for _, v := range values {
		uniqueValues[v] = true
	}
	return len(values) == len(uniqueValues)
}
