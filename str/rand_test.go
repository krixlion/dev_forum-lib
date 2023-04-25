package str

import "testing"

func Test_RandomAlphaString(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test if returned string has correct length",
			args: args{
				length: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RandomAlphaString(tt.args.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("randomAlphaString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != tt.args.length {
				t.Errorf("randomAlphaString() invalid length: got = %v expected length %v", got, tt.args.length)
			}
		})
	}
}
