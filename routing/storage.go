package routing

import (
	"bytes"
	"encoding/json"
)

const (
	storageReqIns storeReqType = iota + 1
	storageReqUpdate
	storageReqRemove
	storageReqFetch
)

type storeReqType uint

type storageRequest struct {
	reqType   storeReqType
	fullPath  string
	component NodeRender
	reply     chan *storageResponse
}

type storageResponse struct {
	err       error
	component NodeRender
}

func (wr *WasmRouter) startStorage(persist bool) {
	var privateStore routeTrunk
	var err error
	if persist {
		privateStore, err = readFromStorage()
		if err != nil {
			println("could not find previous storage, creating new one")
			privateStore = make(routeTrunk)
		}
	} else {
		privateStore = make(routeTrunk)
	}

	for request := range wr.comChan {
		var resp storageResponse
		switch request.reqType {
		case storageReqFetch:
			resp.component, resp.err = privateStore.fetch(request.fullPath)
			request.reply <- &resp
		case storageReqIns:
			privateStore.insert(request.fullPath, request.component)
			request.reply <- &resp
		case storageReqRemove:
			privateStore.remove(request.fullPath)
			request.reply <- &resp
		case storageReqUpdate:
			privateStore.update(request.fullPath, request.component)
			request.reply <- &resp
		}
	}
	//close all connections
	if persist {
		var (
			buffer = bytes.NewBufferString("")
		)
		err = json.NewEncoder(buffer).Encode(&privateStore)
		if err != nil || buffer.Len() > 0 {
			println("could not encode storage", err.Error())
			return
		}
		err = writeToStorage(buffer)
		if err != nil {
			println("could not save to storage", err.Error())
			return
		}
	}
}
