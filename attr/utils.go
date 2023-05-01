package attr

import (
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"strings"
)

func GetElementById(n *html.Node, id string) *html.Node {

	return TraverseForAttr(n, atom.Id.String(), id)
}

func GetNodeByClass(n *html.Node, class string) *html.Node {

	return TraverseForAttr(n, atom.Class.String(), class)
}

func GetNodesByClass(n *html.Node, class string) Selection {

	return traverseForNodesForClass(n, class, make(Selection, 0, 10))
}

func GetNodesByAtom(n *html.Node, kind atom.Atom) {
	traverseForNodesForAtom(n, kind, make(Selection, 0, 10))
}

func getAttrPtrForKey(n *html.Node, key string) *html.Attribute {
	for _, attr := range n.Attr {

		if attr.Key == key {
			return &attr
		}
	}
	return nil
}

func GetAttribute(n *html.Node, key string) (string, bool) {

	for _, attr := range n.Attr {

		if attr.Key == key {
			return attr.Val, true
		}
	}

	return "", false
}

func checkId(n *html.Node, id string) bool {

	if n.Type == html.ElementNode {

		s, ok := GetAttribute(n, atom.Id.String())

		if ok && s == id {
			return true
		}
	}

	return false
}

func checkClass(n *html.Node, class string) bool {

	if n.Type == html.ElementNode {

		s, ok := GetAttribute(n, atom.Class.String())

		if ok && strings.Contains(s, fmt.Sprintf(" %s ", class)) {
			return true
		}
	}

	return false
}

func TraverseForAttr(n *html.Node, attr, id string) *html.Node {
	switch attr {
	case atom.Id.String():
		if checkId(n, id) {
			return n
		}
	case atom.Class.String():

	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {

		res := TraverseForAttr(c, attr, id)

		if res != nil {
			return res
		}
	}

	return nil
}

func traverseForNodesForClass(n *html.Node, class string, found Selection) Selection {

	if checkClass(n, class) {
		found = append(found, n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {

		traverseForNodesForClass(c, class, found)

	}

	return found
}

func traverseForNodesForAtom(n *html.Node, of atom.Atom, found Selection) Selection {

	if n.DataAtom == of && n.Data == of.String() {
		found = append(found, n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {

		traverseForNodesForAtom(c, of, found)

	}

	return found
}
