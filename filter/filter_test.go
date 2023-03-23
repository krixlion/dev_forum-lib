package filter

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/krixlion/dev_forum-lib/internal/gentest"
)

func Test_Parse(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		args    args
		want    []Parameter
		wantErr bool
	}{
		{
			name: "Test if parses simple alphanumeric query with an underscore and multiple params",
			args: args{query: "name1[$eq]=john&last_na2me[$eq]=doe"},
			want: []Parameter{
				{
					Attribute: "name1",
					Operator:  Equal,
					Value:     "john",
				},
				{
					Attribute: "last_na2me",
					Operator:  Equal,
					Value:     "doe",
				},
			},
		},
		{
			name: "Test if allows for dashes and underscores",
			args: args{query: "na-me_[$eq]=john&last_name-[$eq]=doe"},
			want: []Parameter{
				{
					Attribute: "na-me_",
					Operator:  Equal,
					Value:     "john",
				},
				{
					Attribute: "last_name-",
					Operator:  Equal,
					Value:     "doe",
				},
			},
		},
		{
			name:    "Test if fails on dot in param attribute",
			args:    args{query: "na.me[$eq]=john&last.name[$eq]=doe"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) && !tt.wantErr {
				t.Errorf("Parse() got = %+v\n want %+v\n %v", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func Test_matchOperator(t *testing.T) {
	type args struct {
		operator string
	}
	tests := []struct {
		name    string
		args    args
		want    Operator
		wantErr bool
	}{
		{
			name: "Test if matches equal operator",
			args: args{operator: "$eq"},
			want: Equal,
		},
		{
			name: "Test if matches not equal operator",
			args: args{operator: "$neq"},
			want: NotEqual,
		},
		{
			name: "Test if matches greater than operator",
			args: args{operator: "$gt"},
			want: GreaterThan,
		},
		{
			name: "Test if matches lesser than operator",
			args: args{operator: "$lt"},
			want: LesserThan,
		},
		{
			name: "Test if matches greater than or equal operator",
			args: args{operator: "$gte"},
			want: GreaterThanOrEqual,
		},
		{
			name: "Test if matches lesser than or equal operator",
			args: args{operator: "$lte"},
			want: LesserThanOrEqual,
		},
		{
			name:    "Test if fails on random string",
			args:    args{operator: gentest.RandomString(5)},
			wantErr: true,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchOperator(tt.args.operator)
			if (err != nil) != tt.wantErr {
				t.Errorf("matchOperator() set:%d error = %v, wantErr %v", i, err, tt.wantErr)
				return
			}
			if got != tt.want && !tt.wantErr {
				t.Errorf("matchOperator() set:%d got = %+v, want %+v", i, got, tt.want)
			}
		})
	}
}
