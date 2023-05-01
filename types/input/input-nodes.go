package input

import (
	"fmt"
	"github.com/gojrs/go-ui/types"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"strings"
)

type Kinder interface {
	Kind() string
	ColumnKey() string
}

type ItemDescriber interface {
	GetColumns() map[string]ColumnDescriber
	GetId() string
}

type ColumnDescriber interface {
	Kinder
	GetId() string
	GetValue() string
	GetKey() string
	GetPlaceHolder() string
	GetLabel() string
	GetType() string
	GetTitle() string
	GetClasses() []string
	IsHidden() bool
	IsReadonly() bool
	IsRequired() bool
	IsDisabled() bool
}

type BaseItem struct {
	Columns map[string]ColumnDescriber
	Id      string
}

func (bi *BaseItem) GetColumns() map[string]ColumnDescriber { return bi.Columns }
func (bi *BaseItem) GetId() string                          { return bi.Id }

type CollectionDetails struct {
	Id         string            `json:"id"`
	Created    string            `json:"created"`
	Updated    string            `json:"updated"`
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	System     bool              `json:"system"`
	Schema     []ColumnDescriber `json:"schema"`
	Indexes    []string          `json:"indexes"`
	ListRule   any               `json:"listRule"`
	ViewRule   any               `json:"viewRule"`
	CreateRule any               `json:"createRule"`
	UpdateRule any               `json:"updateRule"`
	DeleteRule any               `json:"deleteRule"`
	Options    any               `json:"options"`
}

type BaseInput struct {
	Id          string         `json:"id,omitempty"`
	Value       string         `json:"value,omitempty"`
	Key         string         `json:"key,omitempty"`
	PlaceHolder string         `json:"placeHolder,omitempty"`
	Label       string         `json:"label,omitempty"`
	Type        string         `json:"type,omitempty"`
	Title       string         `json:"title,omitempty"`
	Hidden      bool           `json:"hidden,omitempty"`
	Readonly    bool           `json:"readonly,omitempty"`
	Required    bool           `json:"required,omitempty"`
	Disabled    bool           `json:"disabled,omitempty"`
	System      bool           `json:"system,omitempty"`
	Classes     []string       `json:"classes,omitempty"`
	Class       string         `json:"class,omitempty"`
	Options     map[string]any `json:"options,omitempty"`
}

func (ibo *BaseInput) Kind() string {
	return ibo.Type
}

func (ibo *BaseInput) ColumnKey() string {
	return ibo.Key
}

func (ibo *BaseInput) GetId() string {
	return ibo.Id
}

func (ibo *BaseInput) GetValue() string {
	return ibo.Value
}

func (ibo *BaseInput) GetKey() string {
	return ibo.Key
}

func (ibo *BaseInput) GetPlaceHolder() string {
	return ibo.PlaceHolder
}

func (ibo *BaseInput) GetLabel() string {
	return ibo.Label
}

func (ibo *BaseInput) GetType() string {
	return ibo.Type
}

func (ibo *BaseInput) GetTitle() string {
	return ibo.Title
}

func (ibo *BaseInput) IsHidden() bool {
	return ibo.Hidden
}

func (ibo *BaseInput) IsReadonly() bool {
	return ibo.Readonly
}

func (ibo *BaseInput) IsRequired() bool {
	return ibo.Readonly
}

func (ibo *BaseInput) IsDisabled() bool {
	return ibo.Disabled
}

func (ibo *BaseInput) GetClasses() []string {
	return ibo.Classes
}

type DataList struct {
	Id      string
	Options []OptionParam
}

type SelectInput struct {
	BaseInput
	Multiple bool          `json:"multiple,omitempty"`
	Options  []OptionParam `json:"options,omitempty"`
}

type TextAreaInput struct {
	BaseInput
	Cols      uint `json:"cols,omitempty"`
	Rows      uint `json:"Rows,omitempty"` // per html: cols default == 20 && Rows == 2  we should match
	Maxlength int  `json:"maxlength,omitempty"`
	Minlength int  `json:"Minlength,omitempty"`
}

type RelationInput struct {
	BaseInput
	Options RelationOptions `json:"options"`
}

type RelationOptions struct {
	CollectionId  string      `json:"collectionId"`
	CascadeDelete bool        `json:"cascadeDelete"`
	MinSelect     interface{} `json:"minSelect"`
	MaxSelect     interface{} `json:"maxSelect"`
	DisplayFields []string    `json:"displayFields"`
}

type JsonInput struct {
	BaseInput
	Options RelationOptions `json:"options"`
}

type OptionParam struct {
	Key   string
	Value string
}

type WithListOptions struct {
	BaseInput
	List DataList
}

func NewInputsFromDescriber(inputs ...ColumnDescriber) (nodes map[string][]*html.Node, err error) {
	nodes = make(map[string][]*html.Node)
	for _, kind := range inputs {

		switch option := kind.(type) {
		case *SelectInput:
			nodes[kind.ColumnKey()] = NewSelectNode(option)
		case *WithListOptions:
			nodes[kind.ColumnKey()] = NewInputWithListNode(option)
		case *BaseInput:
			nodes[kind.ColumnKey()] = NewInputNode(option)
		case *TextAreaInput:
			nodes[kind.ColumnKey()] = NewTextAreaNode(option)
		default:
			err = fmt.Errorf("unknown kind")
			return nodes, err
		}

	}

	return nodes, err
}

