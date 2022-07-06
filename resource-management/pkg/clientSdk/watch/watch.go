package watch

import (
	"global-resource-service/resource-management/pkg/common-lib/types/event"
)

// Interface can be implemented by anything that knows how to watch and report changes.
type Interface interface {
	// Stop stops watching. Will close the channel returned by ResultChan(). Releases
	// any resources used by the watch.
	Stop()

	// ResultChan returns a chan which will receive all the events. If an error occurs
	// or Stop() is called, the implementation will close this channel and
	// release any resources used by the watch.
	ResultChan() <-chan event.NodeEvent
}

type emptyWatch chan event.NodeEvent

// NewEmptyWatch returns a watch interface that returns no results and is closed.
// May be used in certain error conditions where no information is available but
// an error is not warranted.
func NewEmptyWatch() Interface {
	ch := make(chan event.NodeEvent)
	close(ch)
	return emptyWatch(ch)
}

// Stop implements Interface
func (w emptyWatch) Stop() {
}

// ResultChan implements Interface
func (w emptyWatch) ResultChan() <-chan event.NodeEvent {
	return chan event.NodeEvent(w)
}
