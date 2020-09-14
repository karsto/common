package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

// IntSlice - a postgres compatible int slice type
type IntSlice []int

// TODO: consider using pgx intSlice

// sourced from https://gist.github.com/adharris/4163702
func (s *IntSlice) Scan(src interface{}) error {
	srcString := ""
	switch src := src.(type) {
	case string:
		srcString = src
	case []byte:
		srcString = string(src)
	default:
		return fmt.Errorf("unsupported type %v", src)
	}

	parsed := parseIntArray(srcString)
	(*s) = IntSlice(parsed)
	return nil
}

func parseIntArray(array string) []int {
	results := make([]int, 0)
	matches := arrayExp.FindAllStringSubmatch(array, -1)
	for _, match := range matches {
		s := match[valueIndex]
		// the string _might_ be wrapped in quotes, so trim them:
		s = strings.Trim(s, "\"")
		asInt, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		results = append(results, asInt)
	}
	return results
}

// Value returns the driver compatible value
func (s IntSlice) Value() (driver.Value, error) {
	var strs []string
	for _, i := range s {
		strs = append(strs, strconv.Itoa(i))
	}
	return "{" + strings.Join(strs, ",") + "}", nil
}
