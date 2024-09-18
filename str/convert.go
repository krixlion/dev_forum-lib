package str

import "strconv"

// ConvertToUint converts given string into an unsigned integer.
// Returns 0 and a nil error on an empty string.
// Returns an error when the given string is not an unsigned integer.
func ConvertToUint(str string) (uint, error) {
	if str == "" {
		return 0, nil
	}

	num, err := strconv.ParseUint(str, 10, 0)
	if err != nil {
		return 0, err
	}

	return uint(num), nil
}

// ConvertToInt converts given string into an integer.
// Returns 0 and a nil error on an empty string.
// Returns an error when the given string is not an integer.
func ConvertToInt(str string) (int, error) {
	if str == "" {
		return 0, nil
	}

	num, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		return 0, err
	}

	return int(num), nil
}
