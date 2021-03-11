package config

import (
	"encoding/json"
	"strings"

	"fmt"

	"github.com/dennwc/dom"
	"github.com/dennwc/dom/storage"
	"github.com/pkg/errors"

	"gitlab.com/256/DebateFrame/client/log"
)

// CurrentConfig is the current Configuration of DebateFrame
var CurrentConfig Configuration

// LocalStorage is the JS equivalent of LocalStorage
var LocalStorage storage.Storage

// Configuration holds everything needed to replicate a DebateFrame instance
// All fields that should be saved must start with an uppercase letter to be saved by json
type Configuration struct {
	FinishedWizard bool
}

func init() {
	LocalStorage = storage.Local()
	doc := dom.GetWindow()
	doc.AddEventListener("beforeunload", func(dom.Event) {
		err := CurrentConfig.saveState()
		if err != nil {
			panic(errors.Wrap(err, "could not save the current Configuration to LocalStorage"))
		}
	})
}

// saveState saves the Configuration to LocalStorage
func (config *Configuration) saveState() error {
	str, err := jsonString(*config)
	if err != nil {
		return errors.Wrap(err, "failed to get json string of Configuration")
	}

	LocalStorage.SetItem("config", str)

	return nil
}

// jsonString converts the Configuration to a string that can be saved in a location and later loaded
func jsonString(config Configuration) (string, error) {
	cfgString, err := json.Marshal(config)
	if err != nil {
		return "", errors.Wrap(err, "failed to convert object to json")
	}

	return string(cfgString), nil
}

// RestoreState gets the last Configuration used, and makes the website use it
func RestoreState() error {
	var err error
	CurrentConfig, err = restoreFromStorage()
	if err != nil {
		return errors.Wrap(err, "failed to restore Configuration from local storage")
	}
	log.WarnMessage("RestoreState doesn't work as intended")
	return nil
}

// restoreFromStorage gets the current Configuration from storage and sets the Configuration to it
func restoreFromStorage() (Configuration, error) {
	jsonString, ok := LocalStorage.GetItem("config")
	if !ok {
		return Configuration{}, fmt.Errorf("failed to load json from local storage")
	}

	var conf Configuration
	dec := json.NewDecoder(strings.NewReader(jsonString))
	err := dec.Decode(&conf)
	if err != nil {
		log.PanicMessage("failed to decode Configuration from json. panicking because given config may be in a bad state now", err)
	}

	return conf, nil
}
