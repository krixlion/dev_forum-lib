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
		want    Filter
		wantErr bool
	}{
		{
			name: "Test if parses simple alphanumeric query with an underscore and multiple params",
			args: args{query: "name1[$eq]=john&last_na2me[$eq]=doe"},
			want: Filter{
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
			want: Filter{
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
			name:    "Test if fails on filter ending with parameter separator",
			args:    args{query: "na-me_[$eq]=john&last_name-[$eq]=doe&"},
			wantErr: true,
		},
		{
			name:    "Test if fails on dot in param attribute",
			args:    args{query: "na.me[$eq]=john&last.name[$eq]=doe"},
			wantErr: true,
		},
		{
			name: "Test if returns nil on empty query",
			args: args{query: ""},
			want: nil,
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

func Test_MatchOperator(t *testing.T) {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MatchOperator(tt.args.operator)
			if (err != nil) != tt.wantErr {
				t.Errorf("MatchOperator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want && !tt.wantErr {
				t.Errorf("MatchOperator() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestFilter_String(t *testing.T) {
	tests := []struct {
		name   string
		filter Filter
		want   string
	}{
		{
			name: "Test if parameter suffix is cut correctly",
			filter: Filter{
				{
					Attribute: "user_id",
					Operator:  Equal,
					Value:     "1",
				},
				{
					Attribute: "name",
					Operator:  Equal,
					Value:     "5",
				},
			},
			want: "user_id[$eq]=1&name[$eq]=5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filter.String(); got != tt.want {
				t.Errorf("Filter.String():\n got = %v\n want %v", got, tt.want)
			}
		})
	}
}

func TestParameter_String(t *testing.T) {
	tests := []struct {
		name  string
		param Parameter
		want  string
	}{
		{
			name: "Test",
			param: Parameter{
				Attribute: "user_id",
				Operator:  Equal,
				Value:     "1",
			},
			want: "user_id[$eq]=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.param.String(); got != tt.want {
				t.Errorf("Parameter.String():\n got = %v\n want = %v", got, tt.want)
			}
		})
	}
}
