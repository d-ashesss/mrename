package observer_test

import (
	"errors"
	mocks "github.com/d-ashesss/mrename/mocks"
	"github.com/d-ashesss/mrename/observer"
	"testing"
)

func TestObserver(t *testing.T) {
	err := errors.New("test error")
	o := observer.New()
	s := mocks.NewSubscriber(t)
	s.On("Notify", observer.Event{Name: "1st event", File: "1st.txt", Result: "first.txt"}).Once()
	s.On("Notify", observer.Event{Name: "2nd event", File: "2nd.txt", Error: err}).Once()
	o.AddSubscriber(s)
	o.PublishResult("1st event", "1st.txt", "first.txt")
	o.PublishError("2nd event", "2nd.txt", err)
}
