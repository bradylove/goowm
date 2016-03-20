package broadcaster

type Event int32

const (
	EventWorkspaceChanged Event = iota
)

var instance *Broadcaster

type EventCallback func()

type Broadcaster struct {
	listeners map[Event][]EventCallback
}

func init() {
	instance = &Broadcaster{make(map[Event][]EventCallback)}
}

func Listen(e Event, cb EventCallback) {
	instance.add(e, cb)
}

func Trigger(e Event) {
	instance.trigger(e)
}

func (b *Broadcaster) add(e Event, cb EventCallback) {
	b.listeners[e] = append(b.listeners[e], cb)
}

func (b *Broadcaster) trigger(e Event) {
	for _, cb := range b.listeners[e] {
		cb()
	}
}
