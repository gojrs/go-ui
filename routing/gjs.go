package routing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"syscall/js"
)

var (
	storageKey = "jrs-route-storage"
	lStore     = js.Global().Get("window").Get("localStorage")
	docJs      = js.Global().Get("document")
)

func writeToStorage(reader io.Reader) error {

	if !lStore.Truthy() {
		return fmt.Errorf("could not get local storage")
	}
	str, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	lStore.Call("setItem", storageKey, string(str))
	return nil
}

func readFromStorage() (router routeTrunk, err error) {
	jItem := lStore.Call("getItem", storageKey)
	if !jItem.Truthy() {
		err = fmt.Errorf("could not get local from storage")
		return router, err
	}
	buffer := bytes.NewBufferString("")
	err = json.NewDecoder(buffer).Decode(router)
	if err != nil {
		err = fmt.Errorf("could not decode %s", err.Error())
		return router, err
	}
	return router, err
}
