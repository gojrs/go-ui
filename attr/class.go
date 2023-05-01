package attr

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"strings"
)

func AddClass(c string, nodes ...*html.Node) {
	for _, node := range nodes {
		aClass := atom.Class
		if node.Attr != nil {
			for i, attribute := range node.Attr {
				if attribute.Key == aClass.String() {
					if strings.Contains(attribute.Val, c) {
						return
					}
					node.Attr[i] = html.Attribute{
						Key: aClass.String(),
						Val: strings.Join([]string{attribute.Val, c}, " "),
					}
				} else {
					node.Attr = append(node.Attr, html.Attribute{
						Key: aClass.String(),
						Val: c,
					})
				}
			}

		} else {
			node.Attr = []html.Attribute{{Key: atom.Class.String(), Val: c}}
		}
	}
}

func RemoveClass(c string, nodes ...*html.Node) {
	for _, node := range nodes {
		aClass := atom.Class
		if node.Attr != nil {
			for i, attribute := range node.Attr {
				if attribute.Key == aClass.String() {
					before, after, didCut := strings.Cut(attribute.Val, c)
					if !didCut {
						return
					}
					node.Attr[i] = html.Attribute{
						Key: aClass.String(),
						Val: strings.Join([]string{before, after}, " "),
					}
				}
			}

		} else {
			continue
		}
	}
}

func SetClasses(classes string, nodes ...*html.Node) {
	for _, n := range nodes {
		classes = strings.TrimSpace(classes)
		if classes == "" {
			RemoveAttr(atom.Class.String(), n)
			continue
		}

		n.Attr = append(n.Attr, html.Attribute{
			Key: atom.Class.String(),
			Val: classes,
		})
	}
}

func AddStyle(s string, nodes ...*html.Node) {
nodeLoop:
	for _, node := range nodes {
		attrPtr := getAttrPtrForKey(node, atom.Style.String())
		//lookingFor := fmt.Sprintf("%s;", s)
		if attrPtr == nil {
			node.Attr = append(node.Attr, html.Attribute{
				Key: atom.Style.String(),
				Val: s,
			})
			continue
		}

		list := strings.Split(attrPtr.Val, ";")
		for i, style := range list {
			if style == s {
				continue nodeLoop
			}
			kv := strings.Split(s, ":")

			if strings.HasPrefix(style, kv[0]) {
				list[i] = s
				attrPtr.Val = strings.Join(list, ";")
				continue nodeLoop
			}
		}
		attrPtr.Val = strings.Join([]string{attrPtr.Val, s}, ";")
	}
}

func RemoveAttr(attrName string, nodes ...*html.Node) {
	for _, n := range nodes {
		for i, a := range n.Attr {
			if a.Key == attrName {
				n.Attr[i], n.Attr[len(n.Attr)-1], n.Attr =
					n.Attr[len(n.Attr)-1], html.Attribute{}, n.Attr[:len(n.Attr)-1]
			}
		}
	}
}

func AddAttr(attr html.Attribute, nodes ...*html.Node) {
	for _, node := range nodes {
		RemoveAttr(attr.Key, node)
		node.Attr = append(node.Attr, attr)
	}
}
