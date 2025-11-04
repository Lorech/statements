package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"statements/pkg/ctime"
	"strings"
	"time"
)

// A filter instance that is implemented as an individual filter in the config.
type Filter interface {
	FieldName() string
	// If the condition matches the requirements.
	Match(value any) bool
}

// An abstract filter used for Go type conformity reasons.
type RawFilter struct {
	// A marshaled JSON field that must be decoded into Filter
	Raw json.RawMessage
}

// Unmarshals JSON data into the raw filter value.
func (r *RawFilter) UnmarshalJSON(data []byte) error {
	r.Raw = append([]byte(nil), data...)
	return nil
}

// Decodes the raw filter into a typesafe filter based on a field map.
func (r RawFilter) DecodeWithFieldMap(fields FieldMap) (Filter, error) {
	var meta struct {
		Field string `json:"field"`
	}
	if err := json.Unmarshal(r.Raw, &meta); err != nil {
		return nil, err
	}

	ftype := fields[meta.Field]

	switch ftype {
	case FieldTypeDate:
		var f DateFilter
		if err := json.Unmarshal(r.Raw, &f); err != nil {
			return nil, err
		}
		return f, nil
	case FieldTypeNumber:
		var f NumberFilter
		if err := json.Unmarshal(r.Raw, &f); err != nil {
			return nil, err
		}
		return f, nil
	case FieldTypeString:
		var f StringFilter
		if err := json.Unmarshal(r.Raw, &f); err != nil {
			return nil, err
		}
		return f, nil
	default:
		return nil, fmt.Errorf("unknown field or type for field %q", meta.Field)
	}
}

type FieldType int

const (
	FieldTypeUnknown FieldType = iota
	FieldTypeDate
	FieldTypeNumber
	FieldTypeString
)

type FieldMap map[string]FieldType

type DateCondition string

const (
	DateLessThan         DateCondition = "LESS_THAN"
	DateLessThanEqual    DateCondition = "LESS_THAN_EQUAL"
	DateGreaterThan      DateCondition = "GREATER_THAN"
	DateGreaterThanEqual DateCondition = "GREATER_THAN_EQUAL"
	DateEqual            DateCondition = "EQUAL"
	DateNotEqual         DateCondition = "NOT_EQUAL"
)

type NumberCondition string

const (
	NumberLessThan         NumberCondition = "LESS_THAN"
	NumberLessThanEqual    NumberCondition = "LESS_THAN_EQUAL"
	NumberGreaterThan      NumberCondition = "GREATER_THAN"
	NumberGreaterThanEqual NumberCondition = "GREATER_THAN_EQUAL"
	NumberEqual            NumberCondition = "EQUAL"
	NumberNotEqual         NumberCondition = "NOT_EQUAL"
)

type StringCondition string

const (
	StringEqual      StringCondition = "EQUAL"
	StringNotEqual   StringCondition = "NOT_EQUAL"
	StringContain    StringCondition = "CONTAIN"
	StringNotContain StringCondition = "NOT_CONTAIN"
)

// A filter applied to a date.
type DateFilter struct {
	Field      string        `json:"field"`
	Condition  DateCondition `json:"condition"`
	Comparison string        `json:"comparison"`
}

// A filter applied to a number.
type NumberFilter struct {
	Field      string          `json:"field"`
	Condition  NumberCondition `json:"condition"`
	Comparison int             `json:"comparison"`
}

// A filter applied to a string.
type StringFilter struct {
	Field      string          `json:"field"`
	Condition  StringCondition `json:"condition"`
	Comparison string          `json:"comparison"`
}

// Returns the name of the field in the filter when dealing with the filter interface.
func (f DateFilter) FieldName() string {
	return f.Field
}

// Checks if a date filter matches a given value.
func (f DateFilter) Match(value any) bool {
	v, ok := value.(string)
	if !ok {
		return false
	}

	t, err := time.Parse(ctime.LittleEndianDateOnly, v)
	if err != nil {
		return false
	}

	c, err := time.Parse(ctime.LittleEndianDateOnly, f.Comparison)
	if err != nil {
		return false
	}

	diff := c.Compare(t)

	switch f.Condition {
	case DateLessThan:
		return diff == -1
	case DateLessThanEqual:
		return diff == -1 || diff == 0
	case DateGreaterThan:
		return diff == 1
	case DateGreaterThanEqual:
		return diff == 1 || diff == 0
	case DateEqual:
		return diff == 0
	case DateNotEqual:
		return diff != 0
	}

	return false
}

// Returns the name of the field in the filter when dealing with the filter interface.
func (f NumberFilter) FieldName() string {
	return f.Field
}

// Checks if a number filter matches a given value.
func (f NumberFilter) Match(value any) bool {
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return false
	}

	var i int
	switch v.Kind() {
	case reflect.Int:
		i = int(v.Int())
	default:
		// Attempt to convert things like int enums.
		if v.Type().ConvertibleTo(reflect.TypeOf("")) {
			i = int(v.Convert(reflect.TypeOf("")).Int())
		} else {
			return false
		}
	}

	switch f.Condition {
	case NumberLessThan:
		return i < f.Comparison
	case NumberLessThanEqual:
		return i <= f.Comparison
	case NumberGreaterThan:
		return i > f.Comparison
	case NumberGreaterThanEqual:
		return i >= f.Comparison
	case NumberEqual:
		return i == f.Comparison
	case NumberNotEqual:
		return i != f.Comparison
	}

	return false
}

// Returns the name of the field in the filter when dealing with the filter interface.
func (f StringFilter) FieldName() string {
	return f.Field
}

// Checks if a string filter matches a given value.
func (f StringFilter) Match(value any) bool {
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return false
	}

	var s string
	switch v.Kind() {
	case reflect.String:
		s = v.String()
	default:
		// Attempt to convert things like string enums.
		if v.Type().ConvertibleTo(reflect.TypeOf("")) {
			s = v.Convert(reflect.TypeOf("")).String()
		} else {
			return false
		}
	}

	switch f.Condition {
	case StringEqual:
		return s == f.Comparison
	case StringNotEqual:
		return s != f.Comparison
	case StringContain:
		return strings.Contains(s, f.Comparison)
	case StringNotContain:
		return !strings.Contains(s, f.Comparison)
	}

	return false
}
