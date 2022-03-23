package clock

import (
	"time"

	"github.com/stretchr/testify/mock"
)

var (
	Now func() time.Time
	On  func(methodName string, arguments ...interface{}) *mock.Call
)

// Когда я пытался сделать
func init() {
	// TODO: fix it
	//if test.IsTest() {
	//	clock = NewMock()
	//} else {
	//	clock = New()
	//}
	clock := New()
	Now = clock.Now
	On = clock.On
}

// Clock is an interface for working with current time.
// There is a problem: other libraries are using real time package,
// so our mock time does not work for any cases that need to count time interval that relates to the current time.
type Clock interface {
	Now() time.Time
	On(methodName string, arguments ...interface{}) *mock.Call
}
