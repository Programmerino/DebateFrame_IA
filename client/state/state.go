package state

import (
	"gitlab.com/256/WebFrame/dyndom"
	"time"

	"fmt"

	"gitlab.com/256/DebateFrame/client/log"
	"gitlab.com/256/WebFrame/waquery"

	"github.com/dennwc/dom"
)

func init() {
	for _, st := range systemStates {
		st.hide()
	}
}

type state struct {
	query string
}

// Default states
var (
	Setup          = state{query: "#setup"}
	DocumentWriter = state{query: "#document"}
)

const (
	defaultView = "flex" // The default display method for elements made visible
)

var visibleStates []state

var systemStates = []state{Setup}

// StartWrap runs the function when the state is ready to start
func StartWrap(st state, fn func(*dyndom.Element)) {
	switchTo(st)
	fn(st.getElement())
}

// switchTo hides all other states and displays the one given
func switchTo(st state) {
	hideAll()
	st.display()
}

// hideAll hides all visible states
func hideAll() {
	for _, vState := range visibleStates {
		vState.hide()
	}
}

func (st *state) display() {
	log.DebugMessage(fmt.Sprintf("State with selector %s is now on display", st.query))
	waquery.FadeIn(&st.getElement().HTMLElement, time.Second)
	visibleStates = append(visibleStates, *st)
}

func (st *state) hide() {
	log.DebugMessage(fmt.Sprintf("State with selector %s is now hidden", st.query))
	waquery.FadeOut(&st.getElement().HTMLElement, time.Second)
	remove(visibleStates, *st)
}

func (st *state) getElement() *dyndom.Element {
	elem := dom.GetDocument().QuerySelector(string(st.query))
	if elem == nil {
		panic(fmt.Sprintf("State query \"%s\" was not found! Cannot proceed...", string(st.query)))
	}
	return dyndom.New(waquery.ToHTML(elem))
}

func remove(s []state, r state) []state {
	for i, v := range s {
		if v.query == r.query {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
