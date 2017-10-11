package dusk

type Event struct {
	handlers []*func(data interface{})
}

func NewEvent() *Event {
	return &Event{
		handlers: []*func(data interface{}){},
	}
}

func (event *Event) Subscribe(fn *func(data interface{})) {
	event.handlers = append(event.handlers, fn)
}

func (event *Event) Unsubscribe(fn *func(data interface{})) {
	for i := range event.handlers {
		if event.handlers[i] == fn {
			event.handlers[i] = event.handlers[len(event.handlers)-1]
			event.handlers = event.handlers[:len(event.handlers)-1]
		}
	}
}

func (event *Event) Call(data interface{}) {
	for i := range event.handlers {
		(*event.handlers[i])(data)
	}
}
