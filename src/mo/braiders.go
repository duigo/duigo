package main

import (
	"fmt"
)

//	基础类
type BraiderBasic struct {
	parent Braider
}

func (b *BraiderBasic) Parent() Braider {
	return b.parent
}

//	program 识别
type BraiderProgram struct {
	BraiderBasic
}

func (b *BraiderProgram) init(parent Braider) *BraiderProgram {
	b.parent = parent
	return b
}

func (b *BraiderProgram) Accept(parser Parser, token *Token) int {
	switch token.Text {
	case "type":
		parser.Push(new(BraiderType).init(b))
		return REPLAY
	case "model":
		parser.Push(new(BraiderModel).init(b))
		return REPLAY
	default:
		if TOKEN_EOF == token.Type {
			return EOP
		}

		return parser.SetError(fmt.Errorf("001 BraiderProgram: unexpected token '%s' '%d'", token.Text, token.Type))
	}
}

//	type 识别
type BraiderType struct {
	BraiderBasic
	state string
}

func (b *BraiderType) init(parent Braider) *BraiderType {
	b.parent = parent
	b.state = "wait:type"
	return b
}

func (b *BraiderType) Accept(parser Parser, token *Token) int {
	switch b.state {
	case "wait:type":
		if "type" == token.Text {
			b.state = "wait:name|{"
			return NDMORE
		}

		return parser.SetError(fmt.Errorf("002 BraiderType: unexpected token '%s' '%d'", token.Text, token.Type))
	case "wait:name|{":
		if "{" == token.Text {
			parser.Push(new(BraiderTypeBody).init(b))
			b.state = "wait:end"
			return REPLAY
		}

		if TOKEN_NAME == token.Type {
			fmt.Println("TYPE-NAME: ", token.Text)
			parser.Push(new(BraiderTypeBody).init(b))
			b.state = "wait:end"
			return NDMORE
		}

		return parser.SetError(fmt.Errorf("003 BraiderType: unexpected token '%s' '%d'", token.Text, token.Type))

	case "wait:end":
		parser.Pop()
		return REPLAY
	default:
		return parser.SetError(fmt.Errorf("004 BraiderType: unexpected token '%s' '%d'", token.Text, token.Type))
	}
}

type BraiderTypeBody struct {
	BraiderBasic
	state string
}

func (b *BraiderTypeBody) init(parent Braider) *BraiderTypeBody {
	b.parent = parent
	b.state = "wait:{"
	return b
}

func (b *BraiderTypeBody) Accept(parser Parser, token *Token) int {
	switch b.state {
	case "wait:{":
		if "{" == token.Text {
			parser.Push(new(BraiderField).init(b))
			b.state = "wait:name|}"
			return NDMORE
		}

		return parser.SetError(fmt.Errorf("005 BraiderType: unexpected token '%s' '%d'", token.Text, token.Type))
	case "wait:name|}":
		if "}" == token.Text {
			parser.Pop()
			return NDMORE
		}

		parser.Push(new(BraiderField).init(b))
		b.state = "wait:name|}"
		return REPLAY
	default:
		return parser.SetError(fmt.Errorf("007 BraiderType: unexpected token '%s' '%d'", token.Text, token.Type))
	}
}

type BraiderField struct {
	BraiderBasic
	state string
}

func (b *BraiderField) init(parent Braider) *BraiderField {
	b.parent = parent
	b.state = "wait:name"
	return b
}

func (b *BraiderField) Accept(parser Parser, token *Token) int {
	switch b.state {
	case "wait:name":
		if TOKEN_NAME == token.Type {
			fmt.Println("FIELD-NAME: ", token.Text)
			b.state = "wait:("
			return NDMORE
		}

		return parser.SetError(fmt.Errorf("008 BraiderField: unexpected token '%s' '%d'", token.Text, token.Type))
	case "wait:(":
		if "(" == token.Text {
			b.state = "wait:)"
			return NDMORE
		}

		parser.Push(new(BraiderTypeRef).init(b))
		b.state = "wait:end"
		return REPLAY
	case "wait:)":
		if ")" == token.Text {
			fmt.Println("FIELD-VIRTUAL: ", "true")
			parser.Push(new(BraiderTypeRef).init(b))
			b.state = "wait:end"
			return NDMORE
		}

		return parser.SetError(fmt.Errorf("010 BraiderField: unexpected token '%s' '%d'", token.Text, token.Type))
	case "wait:end":
		parser.Pop()
		return REPLAY
	default:
		return parser.SetError(fmt.Errorf("011 BraiderField: unexpected token '%s' '%d'", token.Text, token.Type))
	}
}

type BraiderTypeRef struct {
	BraiderBasic
	state string
}

func (b *BraiderTypeRef) init(parent Braider) *BraiderTypeRef {
	b.parent = parent
	b.state = "wait:name|type"
	return b
}

func (b *BraiderTypeRef) Accept(parser Parser, token *Token) int {
	switch b.state {
	case "wait:name|type":
		if "type" == token.Text {
			parser.Push(new(BraiderType).init(b))
			b.state = "wait:("
			return REPLAY
		}

		if TOKEN_NAME == token.Type {
			fmt.Println("REF-NAME: ", token.Text)
			b.state = "wait:("
			return NDMORE
		}

		return parser.SetError(fmt.Errorf("013 BraiderTypeRef: unexpected token '%s' '%d'", token.Text, token.Type))
	case "wait:(":
		if "(" == token.Text {
			parser.Push(new(BraiderTypeRef).init(b))
			b.state = "wait:,|)"
			return NDMORE
		}

		parser.Pop()
		return REPLAY
	case "wait:,|)":
		if "," == token.Text {
			parser.Push(new(BraiderTypeRef).init(b))
			b.state = "wait:,|)"
			return NDMORE
		}

		if ")" == token.Text {
			parser.Pop()
			return NDMORE
		}
		return parser.SetError(fmt.Errorf("013 BraiderTypeRef: unexpected token '%s' '%d'", token.Text, token.Type))
	default:
		return parser.SetError(fmt.Errorf("014 BraiderTypeRef: unexpected token '%s' '%d'", token.Text, token.Type))
	}
}

type BraiderModel struct {
	BraiderBasic
	state string
}

func (b *BraiderModel) init(parent Braider) *BraiderModel {
	b.parent = parent
	b.state = "wait:model"
	return b
}

func (b *BraiderModel) Accept(parser Parser, token *Token) int {
	switch b.state {
	case "wait:model":
		if "model" == token.Text {
			b.state = "wait:name"
			return NDMORE
		}

		return parser.SetError(fmt.Errorf("020 BraiderModel: unexpected token '%s' '%d'", token.Text, token.Type))
	case "wait:name":
		if TOKEN_NAME == token.Type {
			fmt.Println("MODEL: ", token.Text)
			parser.Push(new(BraiderTypeRef).init(b))
			b.state = "wait:end"
			return NDMORE
		}
		return parser.SetError(fmt.Errorf("021 BraiderModel: unexpected token '%s' '%d'", token.Text, token.Type))
	case "wait:end":
		parser.Pop()
		return REPLAY
	default:
		return parser.SetError(fmt.Errorf("015 BraiderModel: unexpected token '%s' '%d'", token.Text, token.Type))
	}
}
