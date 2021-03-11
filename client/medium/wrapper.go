package medium

import (
	"syscall/js"
	"time"

	"github.com/dennwc/dom"
)

/*
var quill = new Quill('#editor', {
  modules: {
    toolbar: false    // Snow includes toolbar by default
  },
  theme: 'snow'
});
*/
type Extension struct {
	Name  string
	Value js.Value
}

// ButtonLabelType is an enum holding the valid values for the ButtonLabels option
type ButtonLabelType int32

const (
	// Default button labels
	Default ButtonLabelType = iota
	// FontAwesome Uses fontawesome icon set for all toolbar icons
	FontAwesome
)

func (label *ButtonLabelType) Value() js.Value {
	switch *label {
	case Default:
		return js.ValueOf(false)
	case FontAwesome:
		return js.ValueOf("fontawesome")
	}
	return js.ValueOf(false)
}

// Options that can be used when creating a new editor
type EditorOptions struct {
	ActiveButtonClass   string               // CSS class added to active buttons in the toolbar.
	ButtonLabels        ButtonLabelType      // Custom content for the toolbar buttons.
	ContentWindow       *dom.Window          // The contentWindow object that contains the contenteditable element. MediumEditor will use this for attaching events, getting selection, etc.
	Delay               time.Duration        // Time to show the toolbar or anchor tag preview.
	DisableReturn       bool                 // Enables/disables the use of the return-key. You can also set specific element behavior by using setting a data-disable-return attribute.
	DisableDoubleReturn bool                 // Allows/disallows two (or more) empty new lines. You can also set specific element behavior by using setting a data-disable-double-return attribute.
	DisableExtraSpaces  bool                 // When set to true, it disallows spaces at the beginning and end of the element. Also it disallows entering 2 consecutive spaces between 2 words.
	DisableEditing      bool                 // Enables/disables adding the contenteditable behavior. Useful for using the toolbar with customized buttons/actions. You can also set specific element behavior by using setting a data-disable-editing attribute.
	ElementsContainer   dom.Node             // Specifies a DOM node to contain MediumEditor's toolbar and anchor preview elements.
	Extensions          []Extension          // Custom extensions to use. See Custom Buttons and Extensions for more details on extensions.
	OwnerDocument       *dom.Document        // The ownerDocument object for the contenteditable element. MediumEditor will use this for creating elements, getting selection, attaching events, etc.
	Spellcheck          bool                 // Enable/disable native contentEditable automatic spellcheck.
	TargetBlank         bool                 // Enables/disables automatically adding the target="_blank" attribute to anchor tags.
	Toolbar             ToolbarOptions       // The toolbar for MediumEditor is implemented as a built-in extension which automatically displays whenever the user selects some text. The toolbar can hold any set of defined built-in buttons, but can also hold any custom buttons passed in as extensions.
	AnchorPreview       AnchorPreviewOptions // The anchor preview is a built-in extension which automatically displays a 'tooltip' when the user is hovering over a link in the editor. The tooltip will display the href of the link, and when click, will open the anchor editing form in the toolbar.
	Placeholder         PlaceholderOptions   // The placeholder handler is a built-in extension which displays placeholder text when the editor is empty.
	Anchor              AnchorOptions        // The anchor form is a built-in button extension which allows the user to add/edit/remove links from within the editor. When 'anchor' is passed in as a button in the list of buttons, this extension will be enabled and can be triggered by clicking the corresponding button in the toolbar.
	Paste               PasteOptions         // The paste handler is a built-in extension which attempts to filter the content when the user pastes. How the paste handler filters is configurable via specific options.
	KeyboardCommands    KeyboardOptions      // The keyboard commands handler is a built-in extension for mapping key-combinations to actions to execute in the editor.
	AutoLink            bool                 // The auto-link handler is a built-in extension which automatically turns URLs entered into the text field into HTML anchor tags (similar to the functionality of Markdown). This feature is OFF by default.
	ImageDragging       bool                 // The image dragging handler is a built-in extension for handling dragging & dropping images into the contenteditable. This feature is ON by default.
}

