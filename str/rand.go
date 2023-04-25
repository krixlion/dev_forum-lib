package str

import "crypto/rand"

func RandomAlphaString(length int) (string, error) {
	const (
		letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 52 possibilities
		letterIdxBits = 6                                                      // 6 bits to represent 64 possibilities / indexes
		letterIdxMask = 1<<letterIdxBits - 1                                   // All 1-bits, as many as letterIdxBits
	)
	result := make([]byte, length)
	bufferSize := int(float64(length) * 1.3)

	for i, j, randomBytes := 0, 0, []byte{}; i < length; j++ {
		if j%bufferSize == 0 {
			randomBytes = make([]byte, bufferSize)
			if _, err := rand.Read(randomBytes); err != nil {
				return "", err
			}
		}
		if idx := int(randomBytes[j%length] & letterIdxMask); idx < len(letterBytes) {
			result[i] = letterBytes[idx]
			i++
		}
	}

	return string(result), nil
}
