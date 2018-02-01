package main

import "io"

type Token struct {
	Text string
	Type int
}

type Lexer interface {
	Locate() Lexer
	Next() Lexer
	Token() *Token
	SetFail(err error) Lexer

	OK() bool
	Fail() bool
	Err() error
}

var NewLexer func(reader io.Reader) Lexer

const (
	NDMORE = 0
	REPLAY = 1
	ERROR  = 2
	EOP    = 99
)

type Braider interface {
	Parent() Braider
	Accept(parser Parser, token *Token) int
}

type Parser interface {
	Push(braider Braider)
	Pop()
	SetError(error) int
}

var Parse func(lexer Lexer, braider Braider, action int) error
