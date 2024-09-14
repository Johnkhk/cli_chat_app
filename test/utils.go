package test

import "time"

type MockTimeProvider struct {
	CurrentTime time.Time
}

func (mtp *MockTimeProvider) Now() time.Time {
	return mtp.CurrentTime
}

func (mtp *MockTimeProvider) Advance(d time.Duration) {
	mtp.CurrentTime = mtp.CurrentTime.Add(d)
}
