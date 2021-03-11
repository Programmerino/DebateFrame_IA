package document

import (
	"syscall/js"

	"gitlab.com/256/DebateFrame/client/medium"
)

// Creates a new medium editor in the provided queryselector
func newEditor(query string) *medium.Editor {
	options := medium.DefaultOptions()
	options.ButtonLabels = medium.FontAwesome
	options.Extensions = []medium.Extension{{Name: "highlighter", Value: js.Global().Get("HighlighterButton").New()}}
	options.Toolbar.Buttons = []string{"bold", "highlighter", "underline", "italic", "h2", "h3", "quote", "anchor"}
	options.TargetBlank = true
	options.AutoLink = true
	options.AnchorPreview.Enabled = true
	options.AnchorPreview.ShowOnEmptyLinks = false
	options.Anchor.LinkValidation = true
	options.Paste.ForcePlainText = false
	options.Paste.CleanPastedHTML = true
	editor := medium.NewEditor(query, options)
	return editor
}