// Value returns the JavaScript value of the options
func (opts *EditorOptions) Value() js.Value {
	optMap := make(map[string]interface{})
	optMap["activeButtonClass"] = opts.ActiveButtonClass
	optMap["buttonLabels"] = opts.ButtonLabels.Value()
	optMap["contentWindow"] = js.Value(opts.ContentWindow.JSValue())
	optMap["delay"] = int(opts.Delay / time.Millisecond)
	optMap["disableReturn"] = opts.DisableReturn
	optMap["disableDoubleReturn"] = opts.DisableDoubleReturn
	optMap["disableExtraSpaces"] = opts.DisableExtraSpaces
	optMap["disableEditing"] = opts.DisableEditing
	optMap["elementsContainer"] = js.Value(opts.ElementsContainer.JSValue())

	extMap := make(map[string]interface{})
	for _, extension := range opts.Extensions {
		extMap[extension.Name] = extension.Value
	}
	optMap["extensions"] = extMap

	optMap["ownerDocument"] = js.Value(opts.OwnerDocument.JSValue())
	optMap["spellcheck"] = opts.Spellcheck
	optMap["targetBlank"] = opts.TargetBlank
	optMap["toolbar"] = opts.Toolbar.Value()
	optMap["anchorPreview"] = opts.AnchorPreview.Value()
	optMap["placeholder"] = opts.Placeholder.Value()
	optMap["anchor"] = opts.Anchor.Value()
	optMap["paste"] = opts.Paste.Value()
	optMap["keyboardCommands"] = opts.KeyboardCommands.Value()
	optMap["autoLink"] = opts.AutoLink
	optMap["imageDragging"] = opts.ImageDragging

	return js.ValueOf(optMap)
}

type ToolbarOptions struct {
	Enabled                      bool          // To disable the toolbar (which also disables the anchor-preview extension), set the value of the Enabled option to false
	AllowMultiParagraphSelection bool          // enables/disables whether the toolbar should be displayed when selecting multiple paragraphs/block elements.
	Buttons                      []string      // The set of buttons to display on the toolbar.
	DiffLeft                     int           // Value in pixels to be added to the X axis positioning of the toolbar.
	DiffTop                      int           // Value in pixels to be added to the Y axis positioning of the toolbar.
	FirstButtonClass             string        // CSS class added to the first button in the toolbar.
	LastButtonClass              string        // CSS class added to the last button in the toolbar.
	RelativeContainer            dom.Element   // DOMElement to append the toolbar to instead of the body. When an element is passed the toolbar will also be positioned relative instead of absolute, which means the editor will not attempt to manually position the toolbar automatically.
	StandardizeSelectionStart    bool          // Enables/disables standardizing how the beginning of a range is decided between browsers whenever the selected text is analyzed for updating toolbar buttons status.
	Static                       StaticOptions // Enable/disable the toolbar always displaying in the same location relative to the medium-editor element.
}

func StrSlice(slice []string) []interface{} {
	interfaceSlice := make([]interface{}, len(slice))
	for i, d := range slice {
		interfaceSlice[i] = d
	}
	return interfaceSlice
}

// Value returns the JavaScript value of the options
func (opts *ToolbarOptions) Value() js.Value {
	if !opts.Enabled {
		return js.ValueOf(false)
	}
	optMap := make(map[string]interface{})
	optMap["allowMultiParagraphSelection"] = js.ValueOf(opts.AllowMultiParagraphSelection)
	optMap["buttons"] = js.ValueOf(StrSlice(opts.Buttons))
	optMap["diffLeft"] = js.ValueOf(opts.DiffLeft)
	optMap["diffTop"] = js.ValueOf(opts.DiffTop)
	optMap["firstButtonClass"] = opts.FirstButtonClass
	optMap["lastButtonClass"] = opts.LastButtonClass
	optMap["relativeContainer"] = js.Value(opts.RelativeContainer.JSValue())
	optMap["standardizeSelectionStart"] = opts.StandardizeSelectionStart
	optMap["static"] = opts.Static.Enabled
	if opts.Static.Enabled {
		optMap["align"] = opts.Static.Align.Value()
		optMap["sticky"] = opts.Static.Sticky
		optMap["stickyTopOffset"] = opts.Static.StickyTopOffset
		optMap["updateOnEmptySelection"] = opts.Static.UpdateOnEmptySelection
	}

	return js.ValueOf(optMap)
}

