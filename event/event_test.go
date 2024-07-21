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
		body        interface{}
		metadata    map[string]string
	}
	tests := []struct {
		name string
		args args
		want Event
	}{
		{
			name: "Test correctly serializes ArticleDeleted event with random data",
			args: args{
				aggregateId: ArticleAggregate,
				eType:       ArticleDeleted,
				body:        randString,
				metadata:    map[string]string{"test": randString},
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
				Metadata:  map[string]string{"test": randString},
				Timestamp: time.Now(),
			},
		},
		{
			name: "Test correctly serializes ArticleUpdated event with random data",
			args: args{
				aggregateId: ArticleAggregate,
				eType:       ArticleUpdated,
				body:        randArticle,
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
			got, err := MakeEvent(tt.args.aggregateId, tt.args.eType, tt.args.body, tt.args.metadata)
			if err != nil {
				t.Errorf("MakeEvent(): error = %v", err)
				return
			}

			if !cmp.Equal(got, tt.want, cmpopts.EquateApproxTime(time.Millisecond)) {
				t.Errorf("MakeEvent():\n got = %+v\n want = %+v\n diff = %+v\n", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}
