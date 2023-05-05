package routing

import (
	"golang.org/x/net/html"
	"io"
)

type NodeRouter interface {
	Name() string
	RouteToPath(path string)
	GetViewNode() *html.Node
	RegisterPath(path string, component NodeRender) error
}

type NodeRender interface {
	Render() io.Reader
	Destroy()
	Name() string
	Guard(userName string) bool
}
