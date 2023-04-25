package str

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_ToLowerSnakeCase(t *testing.T) {
	tests := []struct {
		desc string
		arg  string
		want string
	}{
		{
			desc: "Test on simple data",
			arg:  "UserId",
			want: "user_id",
		},
		{
			desc: "Test on simple data",
			arg:  "UserMightWanttoFixThat",
			want: "user_might_wantto_fix_that",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := ToLowerSnakeCase(tt.arg)

			if !cmp.Equal(got, tt.want) {
				t.Errorf("Wrong output:\n got = %+v\n want = %+v\n", got, tt.want)
				return
			}
		})
	}
}
