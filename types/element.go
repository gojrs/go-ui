package types

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func NewElementFromAtom(a atom.Atom, attrs ...html.Attribute) *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: a,
		Data:     a.String(),
		Attr:     attrs,
	}
}
