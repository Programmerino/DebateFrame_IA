package waiter

import (
	"time"

	"github.com/dennwc/dom"
)

// EventWaiter runs a function once an event stops triggering for the amount of milliseconds specified in tTime
func EventWaiter(node *dom.NodeBase, event string, tTime int, fn func()) {
	sleepTime := tTime / 3
	timeout := 0
	happened := false
	node.AddEventListener(event, func(e dom.Event) {
		happened = true
		timeout = 3
	})
	go func() {
		for {
			if timeout > 0 {
				timeout--
			} else if timeout <= 0 && happened {
				timeout = 3
				happened = false
				fn()
			}
			time.Sleep(time.Duration(int(time.Millisecond) * sleepTime))
		}
	}()
}
