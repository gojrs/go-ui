package input

import (
	"fmt"
	"github.com/gojrs/go-ui/types"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type DescriberTable []DescriberItem
type DescriberItem map[string]ColumnDescriber

type ItemsTable struct {
	iTable, parent, head, body *html.Node
	id                         string
	columns                    []ColumnDescriber
	schemaData                 DescriberTable // outer slice == each tr; map key is column key
}

func (its *ItemsTable) HtmlNode() *html.Node {
	return its.iTable
}

func NewTableElem(id, cap string, pn *html.Node) *ItemsTable {

	var (
		tFoot   = types.NewNodeFromAtom(atom.Tfoot)
		caption = types.NewNodeFromAtom(atom.Caption)

		it = &ItemsTable{
			iTable: types.NewNodeFromAtom(atom.Table, html.Attribute{
				Key: atom.Id.String(),
				Val: fmt.Sprintf("%s-%s", itemTablePrefix, id),
			}),
			id:         id,
			parent:     pn,
			head:       types.NewNodeFromAtom(atom.Thead),
			body:       types.NewNodeFromAtom(atom.Tbody),
			schemaData: make(DescriberTable, 0, 10),
		}
	)

	it.parent.AppendChild(it.iTable)

	caption.AppendChild(&html.Node{
		Type:     html.TextNode,
		DataAtom: 0,
		Data:     cap,
	})
	it.iTable.AppendChild(caption)
	it.iTable.AppendChild(it.head)
	it.iTable.AppendChild(it.body)
	it.iTable.AppendChild(tFoot)
	it.addActionElem()
	return it
}

func (its *ItemsTable) UpdateHeaderData(cols []ColumnDescriber) *ItemsTable {
	its.columns = cols

	row := types.NewNodeFromAtom(atom.Tr, html.Attribute{
		Key: "data-collection-id",
		Val: its.id,
	})
	selectRows := NewInputNode(&BaseInput{
		Id:       fmt.Sprintf("table-thead-tr-%s", its.id),
		Label:    "Select All",
		Type:     "checkbox",
		Title:    "",
		Hidden:   false,
		Readonly: false,
		Required: false,
		Disabled: false,
		System:   false,
		Classes:  nil,
	})
	firstChild := types.NewNodeFromAtom(atom.Th)
	for _, selectRow := range selectRows {
		firstChild.AppendChild(selectRow)
	}
	row.AppendChild(firstChild)
	for _, column := range its.columns {
		colKey := column.ColumnKey()
		th := types.NewNodeFromAtom(atom.Th, html.Attribute{
			Key: "class",
			Val: fmt.Sprintf("col-%s", colKey),
		})
		h3 := types.NewNodeFromAtom(atom.H3)
		h3.AppendChild(&html.Node{
			Type:     html.TextNode,
			DataAtom: 0,
			Data:     column.GetLabel(),
		})
		th.AppendChild(h3)
		row.AppendChild(th)
	}
	its.head.AppendChild(row)
	return its
}

func (its *ItemsTable) LoadData(data DescriberTable) error {
	its.schemaData = data
	for _, rowData := range its.schemaData {
		var (
			tr        = types.NewNodeFromAtom(atom.Tr)
			firstTd   = types.NewNodeFromAtom(atom.Td)
			lastTd    = types.NewNodeFromAtom(atom.Td)
			selectRow = types.NewNodeFromAtom(atom.Select)
		)
		firstTd.AppendChild(selectRow)
		tr.AppendChild(firstTd)

		for _, column := range its.columns {
			base, ok := rowData[column.ColumnKey()]
			if !ok {
				return fmt.Errorf("could not find column %s", column.ColumnKey())
			}
			inputs, err := NewInputsFromDescriber(base)
			if err != nil {
				return fmt.Errorf("could not create inputs %s", err.Error())
			}
			td := types.NewNodeFromAtom(atom.Td)
			for _, nodes := range inputs {
				for _, node := range nodes {
					td.AppendChild(node)
				}
			}
		}
		tr.AppendChild(lastTd)
		its.body.AppendChild(tr)
	}
	return nil
}

func (its *ItemsTable) getBody() *html.Node {
	return its.body
}

func (its *ItemsTable) getFoot() (foot *html.Node) {
	foot = its.iTable.LastChild
	for foot.DataAtom != atom.Tfoot {
		foot = foot.PrevSibling
	}
	return foot
}

func (its *ItemsTable) addActionElem() {

}
