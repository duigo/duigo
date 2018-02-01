package main

import (
	"fmt"
)

//	program 识别
type BraiderProgram struct {
}

func (b *BraiderProgram) Accept(ctx *Context, token *Token) int {
	switch token.Text {
	case "type":
		ctx.Syntaxer.Push(new(BraiderType))
		return REPLAY
	case "model":
		ctx.Syntaxer.Push(new(BraiderModel))
		return REPLAY
	default:
		if TOKEN_EOF == token.Type {
			return EOP
		}

		return ctx.SetFail(fmt.Errorf("001 BraiderProgram: unexpected token '%s' '%d'", token.Text, token.Type))
	}
}

//	type 识别
type BraiderType struct {
	state string
}

func (b *BraiderType) Accept(ctx *Context, token *Token) int {
	switch b.state {
	case "":
		if "type" == token.Text {
			b.state = "wait:name|{"
			return ACCEPT
		}

		return ctx.SetFail(fmt.Errorf("002 BraiderType: unexpected token '%s' '%d'", token.Text, token.Type))
	case "wait:name|{":
		if "{" == token.Text {
			ctx.Syntaxer.Push(new(BraiderTypeBody))
			b.state = "wait:end"
			return REPLAY
		}

		if TOKEN_NAME == token.Type {
			fmt.Println("TYPE-NAME: ", token.Text)
			ctx.Syntaxer.Push(new(BraiderTypeBody))
			b.state = "wait:end"
			return ACCEPT
		}

		return ctx.SetFail(fmt.Errorf("003 BraiderType: unexpected token '%s' '%d'", token.Text, token.Type))

	case "wait:end":
		ctx.Syntaxer.Pop()
		return REPLAY
	default:
		return ctx.SetFail(fmt.Errorf("004 BraiderType: unexpected token '%s' '%d'", token.Text, token.Type))
	}
}

type BraiderTypeBody struct {
	state string
}

func (b *BraiderTypeBody) Accept(ctx *Context, token *Token) int {
	switch b.state {
	case "":
		if "{" == token.Text {
			ctx.Syntaxer.Push(new(BraiderField))
			b.state = "wait:name|}"
			return ACCEPT
		}

		return ctx.SetFail(fmt.Errorf("005 BraiderType: unexpected token '%s' '%d'", token.Text, token.Type))
	case "wait:name|}":
		if "}" == token.Text {
			ctx.Syntaxer.Pop()
			return ACCEPT
		}

		ctx.Syntaxer.Push(new(BraiderField))
		b.state = "wait:name|}"
		return REPLAY
	default:
		return ctx.SetFail(fmt.Errorf("007 BraiderType: unexpected token '%s' '%d'", token.Text, token.Type))
	}
}

type BraiderField struct {
	state string
}

func (b *BraiderField) Accept(ctx *Context, token *Token) int {
	switch b.state {
	case "":
		if TOKEN_NAME == token.Type {
			fmt.Println("FIELD-NAME: ", token.Text)
			b.state = "wait:("
			return ACCEPT
		}

		return ctx.SetFail(fmt.Errorf("008 BraiderField: unexpected token '%s' '%d'", token.Text, token.Type))
	case "wait:(":
		if "(" == token.Text {
			b.state = "wait:)"
			return ACCEPT
		}

		ctx.Syntaxer.Push(new(BraiderTypeRef))
		b.state = "wait:end"
		return REPLAY
	case "wait:)":
		if ")" == token.Text {
			fmt.Println("FIELD-VIRTUAL: ", "true")
			ctx.Syntaxer.Push(new(BraiderTypeRef))
			b.state = "wait:end"
			return ACCEPT
		}

		return ctx.SetFail(fmt.Errorf("010 BraiderField: unexpected token '%s' '%d'", token.Text, token.Type))
	case "wait:end":
		ctx.Syntaxer.Pop()
		return REPLAY
	default:
		return ctx.SetFail(fmt.Errorf("011 BraiderField: unexpected token '%s' '%d'", token.Text, token.Type))
	}
}

type BraiderTypeRef struct {
	state string
}

func (b *BraiderTypeRef) Accept(ctx *Context, token *Token) int {
	switch b.state {
	case "":
		if "type" == token.Text {
			ctx.Syntaxer.Push(new(BraiderType))
			b.state = "wait:("
			return REPLAY
		}

		if TOKEN_NAME == token.Type {
			fmt.Println("REF-NAME: ", token.Text)
			b.state = "wait:("
			return ACCEPT
		}

		return ctx.SetFail(fmt.Errorf("013 BraiderTypeRef: unexpected token '%s' '%d'", token.Text, token.Type))
	case "wait:(":
		if "(" == token.Text {
			ctx.Syntaxer.Push(new(BraiderTypeRef))
			b.state = "wait:,|)"
			return ACCEPT
		}

		ctx.Syntaxer.Pop()
		return REPLAY
	case "wait:,|)":
		if "," == token.Text {
			ctx.Syntaxer.Push(new(BraiderTypeRef))
			b.state = "wait:,|)"
			return ACCEPT
		}

		if ")" == token.Text {
			ctx.Syntaxer.Pop()
			return ACCEPT
		}
		return ctx.SetFail(fmt.Errorf("013 BraiderTypeRef: unexpected token '%s' '%d'", token.Text, token.Type))
	default:
		return ctx.SetFail(fmt.Errorf("014 BraiderTypeRef: unexpected token '%s' '%d'", token.Text, token.Type))
	}
}

type BraiderModel struct {
	state string
}

func (b *BraiderModel) Accept(ctx *Context, token *Token) int {
	switch b.state {
	case "":
		if "model" == token.Text {
			b.state = "wait:name"
			return ACCEPT
		}

		return ctx.SetFail(fmt.Errorf("020 BraiderModel: unexpected token '%s' '%d'", token.Text, token.Type))
	case "wait:name":
		if TOKEN_NAME == token.Type {
			fmt.Println("MODEL: ", token.Text)
			ctx.Syntaxer.Push(new(BraiderTypeRef))
			b.state = "wait:end"
			return ACCEPT
		}
		return ctx.SetFail(fmt.Errorf("021 BraiderModel: unexpected token '%s' '%d'", token.Text, token.Type))
	case "wait:end":
		ctx.Syntaxer.Pop()
		return REPLAY
	default:
		return ctx.SetFail(fmt.Errorf("015 BraiderModel: unexpected token '%s' '%d'", token.Text, token.Type))
	}
}
