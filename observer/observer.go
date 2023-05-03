package observer

type Event struct {
	Name   string
	File   string
	Error  error
	Result string
}

type Subscriber interface {
	Notify(Event)
}

type Observer struct {
	subscribers []Subscriber
}

func New() *Observer {
	return &Observer{subscribers: make([]Subscriber, 0)}
}

func (o *Observer) AddSubscriber(s Subscriber) {
	o.subscribers = append(o.subscribers, s)
}

func (o *Observer) PublishResult(event, file, result string) {
	o.Publish(Event{Name: event, File: file, Result: result})
}

func (o *Observer) PublishError(event, file string, err error) {
	o.Publish(Event{Name: event, File: file, Error: err})
}

func (o *Observer) Publish(event Event) {
	for _, s := range o.subscribers {
		s.Notify(event)
	}
}
