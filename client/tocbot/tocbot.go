package tocbot

import (
	"regexp"
	"syscall/js"

	"fmt"

	"gitlab.com/256/DebateFrame/client/log"

	"github.com/dennwc/dom"
)

/*
tocbot.init({
  // Where to render the table of contents.
  tocSelector: '.js-toc',
  // Where to grab the headings to build the table of contents.
  contentSelector: '.js-toc-content',
  // Which headings to grab inside of the contentSelector element.
  headingSelector: 'h1, h2, h3',

  collapseDepth: 6,
});

tocbot.refresh()

*/

// GenerateTOC creates a table of contents from the given document selector, and places it in the given selector
func GenerateTOC(doc string, to string) {
	elem := dom.GetDocument().QuerySelector(doc)

	// tocbot requires all headers to have an id to jump to location
	nodes := elem.ChildNodes()
	var header = regexp.MustCompile(`H\d`)
	for i, node := range nodes {
		tag := node.JSValue().Get("tagName").String()
		if header.MatchString(tag) {
			node.SetAttribute("id", fmt.Sprintf("h-%v", i))
		}
	}

	input := make(map[string]interface{})
	input["tocSelector"] = to
	input["contentSelector"] = doc
	input["headingSelector"] = "h1, h2, h3, h4, h5, h6"
	input["collapseDepth"] = 6
	js.Global().Get("tocbot").Call("init", input)
	elem.AddEventListener("input", func(e dom.Event) {
		log.DebugMessage("Refreshed table of contents")
		js.Global().Get("tocbot").Call("refresh")
	})
}
