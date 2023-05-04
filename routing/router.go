package routing

import (
	"github.com/gojrs/go-ui/attr"
	"github.com/gojrs/go-ui/types"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Opts struct {
	viewId, name string
	comChan      chan *storageRequest
	vNode        *html.Node
}

type OptFunc func(opts *Opts)

func defaultOpts() Opts {
	id := "DEFAULT-ID"
	return Opts{
		viewId:  id,
		name:    "DEFAULT-NAME",
		comChan: make(chan *storageRequest, 1),
		vNode: types.NewElementFromAtom(atom.Div, html.Attribute{
			Key: "id",
			Val: id,
		}),
	}
}

type WasmRouter struct {
	Opts
}

func NewWasmRouter(opts ...OptFunc) *WasmRouter {
	o := defaultOpts()
	for _, fn := range opts {
		fn(&o)
	}
	w := &WasmRouter{
		Opts: o,
	}
	go w.startStorage(false)
	return w
}

func WithName(name string) OptFunc {
	return func(opts *Opts) {
		opts.name = name
	}
}

func WitId(id string) OptFunc {
	return func(opts *Opts) {
		attr.AddAttr(html.Attribute{
			Key: "id",
			Val: id,
		}, opts.vNode)
		opts.viewId = id
	}
}

func WithViewNode(node *html.Node) OptFunc {
	return func(opts *Opts) {
		id, ok := attr.GetAttribute(opts.vNode, "id")
		if ok {
			opts.viewId = id
		} else {
			attr.AddAttr(html.Attribute{
				Key: "id",
				Val: id,
			}, opts.vNode)
		}
		opts.vNode = node
	}
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
