package wizard

import (
	"gitlab.com/256/WebFrame/dyndom"
	"time"

	"gitlab.com/256/DebateFrame/client/config"

	"gitlab.com/256/DebateFrame/client/log"

	"github.com/dennwc/dom"
	"gitlab.com/256/WebFrame/waquery"
)

// Start is the initialization function for the wizard
func Start(container *dyndom.Element) config.Configuration {
	welcome(container)
	preferences(container)
	return config.Configuration{FinishedWizard: true}
}

func welcome(container *dyndom.Element) {
	welcome := waquery.H1("Welcome to DebateFrame")
	container.AppendChild(dyndom.New(welcome))
	time.Sleep(time.Second * 2)
	waquery.FadeOut(welcome, time.Second)
}

func preferences(container *dyndom.Element) {
	form := form()
	form.Style().Set("display", "none")
	form.Style().Set("opacity", 0)
	container.AppendChild(dyndom.New(form))
	waquery.FadeIn(form, time.Second)
	time.Sleep(time.Second * 1)
	waquery.FadeOut(form, time.Second)
}

func form() *dom.HTMLElement {
	log.WarnMessage("Form is not implemented!")
	return waquery.H1("Just tell us a few things and you're set! (not implemented)")
	/*
		form := dom.NewElement("form")
		margin := divMargin()
	*/
}

func divMargin() *dyndom.Element {
	div := dyndom.CreateElement("div", "uk-margin")
	return div
}
