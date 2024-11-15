package str

import "testing"

func Test_RandomAlphaString(t *testing.T) {
	t.Run("Test if returned string has correct length", func(t *testing.T) {
		length := 5
		got, err := RandomAlphaString(length)
		if err != nil {
			t.Errorf("randomAlphaString() error = %v", err)
			return
		}

		if len(got) != length {
			t.Errorf("randomAlphaString() invalid length: got = %v expected length %v", got, length)
		}
	})
}
