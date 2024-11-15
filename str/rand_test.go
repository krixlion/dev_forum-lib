package str

import "testing"

func Test_RandomAlphaString(t *testing.T) {
	t.Run("Test if returned string has correct length", func(t *testing.T) {
		length := 5
		got, err := RandomAlphaString(length)
		if err != nil {
			t.Errorf("RandomAlphaString() error = %v", err)
			return
		}

		if len(got) != length {
			t.Errorf("RandomAlphaString() invalid length: got = %v expected length %v", got, length)
		}
	})

	t.Run("Test if two equal subsequent calls return different results", func(t *testing.T) {
		length := 5

		first, err := RandomAlphaString(length)
		if err != nil {
			t.Errorf("RandomAlphaString() error = %v", err)
			return
		}

		second, err := RandomAlphaString(length)
		if err != nil {
			t.Errorf("RandomAlphaString() error = %v", err)
			return
		}

		if first == second {
			t.Errorf("RandomAlphaString(): returned the same string twice:\n first = %v\n second = %v\n", first, second)
		}
	})
}
