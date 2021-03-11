package mammoth

import (
	"syscall/js"
	"time"

	"github.com/dennwc/dom"

	"gitlab.com/256/DebateFrame/client/log"
)

// ConvertToHTML generates HTML from the provided file and calls the function given with the HTML
func ConvertToHTML(file *js.Value, fn func(html string)) {
	FileAsArrayBuffer(file, func(buffer *js.Value) {
		options := make(map[string]interface{})
		options["arrayBuffer"] = *buffer

		mammoth := js.Global().Get("mammoth")

		cb := js.NewEventCallback(0, func(e js.Value) {
			fn(e.Get("value").String())
			stopLoadingScreen()
			log.DebugMessage("Done!")
		})
		log.DebugMessage("Passing off to mammoth")
		loadingScreen()
		promise := mammoth.Call("convertToHtml", options)
		promise.Call("then", cb)
		promise.Call("done")
	})
}

func loadingScreen() {
	modalElem := dom.GetDocument().GetElementById("modal-mammothload")
	modalElem.ClassList().Remove("simplehide")
	js.Global().Get("UIkit").Call("modal", dom.GetDocument().GetElementById("modal-mammothload").JSValue()).Call("show")
}

func stopLoadingScreen() {
	modalElem := dom.GetDocument().GetElementById("modal-mammothload")
	js.Global().Get("UIkit").Call("modal", modalElem.JSValue()).Call("hide")
	go func() {
		// Fixes bug where hiding happens at the same time that the showing happens, so after 300 milliseconds, it ensures that the dialog is truly gone
		time.Sleep(time.Millisecond * 300)
		modalElem.ClassList().Add("simplehide")
	}()
}

// FileAsArrayBuffer converts the file given to an array buffer
func FileAsArrayBuffer(file *js.Value, fn func(*js.Value)) {
	log.DebugMessage("Starting ArrayBuffer generation")
	reader := js.Global().Get("FileReader").New()

	cb := js.NewEventCallback(0, func(event js.Value) {
		arrayBuffer := event.Get("target").Get("result")
		log.DebugMessage("Finished ArrayBuffer generation")
		fn(&arrayBuffer)
	})

	reader.Set("onload", cb)

	reader.Call("readAsArrayBuffer", *file)
}
