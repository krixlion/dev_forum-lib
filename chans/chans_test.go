package chans

import (
	"testing"
	"time"

	"github.com/krixlion/dev_forum-lib/internal/gentest"
	"github.com/krixlion/dev_forum-lib/internal/testtypes"
	"github.com/stretchr/testify/assert"
)

func Test_FanIn(t *testing.T) {
	t.Run("Test returned channel receives all messages from multiple channels", func(t *testing.T) {
		want := []testtypes.Article{
			gentest.RandomArticle(5, 10),
			gentest.RandomArticle(6, 11),
			gentest.RandomArticle(7, 12),
		}

		chans := func() []<-chan testtypes.Article {
			chans := make([]<-chan testtypes.Article, 0, len(want))
			for _, v := range want {
				c := make(chan testtypes.Article, 1)
				c <- v
				chans = append(chans, c)
			}
			return chans
		}()

		out := FanIn(chans...)

		var got []testtypes.Article
		for i := 0; i < len(want); i++ {
			got = append(got, <-out)
		}

		if !assert.ElementsMatch(t, got, want) {
			t.Errorf("FanIn(): messages are not equal:\n got = %+v\n want = %+v\n", got, want)
		}
	})

	t.Run("Test returned chan is closed when origin chans are closed", func(t *testing.T) {
		chans := []chan struct{}{
			make(chan struct{}),
			make(chan struct{}),
		}

		out := FanIn(chans[0], chans[1])

		for _, c := range chans {
			close(c)
		}

		if _, ok := <-out; ok {
			t.Errorf("FanIn(): merged chan is open")
		}
	})
}

func TestNonBlockSend(t *testing.T) {
	t.Run("Test does not block when the channel is full", func(t *testing.T) {
		c := make(chan struct{}, 1)
		c <- struct{}{}

		done := make(chan struct{})

		go func() {
			NonBlockSend(c, struct{}{})
			done <- struct{}{}
		}()

		select {
		case <-time.After(time.Millisecond * 10):
			t.Error("NonBlockSend(): timed out waiting for the func to return (probably blocking)")
		case <-done:
			return
		}
	})

	t.Run("Test message is sent when the channel is not full", func(t *testing.T) {
		c := make(chan string, 1)

		want := "test"
		go func() {
			NonBlockSend(c, want)
		}()

		select {
		case <-time.After(time.Millisecond * 10):
			t.Error("NonBlockSend(): timed out waiting for the func to return (probably blocking)")
		case got := <-c:
			if got != want {
				t.Errorf("NonBlockSend():\n got = %v\n want = %v", got, want)
			}
		}
	})
}
