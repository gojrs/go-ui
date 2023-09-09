package routing

import (
	"golang.org/x/net/html"
	"io"
)

type NodeRouter interface {
	Name() string
	RouteToPath(path string)
	Start()
	RegisterPath(path string, component NodeRender) error
}

type NodeRender interface {
	Render() io.Reader
	Destroy()
	Init()
	Name() string
	Guard(userName string) bool
	GetViewNode() *html.Node
}