// AlignType is an enum holding the valid values for the Align option
type AlignType int32

const (
	Left AlignType = iota
	Center
	Right
)

func (align *AlignType) Value() js.Value {
	switch *align {
	case Left:
		return js.ValueOf("left")
	case Center:
		return js.ValueOf("center")
	case Right:
		return js.ValueOf("right")
	}
	return js.ValueOf("")
}

type StaticOptions struct {
	Enabled                bool      // Enable/disable the toolbar always displaying in the same location relative to the medium-editor element.
	Align                  AlignType // When the static option is true, this aligns the static toolbar relative to the medium-editor element.
	Sticky                 bool      // When the static option is true, this enables/disables the toolbar "sticking" to the viewport and staying visible on the screen while the page scrolls.
	StickyTopOffset        int       // When the sticky option is true, this set in pixel a top offset above the toolbar.
	UpdateOnEmptySelection bool      // When the static option is true, this enables/disables updating the state of the toolbar buttons even when the selection is collapsed (there is no selection, just a cursor).
}

type AnchorPreviewOptions struct {
	Enabled                  bool          // To disable the anchor preview, set the value of the Enabled option to false
	HideDelay                time.Duration // Time to show the anchor tag preview after the mouse has left the anchor tag.
	PreviewValueSelector     string        // The default selector to locate where to put the activeAnchor value in the preview. You should only need to override this if you've modified the way in which the anchor-preview extension renders.
	ShowOnEmptyLinks         bool          // Determines whether the anchor tag preview shows up on link with href as "" or "#something". You should set this value to false if you do not want the preview to show up in such use cases.
	ShowWhenToolbarIsVisible bool          // Determines whether the anchor tag preview shows up when the toolbar is visible. You should set this value to true if the static option for the toolbar is true and you want the preview to show at the same time.
}

// Value returns the JavaScript value of the options
func (opts *AnchorPreviewOptions) Value() js.Value {
	if !opts.Enabled {
		return js.ValueOf(false)
	}
	optMap := make(map[string]interface{})
	optMap["hideDelay"] = int64(opts.HideDelay / time.Millisecond)
	optMap["previewValueSelector"] = opts.PreviewValueSelector
	optMap["showOnEmptyLinks"] = opts.ShowOnEmptyLinks
	optMap["showWhenToolbarIsVisible"] = opts.ShowWhenToolbarIsVisible

	return js.ValueOf(optMap)
}

type PlaceholderOptions struct {
	Enabled     bool   // To disable the placeholder, set the value of the Enabled option to false
	Text        string // Defines the default placeholder for empty contenteditables when placeholder is not set to false. You can overwrite it by setting a data-placeholder attribute on the editor elements.
	HideOnClick bool   // Causes the placeholder to disappear as soon as the field gains focus. To hide the placeholder only after starting to type, and to show it again as soon as field is empty, set this option to false.
}

// Value returns the JavaScript value of the options
func (opts *PlaceholderOptions) Value() js.Value {
	if !opts.Enabled {
		return js.ValueOf(false)
	}
	optMap := make(map[string]interface{})
	optMap["text"] = opts.Text
	optMap["hideOnClick"] = opts.HideOnClick

	return js.ValueOf(optMap)
}

type AnchorOptions struct {
	CustomClassOption     string // Custom class name the user can optionally have added to their created links (ie 'button'). If passed as a non-empty string, a checkbox will be displayed allowing the user to choose whether to have the class added to the created link or not.
	CustomClassOptionText string // Text to be shown in the checkbox when the CustomClassOption is being used.
	LinkValidation        bool   // Enables/disables check for common URL protocols on anchor links. Converts invalid url characters (ie spaces) to valid characters using encodeURI
	PlaceholderText       string // Text to be shown as placeholder of the anchor input.
	TargetCheckbox        bool   // Enables/disables displaying a "Open in new window" checkbox, which when checked changes the target attribute of the created link.
	TargetCheckboxText    string // Text to be shown in the checkbox enabled via the TargetCheckbox option.
}

