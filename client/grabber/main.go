package grabber

import (
	"syscall/js"

	"github.com/PuerkitoBio/goquery"
	"github.com/dennwc/dom"
	"gitlab.com/256/DebateFrame/client/log"
)

func testPerms() {
	_, err := goquery.NewDocument("http://www.google.com")
	if err != nil {
		js.Global().Get("UIkit").Call("modal", dom.GetDocument().GetElementById("modal-cors").JSValue()).Call("show")
		log.DebugMessage("CORS not disabled!")
	}
}

func Grab() {
	//testPerms()
}
