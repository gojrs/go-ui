package types

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"syscall/js"
)

type ElementWrapper struct {
	js.Value
}

func NewNodeFromAtom(a atom.Atom, attrs ...html.Attribute) *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: a,
		Data:     a.String(),
		Attr:     attrs,
	}
}

func NewElementFromAtom(a atom.Atom, attrs ...html.Attribute) ElementWrapper {
	elem := js.Global().Get("document").Call("createElement", a.String())
	for _, attribute := range attrs {
		elem.Set(attribute.Key, attribute.Val)
	}
	return ElementWrapper{Value: elem}
}

func NewElementFromNode(src *html.Node) ElementWrapper {
	elem := js.Global().Get("document").Call("createElement", src.DataAtom.String())

	for _, attribute := range src.Attr {
		elem.Set(attribute.Key, attribute.Val)
	}
	return ElementWrapper{elem}
}
