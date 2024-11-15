package event

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/krixlion/dev_forum-lib/internal/testtypes"
)

func Test_MakeEvent(t *testing.T) {
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
				body:        "asJKDa",
				metadata:    map[string]string{"test": "asJKDa"},
			},
			want: Event{
				AggregateId: ArticleAggregate,
				Type:        ArticleDeleted,
				Body: func() []byte {
					data, err := json.Marshal("asJKDa")
					if err != nil {
						panic(err)
					}
					return data
				}(),
				Metadata:  map[string]string{"test": "asJKDa"},
				Timestamp: time.Now(),
			},
		},
		{
			name: "Test correctly serializes ArticleUpdated event with random data",
			args: args{
				aggregateId: ArticleAggregate,
				eType:       ArticleUpdated,
				body: testtypes.Article{
					Id:        "test-id",
					UserId:    "test-user-id",
					Title:     "test-title",
					Body:      "test-body",
					CreatedAt: time.Date(2000, time.April, 1, 1, 1, 1, 1, time.Local),
					UpdatedAt: time.Date(2000, time.April, 1, 1, 1, 1, 1, time.Local),
				},
			},
			want: Event{
				AggregateId: ArticleAggregate,
				Type:        ArticleUpdated,
				Body:        []byte(`{"id":"test-id","user_id":"test-user-id","title":"test-title","body":"test-body","created_at":"2000-04-01T01:01:01.000000001+02:00","updated_at":"2000-04-01T01:01:01.000000001+02:00"}`),
				Timestamp:   time.Now(),
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
