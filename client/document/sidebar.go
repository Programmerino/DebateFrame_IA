package document

import (
	"github.com/dennwc/dom"
)

func generateSidebars() {
	btn := newButton("Save")
	btn.OnClick(OnSave)
	dom.GetDocument().GetElementById("appSideBar").AppendChild(btn)
	btn = newButton("Read Mode")
	btn.OnClick(OnReadMode)
	dom.GetDocument().GetElementById("appSideBar").AppendChild(btn)
	btn.Element.SetId("readMode")
}
