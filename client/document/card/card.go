package card

import (
	"fmt"

	"gitlab.com/256/WebFrame/dyndom"
)

// if strictCards is set to true, then only cards with properties such as the author and year will be added
const strictCards = true

// Card represents a card of evidence
type Card struct {
	Title    string
	Contents string
	URL      string
	Year     uint8
	Author   string
	Element  *dyndom.Element // A reference to the actual element on the page
}

// GenerateElement generates an element card from a debate card
func (card *Card) GenerateElement() {
	cardDiv := dyndom.CreateElement("div", "uk-card", "uk-card-body", "uk-card-default",
		"uk-card-hover", "uk-height-medium", "uk-card-small",
		"docCard")
	title := dyndom.CreateElement("h3", "uk-card-title")
	title.SetTextContent(cleanString(card.Title))
	cardDiv.AppendChild(title)
	author := dyndom.CreateElement("h4")
	author.SetTextContent(cleanString(fmt.Sprintf("%s %v", card.Author, card.Year)))
	cardDiv.AppendChild(author)
	// The contents need to be in a "read more" thing
	/*
		content := dom.NewElement("p")
		content.SetTextContent(card.Contents)
		cardDiv.AppendChild(content)
	*/
	card.Element = cardDiv
}

const limit = 80

func cleanString(str string) string {
	if len(str) > limit {
		str = str[0:limit] + "..."
	}
	return str
}
