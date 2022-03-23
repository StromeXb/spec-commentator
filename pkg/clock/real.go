package clock

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// New returns new clock
func New() Clock {
	return realClock{}
}

type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}

func (realClock) On(methodName string, arguments ...interface{}) *mock.Call {
	return nil
}
