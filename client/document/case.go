package document

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/PuerkitoBio/goquery"
	"github.com/dennwc/dom"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"gitlab.com/256/DebateFrame/client/document/card"
	"gitlab.com/256/DebateFrame/client/filesaver"
	"gitlab.com/256/DebateFrame/client/log"
	"gitlab.com/256/DebateFrame/client/medium"
	"gitlab.com/256/DebateFrame/client/tocbot"
	"gitlab.com/256/WebFrame/dyndom"
	"gitlab.com/256/WebFrame/waquery"
)

// Case is an object that holds information relating to a debate case
type Case struct {
	Name       string
	Cards      []*InfoCard
	Document   *goquery.Document
	Button     *dyndom.Element // The button that switches to the case
	Editor     *medium.Editor
	EditorElem *dyndom.Element
	TOCElem    *dyndom.Element
}

// NewCase creates a new Case object from an HTML string
func NewCase(html string, name string) (*Case, error) {
	cs := Case{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create document from html")
	}
	cs.Document = doc
	cs.Cards = toInfoCards(card.GetCards(doc))
	cs.Name = name
	return &cs, nil
}

// Add adds the case to the screen as an option
func (cs *Case) Add() error {
	html, err := cs.Document.Html()
	if err != nil {
		return errors.Wrap(err, "Failed to convert the case document to HTML")
	}

	link, uuid := newTabLink(cs.Name)
	dom.GetDocument().GetElementById("doctab").AppendChild(&link.Element)
	cs.Button = dyndom.New(waquery.ToHTML(dom.GetDocument().GetElementById(fmt.Sprintf("%s-link", uuid))))

	// EDITOR

	container := newEditorContainer(toCards(cs.Cards))
	dom.GetDocument().GetElementById("editorSwitcher").AppendChild(&container.Element)

	editorID := fmt.Sprintf("%s_editor", dyndom.UUID())
	container.Children("div")[1].SetId(editorID)
	editorQuery := fmt.Sprintf("#%s", editorID)
	editor := newEditor(editorQuery)
	cs.Editor = editor
	cs.EditorElem = container.Children("div")[1]
	cs.Editor.SetContent(html, 0)

	cViewButton := container.Children("div")[0].Child("a")
	cViewButton.AddEventListener("click", func(e dom.Event) {
		fmt.Println(cViewButton.GetAttribute("uk-icon").String())
		if strings.Contains(cViewButton.GetAttribute("uk-icon").String(), "file-text") {
			container.Children("div")[2].ClassList().Add("simplehide")
			container.Children("div")[1].ClassList().Remove("simplehide")
			cViewButton.SetAttribute("uk-icon", "icon: fa-sticky-note-s; ratio: 2")
		} else {
			container.Children("div")[1].ClassList().Add("simplehide")
			container.Children("div")[2].ClassList().Remove("simplehide")
			cViewButton.SetAttribute("uk-icon", "icon: file-text; ratio: 2")
		}
	})

	// TOC

	tocCont, uuid := newTOCCont()
	dom.GetDocument().GetElementById("tocSwitcher").AppendChild(&tocCont.Element)

	tocQuery := fmt.Sprintf("#%s-tocDiv", uuid)
	cs.TOCElem = dyndom.New(waquery.ToHTML(dom.GetDocument().QuerySelector(tocQuery)))

	tocbot.GenerateTOC(editorQuery, tocQuery)

	cardView(toCards(cs.Cards))

	return nil
}

/*
   <div class="uk-card uk-card-body" id="tocCont" uk-height-viewport="expand: true">
       <div id="toc" uk-sticky="offset: 100"></div>
   </div>
*/
func newTOCCont() (*dyndom.Element, string) {
	uuid := newUUID()
	tocCont := dyndom.CreateElement("div", "uk-card", "uk-card-body", "tocCont")
	tocCont.SetAttribute("uk-height-viewport", "expand: true")
	tocCont.SetId(fmt.Sprintf("%s-tocCont", uuid))
	tocCont.AppendChild(newTocDiv(uuid))
	return tocCont, uuid
}

func newTocDiv(uuid string) *dyndom.Element {
	tocDiv := dyndom.CreateElement("div", "tocdiv")
	tocDiv.SetAttribute("uk-sticky", "offset: 100")
	tocDiv.SetId(fmt.Sprintf("%s-tocDiv", uuid))
	return tocDiv
}

func newUUID() string {
	return fmt.Sprintf("a%s", uuid.New().String())
}

func newTabLink(text string) (*dyndom.Element, string) {
	uuid := newUUID()
	item := dyndom.CreateElement("li")
	item.SetId(fmt.Sprintf("%s-item", uuid))
	link := dyndom.CreateElement("a")
	link.SetAttribute("href", "#")
	link.SetInnerHTML(text)
	link.SetId(fmt.Sprintf("%s-link", uuid))
	item.AppendChild(link)
	return item, uuid
}

func newEditorContainer(cards []*card.Card) *dyndom.Element {
	editorCont := dyndom.CreateElement("div", "uk-card", "uk-card-body", "editorCont")
	editorCont.AppendChild(newEditorToolbar())
	editorCont.AppendChild(newEditorDiv())
	cView := cardView(cards)
	cView.ClassList().Add("simplehide")
	editorCont.AppendChild(cView)
	return editorCont
}

