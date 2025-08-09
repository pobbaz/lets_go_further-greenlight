package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

// MarshalJSON method makes Runtime implement json.Marshaler interface
func (r Runtime) MarshalJSON() ([]byte, error) {
	// create string like 170 mins
	jsonValue := fmt.Sprintf("%d mins", r)
	// wraps the string in double quotes
	quotedJSONValue := strconv.Quote(jsonValue)
	// return as []byte
	return []byte(quotedJSONValue), nil
}
