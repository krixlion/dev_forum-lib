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

type Operator string

const (
	Unknown            Operator = "UNKNOWN"
	Equal              Operator = "eq"
	NotEqual           Operator = "neq"
	GreaterThan        Operator = "gt"
	LesserThan         Operator = "lt"
	GreaterThanOrEqual Operator = "gte"
	LesserThanOrEqual  Operator = "lte"
)

type Parameter struct {
	Attribute string
	Operator  Operator
	Value     string
}

// AllOperators returns all registered operators' string representations keyed by their enum.
func AllOperators() map[Operator]string {
	return map[Operator]string{
		Equal:              string(Equal),
		Unknown:            string(Unknown),
		NotEqual:           string(NotEqual),
		GreaterThan:        string(GreaterThan),
		LesserThan:         string(LesserThan),
		GreaterThanOrEqual: string(GreaterThanOrEqual),
		LesserThanOrEqual:  string(LesserThanOrEqual),
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

		operator, err := MatchOperator(rawOperator)
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

func (param Parameter) ToFilter() string {

	query := ""
	query += strings.ToLower(param.Attribute)
	query += operatorOpeningSign
	query += string(param.Operator)
	query += operatorClosingSign
	query += valueAssigmentSign
	query += param.Value
	return query
}

// MatchOperator checks if provided input is a registered operator.
// Returns a non-nil error if the operator is not found.
func MatchOperator(input string) (Operator, error) {
	trimmed := strings.Trim(input, operatorPrefix)
	operator := Operator(trimmed)

	_, ok := AllOperators()[operator]
	if !ok {
		return "", ErrInvalidOperator
	}

	return operator, nil
}
