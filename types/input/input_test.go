package input_test

import (
	"bytes"
	"fmt"
	"github.com/gojrs/go-ui/types"
	"github.com/gojrs/go-ui/types/input"
	"github.com/google/uuid"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"strings"
	"testing"
)

var (
	exampleBase = input.BaseInput{
		Id:          "exampleBase-ID",
		Value:       "exampleBase=Val",
		Key:         "exampleBase",
		PlaceHolder: "exampleBase-PH",
		Label:       "exampleBase-Label",
		Type:        "text",
		Title:       "exampleBase-title",
		Hidden:      false,
		Readonly:    false,
		Required:    false,
		Disabled:    true,
		System:      false,
		Classes:     []string{"row", "col-key-id"},
		Class:       "",
		Options:     nil,
	}
	exampleSelect = input.SelectInput{
		BaseInput: exampleBase,
		Multiple:  false,
		Options: []input.OptionParam{
			{
				Key:   "1",
				Value: "One",
			},
			{
				Key:   "2",
				Value: "Two",
			},
			{
				Key:   "3",
				Value: "Three",
			},
		},
	}
	exampleTextArea = input.TextAreaInput{
		BaseInput: exampleBase,
		Cols:      20,
		Rows:      30,
		Maxlength: 500,
		Minlength: 5,
	}
	exampleWithList = input.WithListOptions{
		BaseInput: exampleBase,
		List: input.DataList{
			Id: "dataList-id",
			Options: []input.OptionParam{
				{
					Key:   "red",
					Value: "Red",
				},
				{
					Key:   "green",
					Value: "Green",
				},
				{
					Key:   "blue",
					Value: "Blue",
				},
			},
		},
	}
)

func makeHeader() []input.ColumnDescriber {
	var (
		baseNode = exampleBase
		//checkBoxNode = exampleBase
		selectNode = exampleSelect
		//textAreaNode = exampleTextArea
		withListNode = exampleWithList
		schema       = []map[string]any{
			map[string]interface{}{
				"system":   false,
				"id":       "mmuzrsrp",
				"name":     "name",
				"type":     "text",
				"required": true,
				"unique":   true,
				"options": map[string]interface{}{
					"min":     3,
					"max":     nil,
					"pattern": "",
				},
			},
			map[string]interface{}{
				"system":   false,
				"id":       "rxdf2ggv",
				"name":     "domainName",
				"type":     "text",
				"required": false,
				"unique":   false,
				"options": map[string]interface{}{
					"min":     nil,
					"max":     nil,
					"pattern": "",
				},
			},
			map[string]interface{}{
				"system":   false,
				"id":       "emrpxodb",
				"name":     "mxaasTypes",
				"type":     "select",
				"required": true,
				"unique":   false,
				"options": map[string]interface{}{
					"maxSelect": 4,
					"values": []interface{}{
						"voice",
						"network",
						"wireless",
						"surveillance",
					},
				},
			},
			map[string]interface{}{
				"system":   false,
				"id":       "wirabiam",
				"name":     "probe",
				"type":     "relation",
				"required": false,
				"unique":   false,
				"options": map[string]interface{}{
					"collectionId":  "le7qimuo80xaq5o",
					"cascadeDelete": false,
					"minSelect":     nil,
					"maxSelect":     nil,
					"displayFields": []interface{}{
						"hostname",
					},
				},
			},
		}
		cols = make([]input.ColumnDescriber, 0, len(schema))
	)
	for _, column := range schema {
		kind := column["type"].(string)
		switch kind {
		case "text":
			name := column["name"].(string)
			baseNode.Id = column["id"].(string)
			baseNode.Key = name
			baseNode.Type = column["type"].(string)
			baseNode.Value = strings.ToTitle(name)
			cols = append(cols, &baseNode)
		case "select":
			name := column["name"].(string)
			selectNode.Id = column["id"].(string)
			selectNode.Key = name
			selectNode.Type = column["type"].(string)
			selectNode.Value = strings.ToTitle(name)
			opts := column["options"].([]string)
			for _, opt := range opts {
				selectNode.Options = append(selectNode.Options, input.OptionParam{
					Key:   opt,
					Value: strings.ToTitle(opt),
				})
			}
			cols = append(cols, &selectNode)
		case "relation":

			name := column["name"].(string)
			withListNode.Id = column["id"].(string)
			withListNode.Key = name
			withListNode.Type = column["type"].(string)
			withListNode.Value = strings.ToTitle(name)
			options := column["Options"].(map[string]any)
			opts := options["displayFields"].([]string)
			for _, opt := range opts {
				withListNode.List.Options = append(withListNode.List.Options, input.OptionParam{
					Key:   opt,
					Value: strings.ToTitle(opt),
				})
			}
			cols = append(cols, &withListNode)
		default:
			continue
		}
	}

	return cols
}

