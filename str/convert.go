package str

import "strconv"

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
