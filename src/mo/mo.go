package main

import (
	"io"
	"reflect"
	"fmt"
)

type Token struct {
	Text string
	Type int
}

type Result interface {
	SetFail(err error)
	Clear()

	OK() bool
	Fail() bool
	Err() error

	Error() string
}

var NewResult func() Result

type Lexer interface {
	Locate() Lexer
	Next() Lexer
	Token() *Token
	SetFail(err error) Lexer
	Result() Result
	EOF() bool
}

var NewLexer func(r Result, reader io.Reader) Lexer

const (
	ACCEPT = 0
	REPLAY = 1
	ERROR  = 2
	EOP    = 9
)

type Braider interface {
	Accept(ctx *Context, token *Token) int
}

type Syntaxer interface {
	Push(braider Braider) Braider
	Pop() Braider
	Top() Braider
}

var NewSyntaxer func(r Result, topBraider Braider) Syntaxer

type Context struct {
	Syntaxer Syntaxer
	Lexer    Lexer
	Result   Result
}

func (ctx *Context) Walk(action int) *Context {
	for {
		switch action {
		case ACCEPT:
			if ctx.Lexer.Next().Result().Fail() {
				if !ctx.Lexer.EOF() {
					return ctx
				}
			}
			action = ctx.Syntaxer.Top().Accept(ctx, ctx.Lexer.Token())
		case REPLAY:
			action = ctx.Syntaxer.Top().Accept(ctx, ctx.Lexer.Token())
		case ERROR:
			return ctx
		case EOP:
			return ctx
		}
	}
}

func (ctx *Context) SetFail(err error) int {
	ctx.Result.SetFail(err)
	return ERROR
}

var NewContext func(r Result, lexer Lexer, syntaxer Syntaxer) *Context

type defaultSyntaxer struct {
	braiders []Braider
}

func (p *defaultSyntaxer) Push(braider Braider) Braider {
	t := reflect.TypeOf(braider)
	fmt.Println("\tPUSH:br ", t.String())

	p.braiders = append(p.braiders, braider)
	return braider
}

func (p *defaultSyntaxer) Pop() Braider {
	size := len(p.braiders)
	if size > 0 {
		oldTop := p.braiders[size-1]

		t := reflect.TypeOf(oldTop)
		fmt.Println("\tPOP:  ", t.String())

		p.braiders = p.braiders[:size-1]
		return oldTop
	}

	return nil
}

func (p *defaultSyntaxer) Top() Braider {
	size := len(p.braiders)
	if size > 0 {
		return p.braiders[size-1]
	}

	return nil
}

type defaultResult struct {
	err error
}

func (r *defaultResult) Clear() {
	r.err = nil
}

func (r *defaultResult) SetFail(err error) {
	r.err = err
}

func (r *defaultResult) OK() bool {
	return nil == r.err
}

func (r *defaultResult) Fail() bool {
	return nil != r.err
}

func (r *defaultResult) Err() error {
	return r.err
}

func (r *defaultResult) Error() string {
	if nil == r.err {
		return ""
	}

	return r.err.Error()
}

func init() {

	NewSyntaxer = func(r Result, topBraider Braider) Syntaxer {
		return &defaultSyntaxer{
			braiders: append(make([]Braider, 0, 5), topBraider),
		}
	}

	NewResult = func() Result {
		return &defaultResult{}
	}

	NewContext = func(r Result, lexer Lexer, syntaxer Syntaxer) *Context {
		return &Context{
			Lexer:    lexer,
			Syntaxer: syntaxer,
			Result:   r,
		}
	}
}
