package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// StringSlice - a postgres compatible string slice type
type StringSlice []string

// sourced from https://gist.github.com/adharris/4163702
func (s *StringSlice) Scan(src interface{}) error {
	srcString := ""
	switch src := src.(type) {
	case string:
		srcString = src
	case []byte:
		srcString = string(src)
	default:
		return fmt.Errorf("unsupported type %v", src)
	}

	parsed := parseStringArray(srcString)
	(*s) = StringSlice(parsed)
	return nil
}

// Parse the output string from the array type.
// Regex used: (((?P<value>(([^",\\{}\s(NULL)])+|"([^"\\]|\\"|\\\\)*")))(,)?)
func parseStringArray(array string) []string {
	results := make([]string, 0)
	matches := arrayExp.FindAllStringSubmatch(array, -1)
	for _, match := range matches {
		s := match[valueIndex]
		// the string _might_ be wrapped in quotes, so trim them:
		s = strings.Trim(s, "\"")
		results = append(results, s)
	}
	return results
}

// Value returns the driver compatible value
func (s StringSlice) Value() (driver.Value, error) {
	var strs []string
	for _, i := range s {
		strs = append(strs, fmt.Sprintf(`"%s"`, i))
	}
	return "{" + strings.Join(strs, ",") + "}", nil
}