func NewInputNode(param *BaseInput) []*html.Node {
	var (
		input = types.NewElementFromAtom(atom.Input, defaultInputAttr(param)...)
		label = &html.Node{
			Type:     html.ElementNode,
			DataAtom: atom.Label,
			Data:     atom.Label.String(),
			Attr: []html.Attribute{
				{
					Key: atom.For.String(),
					Val: param.Id,
				},
			},
		}
	)

	label.AppendChild(&html.Node{
		Type:     html.TextNode,
		DataAtom: 0,
		Data:     param.Label,
	})

	return []*html.Node{label, input}
}

func NewTextAreaNode(param *TextAreaInput) []*html.Node {
	input := types.NewElementFromAtom(atom.Textarea, defaultInputAttr(param)...)
	label := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Label,
		Data:     atom.Label.String(),
		Attr: []html.Attribute{
			{
				Key: atom.For.String(),
				Val: param.Id,
			},
		},
	}
	lt := &html.Node{
		Type:     html.TextNode,
		DataAtom: 0,
		Data:     param.Label,
	}
	txt := &html.Node{
		Type:     html.TextNode,
		DataAtom: 0,
		Data:     "",
	}

	if param.Required {
		input.Attr = append(input.Attr, html.Attribute{
			Key: atom.Required.String(),
			Val: atom.Required.String(),
		})
	}
	if param.Hidden {
		input.Attr = append(input.Attr, html.Attribute{
			Key: atom.Hidden.String(),
			Val: atom.Hidden.String(),
		})
	}
	if param.Readonly {
		input.Attr = append(input.Attr, html.Attribute{
			Key: atom.Readonly.String(),
			Val: atom.Readonly.String(),
		})
	}
	if param.Disabled {
		input.Attr = append(input.Attr, html.Attribute{
			Key: atom.Disabled.String(),
			Val: atom.Disabled.String(),
		})
	}
	label.AppendChild(lt)
	input.AppendChild(txt)
	return []*html.Node{label, input}
}

func NewInputWithListNode(param *WithListOptions) []*html.Node {
	var label, input, datalist html.Node
	var (
		kind   = html.ElementNode
		idKey  = atom.Id
		forKey = atom.For
	)

	input.Type = kind
	input.DataAtom = atom.Input
	input.Data = atom.Input.String()
	input.Attr = defaultInputAttr(param)

	label.DataAtom = atom.Label
	label.Data = atom.Label.String()
	label.AppendChild(&html.Node{
		Type:     html.TextNode,
		DataAtom: 0,
		Data:     param.Label,
	})
	label.Type = kind
	label.Attr = []html.Attribute{
		{
			Key: forKey.String(),
			Val: param.Id,
		},
	}

	datalist.Type = kind
	datalist.DataAtom = atom.Datalist
	datalist.Data = atom.Datalist.String()
	datalist.Attr = []html.Attribute{
		{
			Key: idKey.String(),
			Val: param.List.Id,
		},
	}

	for _, option := range param.List.Options {
		opt := &html.Node{
			Type:     kind,
			DataAtom: atom.Option,
			Data:     atom.Option.String(),
			Attr: []html.Attribute{
				{
					Key: "value",
					Val: option.Key,
				},
			},
		}
		opt.AppendChild(&html.Node{
			Type:     html.TextNode,
			DataAtom: 0,
			Data:     option.Value,
		})
		datalist.AppendChild(opt)
	}

	return []*html.Node{&label, &input, &datalist}
}

func NewSelectNode(param *SelectInput) []*html.Node {
	var label *html.Node
	var (
		sel  = types.NewElementFromAtom(atom.Select, defaultInputAttr(param)...)
		aFor = atom.For
	)

	label = types.NewElementFromAtom(atom.Label, []html.Attribute{
		{
			Key: aFor.String(),
			Val: param.Id,
		}}...)

	label.AppendChild(&html.Node{
		Type:     html.TextNode,
		DataAtom: 0,
		Data:     param.Label,
	})

	for _, option := range param.Options {
		opt := types.NewElementFromAtom(atom.Select, []html.Attribute{
			{
				Key: atom.Value.String(),
				Val: option.Key,
			}}...)

		opt.AppendChild(&html.Node{
			Type:     html.TextNode,
			DataAtom: 0,
			Data:     option.Value,
		})
		sel.AppendChild(opt)
	}

	return []*html.Node{label, sel}
}

func defaultInputAttr(param ColumnDescriber) (attrs []html.Attribute) {

	attrs = []html.Attribute{
		{
			Key: atom.Id.String(),
			Val: param.GetId(),
		},
		{
			Key: atom.Class.String(),
			Val: strings.Join(param.GetClasses(), " "),
		},
		{
			Key: atom.Value.String(),
			Val: param.GetValue(),
		},
		{
			Key: atom.Placeholder.String(),
			Val: param.GetPlaceHolder(),
		},
		{
			Key: atom.Title.String(),
			Val: param.GetTitle(),
		},
		{
			Key: fmt.Sprintf("data-col-key"),
			Val: param.GetKey(),
		},
	}
	switch param.GetType() {
	case atom.Select.String(), "relation", "json":

	default:
		attrs = append(attrs, html.Attribute{
			Key: atom.Type.String(),
			Val: param.GetType(),
		})
	}
	if param.IsHidden() {
		attrs = append(attrs, html.Attribute{
			Key: atom.Hidden.String(),
			Val: atom.Hidden.String(),
		})
	}
	if param.IsDisabled() {
		attrs = append(attrs, html.Attribute{
			Key: atom.Disabled.String(),
			Val: atom.Disabled.String(),
		})
	}
	if param.IsReadonly() {
		attrs = append(attrs, html.Attribute{
			Key: atom.Readonly.String(),
			Val: atom.Readonly.String(),
		})
	}
	if param.IsRequired() {
		attrs = append(attrs, html.Attribute{
			Key: atom.Required.String(),
			Val: atom.Required.String(),
		})
	}
	return attrs
}
