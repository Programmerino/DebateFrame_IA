package document

import (
	"syscall/js"

	"gitlab.com/256/WebFrame/dyndom"

	"github.com/davecgh/go-xdr/xdr"
	"github.com/dennwc/dom"
	"github.com/pkg/errors"
)

var currentCase *Case

// Run starts the document writer code
func Run(container *dyndom.Element) error {
	initCase, err := NewCase("", "Untitled Document")
	if err != nil {
		return errors.Wrap(err, "Failed to create new blank case")
	}
	err = initCase.Add()
	if err != nil {
		return errors.Wrap(err, "Failed to add the case to the screen")
	}
	err = initCase.SetActive()
	if err != nil {
		return errors.Wrap(err, "Failed to set the initial case as the active case")
	}
	err = eventListeners(container)
	if err != nil {
		return errors.Wrap(err, "failed to create event listeners")
	}

	generateSidebars()
	return nil
}

func newDrop() *js.Value {
	input := make(map[string]interface{})

	input["url"] = "#"

	drop := js.Global().Get("Dropzone").New(dom.GetDocument().GetElementsByTagName("body")[0].JSValue(), input)
	return &drop
}

func newButton(content string) *dom.Button {
	btn := dom.NewButton(content)
	btn.SetAttribute("class", "uk-button uk-width-1-1")
	return btn
}

const compression = 1

func objToBytes(object interface{}) ([]byte, error) {
	objbytes, err := xdr.Marshal(object)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode object to bytes")
	}
	return objbytes, nil
	/*
		var b bytes.Buffer
		// Anything more than 2 just makes the result file 0 bytes long...
		w, err := bzip2.NewWriter(&b, &bzip2.WriterConfig{Level: compression})
		if err != nil {
			return nil, errors.Wrap(err, "failed to initialize the compressor")
		}
		defer w.Close()
		_, err = w.Write(objbytes)
		if err != nil {
			return nil, errors.Wrap(err, "failed to apply compression on xdr marshalled-file")
		}

		return b.Bytes(), nil
	*/
}

func bytesToObj(xdrBytes []byte, object interface{}) error {
	/*
		rd, err := bzip2.NewReader(bytes.NewReader(xdrBytes), &bzip2.ReaderConfig{})
		defer rd.Close()
		if err != nil {
			return errors.Wrap(err, "unexpected error when creating new bzip2 reader")
		}

		output, err := ioutil.ReadAll(rd)
		if err != nil {
			return errors.Wrap(err, "failed to read data from the bzip2 reader")
		}
	*/
	_, err := xdr.Unmarshal(xdrBytes, object)
	if err != nil {
		return errors.Wrap(err, "failed to decode bytes to Go object")
	}
	return err
}
