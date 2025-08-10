package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Runtime still has int32 as its underlying type
type Runtime int32

// Error returned when runtime format is wrong

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

// MarshalJSON method makes Runtime implement json.Marshaler interface
func (r Runtime) MarshalJSON() ([]byte, error) {
	// create string like 170 mins
	jsonValue := fmt.Sprintf("%d mins", r)
	// wraps the string in double quotes
	quotedJSONValue := strconv.Quote(jsonValue)
	// return as []byte
	return []byte(quotedJSONValue), nil
}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	// remove surrounding double quotes
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Split into parts: number and unit
	parts := strings.Split(unquotedJSONValue, " ")

	// Check that format is exactly "<number> mins"
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	// Convert the first part (number) into int32
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Store the value in the Runtime type
	*r = Runtime(i)
	return nil
}
