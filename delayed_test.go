package errorutil

import (
	"github.com/pkg/errors"
	"testing"
	"time"
)

func TestDelay(t *testing.T) {

	//	err := errorutil.NewRetryableError("Delayed Error Test")
	tests := []struct {
		err  error
		want time.Duration
	}{
		{nil, 0},
		{WithDelay(nil, 0), 0},
		{WithDelay(nil, 10), 0},
		{WithDelay(errors.New("foo"), 0), 0},
		{WithDelay(errors.New("foo"), 10), 10},
	}
	err := WithDelay(errors.New("foobar"), 10)
	if err.Error() != "foobar" {
		t.Errorf("WithDelay(%q): got: %v, want %v", err, err.Error(), "foobar")
	}
	for _, tt := range tests {
		got := Delay(tt.err)
		if got != tt.want {
			t.Errorf("WithDelay(%q): got: %v, want %v", tt.err, got, tt.want)
		}
	}
}

func TestNewDelayedError(t *testing.T) {
	err := NewDelayedError("delayed Error", 10)
	if err.Error() != "delayed Error" {
		t.Errorf("NewDelayedError expected %s, got %s", "delayed Error", err.Error())
	}
}
