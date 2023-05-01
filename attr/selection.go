package attr

import "golang.org/x/net/html"

type Selection []*html.Node

func (sel Selection) AddClass(c string) {
	AddClass(c, sel...)
}

func (sel Selection) RemoveClass(c string) {
	RemoveClass(c, sel...)
}

func (sel Selection) AddAttr(attr html.Attribute) {
	AddAttr(attr, sel...)
}

func (sel Selection) SetClasses(c string) {
	SetClasses(c, sel...)
}

func (sel Selection) AddStyle(s string) {
	AddStyle(s, sel...)
}
