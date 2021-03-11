package main

import (
	"gitlab.com/256/WebFrame/dyndom"
	"syscall/js"

	"github.com/dennwc/dom"
	"github.com/pkg/errors"
	"gitlab.com/256/DebateFrame/client/config"
	"gitlab.com/256/DebateFrame/client/document"
	"gitlab.com/256/DebateFrame/client/grabber"
	"gitlab.com/256/DebateFrame/client/log"
	"gitlab.com/256/DebateFrame/client/state" 
	"gitlab.com/256/DebateFrame/client/wizard"
)

// main is ran when the page finished loading 
func main() {
	log.DebugMessage("Client side code started!")
	err := getDeps()
	if err != nil {
		panic(errors.Wrap(err, "failed to load needed resources"))
	}
	err = config.RestoreState()
	if err != nil {
		log.WarnMessage(err.Error())
		log.DebugMessage("Failed to restore previous state... Clearing LocalStorage and starting over")
		config.LocalStorage.Clear()
		config.CurrentConfig = config.Configuration{}
	}
	if !config.CurrentConfig.FinishedWizard { 
		state.StartWrap(state.Setup, func(container *dyndom.Element) {
			config.CurrentConfig = wizard.Start(container)
		})
	}

	state.StartWrap(state.DocumentWriter, func(container *dyndom.Element) {
		err := document.Run(container)
		if err != nil {
			panic(errors.Wrap(err, "failed to run the document state"))
		}
	})

	grabber.Grab()
	js.Global().Get("UIkit").Call("modal", dom.GetDocument().GetElementById("modal-loading").JSValue()).Call("hide")
	select {}
}

func getDeps() error {
	log.DebugMessage("No dependencies to get!")
	return nil
}