func makeTestData() input.DescriberTable {
	size := 15
	data := []map[string]any{
		{
			"collectionId":   "5uejgttm07xch9i",
			"collectionName": "customers",
			"created":        "2023-03-13 02:27:28.536Z",
			"domainName":     "ams.net",
			"id":             "f8jerx0bo2s8cwq",

			"name":        "AMS.NET",
			"probe":       []string{"qgn9i11kyppke9a"},
			"updated":     "2023-03-13 02:43:22.768Z",
			"mxaasTypes0": []string{"voice", "network"},
		},
		{
			"collectionId":   "5uejgttm07xch9i",
			"collectionName": "customers",
			"created":        "2023-03-13 02:30:22.446Z",
			"domainName":     "lvjusd.org",
			"id":             "yfn03u7d21wfuzx",
			"mxaasTypes":     []string{"voice"},
			"name":           "Livermore Joint Unified School District",
			"probe":          []string{"677bt5fbu800"},
			"updated":        "2023-03-13 02:39:36.158Z",
		},
	}
	inputData := make(input.DescriberTable, 0, size)
	for i, item := range data {
		var (
			inputItem    = make(input.DescriberItem)
			baseNode     = exampleBase
			checkBoxNode = exampleBase
			selectNode   = exampleSelect
			textAreaNode = exampleTextArea
			withListNode = exampleWithList
			baseItem     = &input.BaseItem{
				Columns: inputItem,
			}
		)
		for key, val := range item {
			switch key {
			case "mxaasTypes", "probe":
				selectNode.Key = key
				selectNode.Id = fmt.Sprintf("testing-%s-%s-%d", key, val, i)
				opts := val.([]string)
				for _, opt := range opts {
					selectNode.Options = append(selectNode.Options, input.OptionParam{
						Key:   opt,
						Value: opt,
					})
				}
				inputItem[key] = &selectNode
			case "collectionId", "name":
				baseNode.Key = key
				baseNode.Id = fmt.Sprintf("testing-%s-%s-%d", key, val, i)
				baseNode.Value = val.(string)
				inputItem[key] = &baseNode
			case "collectionName":
				checkBoxNode.Key = key
				checkBoxNode.Id = fmt.Sprintf("testing-%s-%s-%d", key, val.(string), i)
				checkBoxNode.Value = val.(string)
				checkBoxNode.Required = true
				inputItem[key] = &checkBoxNode
			case "created":
				withListNode.Key = key
				withListNode.Id = fmt.Sprintf("testing-%s-%s-%d", key, val.(string), i)
				withListNode.Value = val.(string)
				inputItem[key] = &withListNode
			case "domainName":
				textAreaNode.Key = key
				textAreaNode.Value = val.(string)
				textAreaNode.Id = fmt.Sprintf("testing-%s-%s-%d", key, val.(string), i)
				inputItem[key] = &textAreaNode
			case "id":
				baseNode.Key = key
				baseNode.Id = fmt.Sprintf("testing-%s-%s-%d", key, val, i)
				baseNode.Value = val.(string)
				baseNode.Hidden = true
				baseItem.Id = val.(string)
				inputItem[key] = &baseNode
			}
		}
		inputData = append(inputData, inputItem)
	}
	return inputData
}

func TestNewElementFromAtom(t *testing.T) {
	thing := types.NewNodeFromAtom(atom.Div, html.Attribute{
		Key: atom.Class.String(),
		Val: "row",
	})

	if thing.Data == "div" {
		println(thing.Data)
	}

	anotherTing := types.NewNodeFromAtom(atom.Input, html.Attribute{
		Key: atom.Id.String(),
		Val: "i2345",
	})
	thing.AppendChild(anotherTing)
	buffer := bytes.NewBufferString("")
	err := html.Render(buffer, thing)
	if err != nil {
		t.Error(err)
	}
	str := buffer.String()
	println(str)
}

func TestNewInputsFromDescriber(t *testing.T) {
	var (
		buffer = bytes.NewBufferString("")
		cols   = makeHeader()
	)
	buffer.Reset()

	items, err := input.NewInputsFromDescriber(cols...)
	if err != nil {
		t.Error(err)
	}
	div := types.NewNodeFromAtom(atom.Div)
	for key, rows := range items {
		t.Log(key)
		for _, col := range rows {
			div.AppendChild(col)
		}
	}
	err = html.Render(buffer, div)
	if err != nil {
		t.Error(err)
	}

}

func TestNewFormElem(t *testing.T) {
	var (
		buffer       = bytes.NewBufferString("")
		baseNode     = exampleBase
		checkBoxNode = exampleBase
		selectNode   = exampleSelect
		textAreaNode = exampleTextArea
		withListNode = exampleWithList
		cols         = []input.ColumnDescriber{&checkBoxNode, &baseNode, &selectNode, &textAreaNode, &withListNode}
	)
	checkBoxNode.Key = "checkbox-key"
	selectNode.Key = "select-key"
	textAreaNode.Key = "text_area-kwy"
	withListNode.Key = "datalist-key"
	thing := input.NewFormElem("form-id", "This is a caption", types.NewNodeFromAtom(atom.Div))
	thing.WithSchemaData(cols...)
	thing.BuildForm()

	err := html.Render(buffer, thing.HtmlNode())
	if err != nil {
		t.Error(err)
	}

	str := buffer.String()
	println(str)
}

func TestNewTableElem(t *testing.T) {
	var (
		header = makeHeader()
		data   = makeTestData()
		parent = &html.Node{
			Type:      html.ElementNode,
			DataAtom:  atom.Div,
			Data:      atom.Div.String(),
			Namespace: "",
			Attr:      make([]html.Attribute, 0, 5),
		}
		iTable = input.NewTableElem(uuid.New().String(), "this is a table", parent).UpdateHeaderData(header)
		err    = iTable.LoadData(data)
	)
	if err != nil {
		t.Error(err)
	}

	node := iTable.HtmlNode()
	println(node.Data)
}
