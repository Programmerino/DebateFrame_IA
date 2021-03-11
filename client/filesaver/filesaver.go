package filesaver

import (
	"encoding/base64"
	"gitlab.com/256/WebFrame/dyndom"
	"syscall/js"

	"github.com/dennwc/dom"
)

// Save saves contents from bytes form to a specified file with the given mimetype
func Save(contents []byte, filename string, mimetype string) {
	b64 := base64.StdEncoding.EncodeToString(contents)
	blob := Base64ToBlob(b64, mimetype)
	elem := dyndom.CreateElement("a")
	elem.SetAttribute("href", js.Global().Get("window").Get("URL").Call("createObjectURL", *blob))
	elem.SetAttribute("download", filename)
	elem.Style().Set("display", "none")
	dom.GetDocument().GetElementsByTagName("body")[0].AppendChild(&elem.Element)
	elem.JSValue().Call("click")
	dom.GetDocument().GetElementsByTagName("body")[0].RemoveChild(&elem.Element)
}

// Base64ToBlob converts a base64 string to a JavaScript blob
func Base64ToBlob(b64Data string, mimeType string) *js.Value {
	blob := js.Global().Call("b64toBlob", b64Data, mimeType)
	return &blob
}
