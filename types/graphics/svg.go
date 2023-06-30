package graphics

import (
	"golang.org/x/net/html"
)

type SvgElement struct {
	svgNode *html.Node
	Height  uint
	Width   uint
	Defs    map[string]*html.Node
}

func (svg *SvgElement) CreateDef(kind, id string, attrs ...html.Attribute) {

}
