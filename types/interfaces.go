package types

import "golang.org/x/net/html"

type Html interface {
	HtmlNode() *html.Node
}