// Value returns the JavaScript value of the options
func (opts *AnchorOptions) Value() js.Value {
	optMap := make(map[string]interface{})
	optMap["customClassOption"] = opts.CustomClassOption
	optMap["customClassOptionText"] = opts.CustomClassOptionText
	optMap["linkValidation"] = opts.LinkValidation
	optMap["placeholderText"] = opts.PlaceholderText
	optMap["targetCheckbox"] = opts.TargetCheckbox
	optMap["targetCheckboxText"] = opts.TargetCheckboxText

	return js.ValueOf(optMap)
}

type Replacement struct {
	RegExp      string
	Replacement string
}

// Value returns the JavaScript value of the Replacement
func (repl *Replacement) Value() js.Value {
	replRay := []string{repl.RegExp, repl.Replacement}
	return js.ValueOf(replRay)
}

type PasteOptions struct {
	ForcePlainText    bool          // Forces pasting as plain text.
	CleanPastedHTML   bool          // Cleans pasted content from different sources, like google docs etc.
	CleanReplacements []Replacement // Replacement structs to use during paste when ForcePlainText or CleanPastedHTML are true OR when calling cleanPaste(text) helper method.
	CleanAttrs        []string      // List of element attributes to remove during paste when CleanPastedHTML is true or when calling cleanPaste(text) or pasteHTML(html,options) helper methods.
	CleanTags         []string      // List of element tag names to remove during paste when CleanPastedHTML is true or when calling cleanPaste(text) or pasteHTML(html,options) helper methods.
	UnwrapTags        []string      // List of element tag names to unwrap (remove the element tag but retain its child elements) during paste when CleanPastedHTML is true or when calling cleanPaste(text) or pasteHTML(html,options) helper methods.
}

// Value returns the JavaScript value of the options
func (opts *PasteOptions) Value() js.Value {
	optMap := make(map[string]interface{})
	optMap["forcePlainText"] = opts.ForcePlainText
	optMap["cleanPastedHTML"] = opts.CleanPastedHTML
	optMap["cleanAttrs"] = StrSlice(opts.CleanAttrs)
	optMap["cleanTags"] = StrSlice(opts.CleanTags)
	optMap["unwrapTags"] = StrSlice(opts.UnwrapTags)

	replRay := []interface{}{}
	for _, repl := range opts.CleanReplacements {
		replRay = append(replRay, repl.Value())
	}
	optMap["cleanReplacements"] = replRay
	return js.ValueOf(optMap)
}

type Binding struct {
	Command string // argument passed to editor.execAction() when key-combination is used
	Key     rune   // keyboard character that triggers this command
	Meta    bool   // whether the ctrl/meta key has to be active or inactive
	Shift   bool   // whether the shift key has to be active or inactive
}

// Value returns the JavaScript value of the binding
func (bind *Binding) Value() js.Value {
	bindingMap := make(map[string]interface{})
	bindingMap["command"] = bind.Command
	bindingMap["key"] = string(bind.Key)
	bindingMap["meta"] = bind.Meta
	bindingMap["shift"] = bind.Shift
	return js.ValueOf(bindingMap)
}

type KeyboardOptions struct {
	Enabled  bool // To disable the keyboard commands, set the value of the Enabled option to false:
	Commands []Binding
}

// Value returns the JavaScript value of the options
func (opts *KeyboardOptions) Value() js.Value {
	if !opts.Enabled {
		return js.ValueOf(false)
	}
	optMap := make(map[string]interface{})

	bindingRay := []interface{}{}
	for _, binding := range opts.Commands {
		bindingRay = append(bindingRay, binding.Value())
	}

	optMap["commands"] = bindingRay
	return js.ValueOf(optMap)

}