func newEditorToolbar() *dyndom.Element {
	toolbarDiv := dyndom.CreateElement("div", "toolbar")
	toolbarDiv.AppendChild(newCardViewButton())
	toolbarDiv.AppendChild(newDownloadButton())
	return toolbarDiv
}

func newCardViewButton() *dyndom.Element {
	cardViewButton := newToolbarButton("fa-sticky-note-s")
	return cardViewButton
}

func newDownloadButton() *dyndom.Element {
	downloadButton := newToolbarButton("download")
	downloadButton.AddEventListener("click", func(e dom.Event) {
		caseSave()
	})
	return downloadButton
}

func newToolbarButton(iconName string) *dyndom.Element {
	button := dyndom.CreateElement("a", "uk-icon", "toolbarButton")
	button.SetAttribute("href", "#")
	button.SetAttribute("uk-icon", fmt.Sprintf("icon: %s; ratio: 2", iconName))
	return button
}

func newEditorDiv() *dyndom.Element {
	div := dyndom.CreateElement("div", "editor")
	return div
}

// Update updates a debate case with the new HTML contents
func (cs *Case) Update(html string) error {
	var err error
	cs, err = NewCase(html, cs.Name)
	if err != nil {
		return errors.Wrap(err, "failed to update the debate case")
	}
	return nil
}

// SetActive sets the case provided as the active case that is on screen
func (cs *Case) SetActive() error {
	currentCase = cs
	cs.Button.JSValue().Call("click")
	return nil
}

// SaveableCase is a version of Case that stores the document as HTML instead of as a goquery.Document to avoid issues with recursion limits
type SaveableCase struct {
	Name     string
	Cards    []*InfoCard
	Document string
}

// Saveable returns a version of the case that is saveable
func (cs *Case) Saveable() *SaveableCase {
	var err error
	scase := SaveableCase{}
	scase.Cards = cs.Cards
	scase.Document, err = cs.Document.Html()
	if err != nil {
		log.PanicMessage("Failed to convert the goquery document to HTML", err)
	}
	scase.Name = cs.Name
	return &scase
}

// Normalize converts a Saveable case to one that can be used by DebateFrame
func (saveable *SaveableCase) Normalize() *Case {
	var err error
	cs := Case{}
	cs.Name = saveable.Name
	cs.Cards = saveable.Cards
	cs.Document, err = goquery.NewDocumentFromReader(strings.NewReader(saveable.Document))
	if err != nil {
		log.PanicMessage("Failed to convert the document HTML to a goquery document", err)
	}
	return &cs
}

// InfoCard represents a card of evidence without the attached element
type InfoCard struct {
	Title    string
	Contents string
	URL      string
	Year     uint8
	Author   string
}

// InfoCard -> Card

// toCard converts an info card to a normal card
func (icard *InfoCard) toCard() *card.Card {
	card := card.Card{}
	card.Title = icard.Title
	card.Contents = icard.Contents
	card.URL = icard.URL
	card.Year = icard.Year
	card.Author = icard.Author
	return &card
}

// toInfoCards converts normal cards to the InfoCard format
func toCards(cards []*InfoCard) []*card.Card {
	cardSlice := []*card.Card{}
	for _, icard := range cards {
		cardSlice = append(cardSlice, icard.toCard())
	}
	return cardSlice
}

// Card -> InfoCard

// toInfoCard converts a normal card to the InfoCard format
func toInfoCard(card *card.Card) *InfoCard {
	icard := InfoCard{}
	icard.Title = card.Title
	icard.Contents = card.Contents
	icard.URL = card.URL
	icard.Year = card.Year
	icard.Author = card.Author
	return &icard
}

// toInfoCards converts normal cards to the InfoCard format
func toInfoCards(cards []*card.Card) []*InfoCard {
	icardSlice := []*InfoCard{}
	for _, tcard := range cards {
		icardSlice = append(icardSlice, toInfoCard(tcard))
	}
	return icardSlice
}

// Saving and Loading functionality

func caseSave() {
	log.DebugMessage("Case Save initiated!")
	log.DebugMessage("Converting to gob")
	scase := currentCase.Saveable()
	bytes, err := objToBytes(&scase)
	if err != nil {
		log.PanicMessage("Failed to convert the object into a Gob byte array", err)
	}
	log.DebugMessage("Saving gob")
	filesaver.Save(bytes, "Case.dfc", "application/vnd.dframe-case")
}

func caseLoad(file *js.Value) error {
	temp := SaveableCase{}
	err := bytesToObj(blobToBytes(*file), &temp)
	if err != nil {
		fmt.Println("error encountered")
	}
	cs := temp.Normalize()
	err = cs.Add()
	if err != nil {
		return errors.Wrap(err, "Failed to add the DebateFrame case to the list of cases")
	}
	err = cs.SetActive()
	if err != nil {
		return errors.Wrap(err, "Failed to set the DebateFrame format case as the current case")
	}
	/*
		log.DebugMessage("Case Save initiated!")
		log.DebugMessage("Converting to gob")
		scase := currentCase.Saveable()
		bytes, err := objToBytes(&scase)
		if err != nil {
			log.PanicMessage("Failed to convert the object into a Gob byte array", err)
		}
		log.DebugMessage("Saving gob")
		filesaver.Save(bytes, "Case.dfc", "application/vnd.dframe-case")
	*/
	return nil
}
