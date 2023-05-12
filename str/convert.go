package str

import "strconv"

func ConvertToUint(str string) (uint, error) {
	if str == "" {
		return 0, nil
	}

	num, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(num), nil
}
