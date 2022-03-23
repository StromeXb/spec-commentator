package clock

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// MockClock is a mock of Now() related functions
type MockClock struct {
	mock.Mock
}

func NewMock() Clock {
	return new(MockClock)
}

// Now returns current time
func (c *MockClock) Now() time.Time {
	args := c.Called()
	return args.Get(0).(time.Time)
}
