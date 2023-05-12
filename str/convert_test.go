package str

import "testing"

func Test_convertToUint(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    uint
		wantErr bool
	}{
		{
			name: "Test if empty string returns 0",
			args: args{
				str: "",
			},
			want: 0,
		},
		{
			name: "Test if works on a simple int value",
			args: args{
				str: "53",
			},
			want: 53,
		},
		{
			name: "Test if fails on float values",
			args: args{
				str: "55.5",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToUint(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToUint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("convertToUint() = %v, want %v", got, tt.want)
			}
		})
	}
}
