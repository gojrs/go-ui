package input

import (
	"fmt"
	"github.com/gojrs/go-ui/types"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	itemTablePrefix = "itemTable"
)

type Noder interface {
	HtmlNode() *html.Node
}

type ItemTable struct {
	iTable, parent, body *html.Node
	id                   string
	schemaData           []ColumnDescriber
}

func (it *ItemTable) HtmlNode() *html.Node {
	return it.iTable
}

func NewFormElem(id, cap string, pn *html.Node) *ItemTable {

	var (
		foot    = types.NewElementFromAtom(atom.Tfoot)
		caption = types.NewElementFromAtom(atom.Caption)
		it      = &ItemTable{
			iTable: types.NewElementFromAtom(atom.Table, html.Attribute{
				Key: atom.Id.String(),
				Val: fmt.Sprintf("%s-%s", itemTablePrefix, id),
			}),
			id:         id,
			parent:     pn,
			body:       types.NewElementFromAtom(atom.Tbody),
			schemaData: nil,
		}
	)

	it.parent.AppendChild(it.iTable)

	caption.AppendChild(&html.Node{
		Type:     html.TextNode,
		DataAtom: 0,
		Data:     cap,
	})
	it.iTable.AppendChild(caption)
	it.iTable.AppendChild(it.body)
	it.iTable.AppendChild(foot)
	it.addActionButtons()
	return it
}

func (it *ItemTable) WithSchemaData(kinders ...ColumnDescriber) *ItemTable {
	it.schemaData = kinders
	return it
}

func (it *ItemTable) BuildForm() error {
	fields, err := NewInputsFromDescriber(it.schemaData...)
	if err != nil {
		return err
	}

	for colKey, set := range fields {
		var tr *html.Node
		for _, node := range set {

			thId := fmt.Sprintf("%s-th-%s", itemTablePrefix, colKey)
			switch node.DataAtom {
			case atom.Label:
				tr = types.NewElementFromAtom(atom.Tr, html.Attribute{
					Key: atom.Id.String(),
					Val: fmt.Sprintf("%s-tr-%s", itemTablePrefix, colKey),
				})
				it.getBody().AppendChild(tr)
				th := types.NewElementFromAtom(atom.Th, html.Attribute{
					Key: atom.Id.String(),
					Val: thId,
				}, html.Attribute{
					Key: atom.Class.String(),
					Val: thId,
				}, html.Attribute{
					Key: atom.Scope.String(),
					Val: atom.Col.String(),
				})
				th.AppendChild(node)
				tr.AppendChild(th)
			case atom.Input, atom.Select, atom.Textarea:
				td := types.NewElementFromAtom(atom.Td, html.Attribute{
					Key: atom.Headers.String(),
					Val: thId,
				}, html.Attribute{
					Key: atom.Class.String(),
					Val: thId,
				})
				td.AppendChild(node)
				tr.AppendChild(td)
			case atom.Datalist:
				body := it.getBody()

				body.LastChild.LastChild.AppendChild(node)

			default:
				err = fmt.Errorf("unknown tag type %s", node.DataAtom.String())
			}
		}
	}

	return nil
}

func (it *ItemTable) AppendChildToBody(node *html.Node) {
	it.getBody().AppendChild(node)
}

func (it *ItemTable) addActionButtons() {
	var (
		foot     = it.getFoot()
		dataAttr = html.Attribute{
			Key: fmt.Sprintf("data-%s-id", itemTablePrefix),
			Val: it.id,
		}
		containerDiv = types.NewElementFromAtom(atom.Div, html.Attribute{
			Key: atom.Class.String(),
			Val: fmt.Sprintf("%s-btn-group", itemTablePrefix),
		})
		saveBtn = types.NewElementFromAtom(atom.Button, html.Attribute{
			Key: atom.Onclick.String(),
			Val: fmt.Sprintf("%s.%s.save()", itemTablePrefix, it.id),
		}, dataAttr)
		deleteBtn = types.NewElementFromAtom(atom.Button, html.Attribute{
			Key: atom.Onclick.String(),
			Val: fmt.Sprintf("%s.%s.delete()", itemTablePrefix, it.id),
		}, dataAttr)
	)

	// we found the foot
	foot.AppendChild(containerDiv)

	saveBtn.AppendChild(&html.Node{
		Type:     html.TextNode,
		DataAtom: 0,
		Data:     "Save",
	})

	deleteBtn.AppendChild(&html.Node{
		Type:     html.TextNode,
		DataAtom: 0,
		Data:     "Delete",
	})
	containerDiv.AppendChild(saveBtn)
	containerDiv.AppendChild(deleteBtn)
}

func (it *ItemTable) getBody() *html.Node {
	return it.body
}

func (it *ItemTable) getFoot() (foot *html.Node) {
	foot = it.iTable.LastChild
	for foot.DataAtom != atom.Tfoot {
		foot = foot.PrevSibling
	}
	return foot
}
