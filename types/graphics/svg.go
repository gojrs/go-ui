package graphics

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type SvgElement struct {
	svgNode *html.Node
	Height  uint
	Width   uint
	Defs    map[string]*html.Node
}

func (svg *SvgElement) CreateDef(kind, id string, attrs ...html.Attribute) {
	me := atom.Shape
}