type Editor struct {
	inst *js.Value
}

// Editor creates a new instance of MediumEditorq
func NewEditor(querySelector string, options EditorOptions) *Editor {
	inst := js.Global().Get("MediumEditor").New(querySelector, options.Value())
	return &Editor{inst: &inst}
}

// Destroy tears down the editor if already setup
func (editor *Editor) Destroy() {
	editor.inst.Call("destroy")
}

// Setup initializes this instance of the editor if it has been destroyed.
func (editor *Editor) Setup() {
	editor.inst.Call("setup")
}

// AddElements will dynamically add one or more elements to an already initialized instance of MediumEditor.
func (editor *Editor) AddElements(elements ...dom.Element) {
	interfaceSlice := make([]interface{}, len(elements))
	for i, d := range elements {
		interfaceSlice[i] = js.Value(d.JSValue())
	}
	editor.inst.Call("addElements", interfaceSlice)
}

// RemoveElements will remove one or more elements from an already initialized instance of MediumEditor.
func (editor *Editor) RemoveElements(elements ...dom.Element) {
	interfaceSlice := make([]interface{}, len(elements))
	for i, d := range elements {
		interfaceSlice[i] = js.Value(d.JSValue())
	}
	editor.inst.Call("removeElements", interfaceSlice)
}

// SetContent sets the html content for the first editor element, or the element at index. Ensures the the editableInput event is triggered.
func (editor *Editor) SetContent(html string, index int) {
	editor.inst.Call("setContent", html, index)
}

// GetContent returns the trimmed html content for the first editor element, or the element at index.
func (editor *Editor) GetContent(index int) string {
	return editor.inst.Call("getContent", index).String()
}

// DefaultOptions returns an EditorOptions with all the default options selected
func DefaultOptions() EditorOptions {
	options := EditorOptions{}
	options.ActiveButtonClass = "medium-editor-button-active"
	options.ContentWindow = dom.GetWindow()
	bodyNode := dom.Node(&dom.Body.NodeBase)
	options.ElementsContainer = bodyNode
	options.OwnerDocument = dom.GetDocument()
	options.Spellcheck = true
	toolbar := ToolbarOptions{}
	toolbar.Enabled = true
	toolbar.AllowMultiParagraphSelection = true
	toolbar.Buttons = []string{"bold", "italic", "underline", "anchor", "h2", "h3", "quote"}
	toolbar.DiffTop = -10
	toolbar.FirstButtonClass = "medium-editor-button-first"
	toolbar.LastButtonClass = "medium-editor-button-last"
	toolbar.Static.Align = Center
	options.Toolbar = toolbar
	anPrev := AnchorPreviewOptions{}
	anPrev.Enabled = true
	anPrev.HideDelay = time.Millisecond * 500
	anPrev.PreviewValueSelector = "a"
	anPrev.ShowOnEmptyLinks = true
	options.AnchorPreview = anPrev
	placeholder := PlaceholderOptions{}
	placeholder.Enabled = true
	placeholder.Text = "Type your text"
	placeholder.HideOnClick = true
	options.Placeholder = placeholder
	anchor := AnchorOptions{}
	anchor.CustomClassOptionText = "Button"
	anchor.PlaceholderText = "Paste or type a link"
	anchor.TargetCheckboxText = "Open in new window"
	options.Anchor = anchor
	paste := PasteOptions{}
	paste.ForcePlainText = true
	paste.CleanAttrs = []string{"class", "style", "dir"}
	paste.CleanTags = []string{"meta"}
	options.Paste = paste
	keyb := KeyboardOptions{}
	keyb.Commands = []Binding{
		Binding{
			Command: "bold",
			Key:     'b',
			Meta:    true,
		},
		Binding{
			Command: "italic",
			Key:     'i',
			Meta:    true,
		},
		Binding{
			Command: "underline",
			Key:     'u',
			Meta:    true,
		},
	}
	options.KeyboardCommands = keyb
	options.ImageDragging = true
	return options
}
