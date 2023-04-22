package event

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/krixlion/dev_forum-lib/internal/gentest"
)

func Test_MakeEvent(t *testing.T) {
	randString := gentest.RandomString(5)
	randArticle := gentest.RandomArticle(1, 2)
	type args struct {
		aggregateId AggregateId
		eType       EventType
		data        interface{}
	}
	tests := []struct {
		name string
		args args
		want Event
	}{
		{
			name: "Test is correctly serializes ArticleDeleted event with random data",
			args: args{
				aggregateId: ArticleAggregate,
				eType:       ArticleDeleted,
				data:        randString,
			},
			want: Event{
				AggregateId: ArticleAggregate,
				Type:        ArticleDeleted,
				Body: func() []byte {
					data, err := json.Marshal(randString)
					if err != nil {
						panic(err)
					}
					return data
				}(),
				Timestamp: time.Now(),
			},
		},
		{
			name: "Test is correctly serializes ArticleUpdated event with random data",
			args: args{
				aggregateId: ArticleAggregate,
				eType:       ArticleUpdated,
				data:        randArticle,
			},
			want: Event{
				AggregateId: ArticleAggregate,
				Type:        ArticleUpdated,
				Body: func() []byte {
					data, err := json.Marshal(randArticle)
					if err != nil {
						panic(err)
					}
					return data
				}(),
				Timestamp: time.Now(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeEvent(tt.args.aggregateId, tt.args.eType, tt.args.data)
			if err != nil {
				t.Errorf("MakeEvent(): error = %v", err)
				return
			}

			if !cmp.Equal(got, tt.want, cmpopts.EquateApproxTime(time.Millisecond*5)) {
				t.Errorf("MakeEvent():\n got = %+v\n want = %+v\n diff = %+v\n", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}
