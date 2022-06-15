package pkg

import (
	"testing"
	"time"
)

func fakeNow() time.Time {
	return time.Date(2019, time.January, 1, 1, 23, 1, 0, time.UTC)
}

func TestGet(t *testing.T) {
	d, _ := time.ParseDuration("5m")
	expected := fakeNow().Truncate(d)
	if expected.Second() != 0 {
		t.Fatal("Expected seconds to be zero")
	}
}
