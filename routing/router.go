package routing

import (
	"github.com/gojrs/go-ui/types"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type WasmRouter struct {
	viewId, name string
	comChan      chan *storageRequest
	vNode        *html.Node
}

func NewWasmRouter(id, name string) *WasmRouter {
	w := &WasmRouter{
		viewId:  id,
		name:    name,
		comChan: make(chan *storageRequest, 1),
		vNode: types.NewElementFromAtom(atom.Div, html.Attribute{
			Key: "id",
			Val: id,
		}),
	}
	go w.startStorage(false)
	return w
}

func (wr *WasmRouter) Name() string                                   { return wr.name }
func (wr *WasmRouter) RouteToPath(path string)                        {}
func (wr *WasmRouter) SetViewNodeId(id string)                        { wr.viewId = id }
func (wr *WasmRouter) GetViewNode() *html.Node                        { return wr.vNode }
func (wr *WasmRouter) RegisterPath(path string, component NodeRender) {}

func (wr *WasmRouter) Start() {
	channel := make(chan struct{})

	<-channel
}
