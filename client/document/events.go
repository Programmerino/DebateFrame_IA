package document

import (
	"path/filepath"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	"gitlab.com/256/WebFrame/dyndom"

	"github.com/PuerkitoBio/goquery"
	"github.com/dennwc/dom"
	"golang.org/x/net/html"

	"gitlab.com/256/DebateFrame/client/filesaver"
	"gitlab.com/256/DebateFrame/client/log"
	"gitlab.com/256/DebateFrame/client/mammoth"
	"gitlab.com/256/DebateFrame/client/waiter"
	"gitlab.com/256/WebFrame/waquery"
)

func eventListeners(container *dyndom.Element) error {
	drop := newDrop()
	cb := js.NewEventCallback(0, func(file js.Value) {
		log.DebugMessage("Drop event fired")
		fullFile := file.Get("name").String()
		fileNoExt := strings.TrimSuffix(fullFile, filepath.Ext(fullFile))
		fileExt := filepath.Ext(fullFile)
		switch fileExt {
		case ".docx":
			mammoth.ConvertToHTML(&file, func(html string) {
				html = postProcess(html)
				cs, err := NewCase(html, fileNoExt)
				if err != nil {
					log.PanicMessage("Failed to create new case from HTML", err)
				}
				err = cs.Add()
				if err != nil {
					log.PanicMessage("Failed to add the docx document to the screen", err)
				}
				err = cs.SetActive()
				if err != nil {
					log.PanicMessage("Failed to set the docx file as the current case!", err)
				}
			})
		case ".dfc":
			log.DebugMessage("DebateFrame case detected!")
			caseLoad(&file)
		}

	})
	drop.Call("on", "addedfile", cb)

	waiter.EventWaiter(&dom.GetDocument().NodeBase, "scroll", 50, func() {
		activeLinks := dom.GetDocument().QuerySelectorAll(".is-active-link")
		if len(activeLinks) > 0 {
			active := activeLinks[0]
			if !waquery.IsVisible(waquery.ToHTML(active)) {
				topPos := active.JSValue().Get("offsetTop")
				currentCase.TOCElem.JSValue().Set("scrollTop", topPos)
			}
		}
	})

	return nil
}

// blobToBytes converts a Blob to []byte.
func blobToBytes(blob js.Value) []byte {
	js.Global().Call("blobToBytes", blob)
	var res string
	for {
		time.Sleep(time.Millisecond * 50)
		res = js.Global().Get("result").String()
		if res != "" {
			js.Global().Set("result", "")
			break
		}
	}
	bytesStr := strings.Split(res, ",")
	byteSlice := []byte{}
	for _, byteStr := range bytesStr {
		u, _ := strconv.ParseUint(byteStr, 0, 8)
		byteSlice = append(byteSlice, byte(u))
	}
	return byteSlice
}

var template = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Document</title>
</head>
<body>
    
</body>
</html>
`

// OnSave is the event listener for when the save button is pressed
func OnSave(e dom.Event) {
	docHTML := currentCase.EditorElem.JSValue().Get("innerHTML").String()

	templateDoc, err := html.Parse(strings.NewReader(template))
	if err != nil {
		panic("couldn't parse template")
	}

	goDoc := goquery.NewDocumentFromNode(templateDoc)
	goDoc.Find("body").AppendHtml(docHTML)
	genHTML, err := goDoc.Html()
	if err != nil {
		panic("couldn't generate html from document")
	}
	filesaver.Save([]byte(genHTML), "test.html", "text/html")
}

var toggle = false

// OnReadMode is the event listener for when the ReadMode button is pressed
func OnReadMode(e dom.Event) {
	if !toggle {
		dom.GetDocument().GetElementById("editorSwitcher").ClassList().Add("readmode")
		dom.GetDocument().QuerySelector("#readMode").SetTextContent("Default View")
		waquery.AddStyle(fmt.Sprintf("%s mark { visibility:visible !important; }", docQuery))
	} else {
		dom.GetDocument().GetElementById("editorSwitcher").ClassList().Add("readmode")
		dom.GetDocument().QuerySelector("#readMode").SetTextContent("Read Mode")
	}
	toggle = !toggle
}
