package document

import (
	"fmt"

	"gitlab.com/256/DebateFrame/client/document/card"
	"gitlab.com/256/DebateFrame/client/log"
	"gitlab.com/256/DebateFrame/client/waiter"
	"gitlab.com/256/WebFrame/dyndom"
)

func cardView(cards []*card.Card) *dyndom.Element {
	cView := cardViewElem()
	fmt.Println(cView.Children("div"))
	parent := cView.Children("div")[1]
	search := cView.Child("div")
	for _, card := range cards {
		if card.Element == nil {
			card.GenerateElement()
		}
		parent.AppendChild(card.Element)
	}
	lastQuery := ""

	waiter.EventWaiter(&search.NodeBase, "input", 300, func() {
		query := search.JSValue().Get("value").String()
		if query != "" && lastQuery != query {
			go inputEvent(cards, query)
			lastQuery = query
		} else if query == "" {
			for i, card := range cards {
				if card.Element == nil {
					log.PanicMessage("Given card has no associated element! Please associate an element with the card first!", nil)
				} else {
					card.Element.ClassList().Add("nonmatch")
					card.Element.Style().Set("order", i)
				}
			}
		}
	})

	return cView
}

/*
<div>
   <div class="uk-search uk-search-large">
       <span uk-search-icon=""></span>
       <input id="cardSearch" class="uk-search-input" type="search" placeholder="Search..." />
   </div>
   <div id="cards" class="uk-flex uk-flex-around uk-flex-wrap">

   </div>
</div>
*/

func cardViewElem() *dyndom.Element {
	div := dyndom.CreateElement("div")
	div.AppendChild(searchDivElem())
	div.AppendChild(cardDivElem())
	return div
}

func searchDivElem() *dyndom.Element {
	div := dyndom.CreateElement("div", "uk-search", "uk-search-large")
	icon := dyndom.CreateElement("span")
	icon.SetAttribute("uk-search-icon", "")
	input := dyndom.CreateElement("input")
	input.SetId("cardSearch")
	input.ClassList().Add("uk-search-input")
	input.SetAttribute("type", "search")
	input.SetAttribute("placeholder", "Search...")
	div.AppendChild(icon)
	div.AppendChild(input)
	return div
}

func cardDivElem() *dyndom.Element {
	div := dyndom.CreateElement("div", "uk-flex", "uk-flex-around", "uk-flex-wrap")
	div.SetId("cards")
	return div
}

var running = false

func inputEvent(cards []*card.Card, str string) {
	if running {
		return
	}
	running = true
	card.Filter(cards, str)
	go func() {
		running = false
	}()
}
