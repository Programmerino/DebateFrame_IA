package log

import (
	"fmt"
	"time"

	"syscall/js"

	"github.com/dennwc/dom"
	"gitlab.com/256/WebFrame/waquery"
)

// Levels determines how many messages will be printed
var Levels = []string{"DEBUG", "WARN"}

// DebugMessage prints message if debugging messages are enabled
func DebugMessage(msg string, a ...interface{}) {
	if sliceExists(Levels, "DEBUG") {
		fmt.Printf("[FRAME DEBUG] %s\n", fmt.Sprintf(msg, a...))
	}
}

// WarnMessage prints message if warning messages are enabled
func WarnMessage(msg string, a ...interface{}) {
	if sliceExists(Levels, "WARN") {
		fmt.Printf("[FRAME WARN] %s\n", fmt.Sprintf(msg, a...))
	}
}

// PanicMessage prints message then panics
func PanicMessage(msg string, err error) {
	modal := dom.GetDocument().GetElementById("modal-panic")
	js.Global().Get("UIkit").Call("modal", modal.JSValue()).Call("show")
	waquery.AddStyle(".uk-modal * { opacity: 1; }")
	modal.AddEventListener("shown", func(e dom.Event) {
		fmt.Printf("[FRAME PANIC] %s\n", msg)
		panic(err)
	})
	time.Sleep(time.Second)
	fmt.Printf("[FRAME PANIC] %s\n", msg)
	panic(err)
}

func sliceExists(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
