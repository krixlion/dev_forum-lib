package filter

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var ErrInvalidParam error = errors.New("invalid param")
var ErrValueNotFound error = errors.New("field's value not found")
var ErrInvalidOperator error = errors.New("field's operator invalid")

const (
	parameterSeparator  = "&"
	valueAssigmentSign  = "="
	operatorOpeningSign = "["
	operatorClosingSign = "]"
	operatorPrefix      = "$"
)

type Operator int

const (
	Unknown Operator = 1 << iota
	Equal
	NotEqual
	GreaterThan
	LesserThan
	GreaterThanOrEqual
	LesserThanOrEqual
)

type Parameter struct {
	Attribute string
	Operator  Operator
	Value     string
}

func (param Parameter) ToFilter() string {

	query := ""
	query += strings.ToLower(param.Attribute)
	query += operatorOpeningSign
	query += param.Operator.String()
	query += operatorClosingSign
	query += valueAssigmentSign
	query += param.Value
	return query
}

func (s Operator) String() string {
	switch s {
	case Equal:
		return operatorPrefix + "eq"
	case NotEqual:
		return operatorPrefix + "neq"
	case GreaterThan:
		return operatorPrefix + "gt"
	case LesserThan:
		return operatorPrefix + "lt"
	case GreaterThanOrEqual:
		return operatorPrefix + "gte"
	case LesserThanOrEqual:
		return operatorPrefix + "lte"
	default:
		return "unknown operator"
	}
}

func matchOperator(operator string) (Operator, error) {
	switch operator {
	case Equal.String():
		return Equal, nil
	case NotEqual.String():
		return NotEqual, nil
	case GreaterThan.String():
		return GreaterThan, nil
	case LesserThan.String():
		return LesserThan, nil
	case GreaterThanOrEqual.String():
		return GreaterThanOrEqual, nil
	case LesserThanOrEqual.String():
		return LesserThanOrEqual, nil
	default:
		return Unknown, ErrInvalidOperator
	}
}

// Parse parses input query string.
//
// Example input:
//
//	params, err := filter.Parse("name[$eq]=john&last_name[$eq]=doe")
//	fmt.Printf("%+v", params)
//
// Output:
//
//	[{Field:name Operator:$eq Value:john} {Field:last_name Operator:$eq Value:doe}]%
func Parse(query string) ([]Parameter, error) {
	parsedParams := []Parameter{}
	params := strings.Split(query, parameterSeparator)

	for _, param := range params {
		beforeValue, value, found := strings.Cut(param, valueAssigmentSign)
		if !found {
			return nil, ErrValueNotFound
		}

		// Allow alphanumeric, lowercase names with underscore and dash followed by a prefixed operator within brackets.
		exp := fmt.Sprintf(`^[a-z0-9_-]+\%s[\%s[a-z]+\%s$`, operatorOpeningSign, operatorPrefix, operatorClosingSign)
		re, err := regexp.Compile(exp)
		if err != nil {
			return nil, err
		}

		if !re.MatchString(beforeValue) {
			return nil, ErrInvalidParam
		}

		parsed := strings.Split(beforeValue, operatorOpeningSign)
		attribute := parsed[0]
		rawOperator := parsed[1]
		rawOperator = strings.Trim(rawOperator, operatorClosingSign)

		operator, err := matchOperator(rawOperator)
		if err != nil {
			return nil, err
		}

		parsedParams = append(parsedParams, Parameter{
			Attribute: attribute,
			Operator:  operator,
			Value:     value,
		})

	}

	return parsedParams, nil
}
