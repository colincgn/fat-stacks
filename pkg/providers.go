package pkg

import "time"

type TimeProvider interface {
	Now() time.Time
}

type LiveTimeProvider struct{}

func (LiveTimeProvider) Now() time.Time {
	return time.Now()
}
