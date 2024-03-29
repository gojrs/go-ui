package routing

import (
	"github.com/gojrs/go-ui/attr"
	"github.com/gojrs/go-ui/types"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"syscall/js"
)

var currentRenderer NodeRender

type Opts struct {
	viewId, name string
	comChan      chan *storageRequest
	vNode        *html.Node
	vChildNode   *html.Node
}

type OptFunc func(opts *Opts)

func defaultOpts() Opts {
	id := "DEFAULT-VIEW-ID"
	return Opts{
		viewId:  id,
		name:    "DEFAULT-NAME",
		comChan: make(chan *storageRequest, 1),
		vNode: types.NewNodeFromAtom(atom.Div, html.Attribute{
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
	emptyObj := make(map[string]any)
	js.Global().Set("Router", emptyObj)
	ns := js.Global().Get("Router")
	ns.Set("Link", js.FuncOf(w.Link))
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

func (wr *WasmRouter) Name() string { return wr.name }
func (wr *WasmRouter) RouteToPath(path string) {
	respC := make(chan *storageResponse, 1)
	wr.comChan <- &storageRequest{
		fullPath: path,
		reqType:  storageReqFetch,
		reply:    respC,
	}
	answer := <-respC
	if wr.vChildNode != nil {
		wr.vNode.RemoveChild(wr.vChildNode)
	}

	if currentRenderer != nil {
		currentRenderer.Destroy()
	}

	if answer.err != nil {
		println(answer.err.Error())
		return
	}
	currentRenderer = answer.component
	currentRenderer.Init()
	wr.vChildNode = currentRenderer.GetViewNode()
	wr.vNode.AppendChild(currentRenderer.GetViewNode())

	bs, err := io.ReadAll(currentRenderer.Render())
	if err != nil {
		println(err.Error())
		return
	}

	viewJS := docJs.Call("getElementById", wr.viewId)
	if !viewJS.Truthy() {
		println("could not find", wr.viewId)
		return
	}
	viewJS.Set("innerHTML", string(bs))
}

func (wr *WasmRouter) Link(this js.Value, i []js.Value) interface{} {
	var (
		path = i[0].String()
	)
	wr.RouteToPath(path)

	return nil
}
func (wr *WasmRouter) SetViewNodeId(id string) { wr.viewId = id }
func (wr *WasmRouter) GetViewNode() *html.Node { return wr.vNode }
func (wr *WasmRouter) RegisterPath(path string, component NodeRender) error {
	respChan := make(chan *storageResponse)
	wr.comChan <- &storageRequest{
		reqType:   storageReqIns,
		fullPath:  path,
		component: component,
		reply:     respChan,
	}
	ok := <-respChan
	if ok != nil {
		if ok.err != nil {
			return ok.err
		}
	}

	return nil
}

func (wr *WasmRouter) Start() {
	channel := make(chan struct{})

	var (
		objName = "Router"
		ns      = js.Global().Get(objName)
	)
	if !ns.Truthy() {
		emptyObj := make(map[string]any)
		js.Global().Set(objName, emptyObj)
		js.Global().Get(objName).Set("Link", wr.Link)
	}

	<-channel
}

type registerRouteRequest struct {
	path      string
	component NodeRender
}

type routeToRequest struct {
	path   string
	params []js.Value
	respC  chan *routeToResponse
}
type routeToResponse struct {
	component NodeRender
	params    []js.Value
	err       error
}
