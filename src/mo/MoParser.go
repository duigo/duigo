package main

import (
	"fmt"
	"reflect"
)

type MoParser struct {
	Braider Braider
	Lexer   Lexer
	Err     error
}

func (p *MoParser) Push(braider Braider) {
	t := reflect.TypeOf(braider)
	fmt.Println("\tPUSH: ", t.String())

	p.Braider = braider
}

func (p *MoParser) Pop() {
	t := reflect.TypeOf(p.Braider)
	fmt.Println("\tPOP:  ", t.String())

	p.Braider = p.Braider.Parent()
}

func (p *MoParser) SetError(err error) int {
	p.Err = err
	return ERROR
}

func init() {
	Parse = func(topLexer Lexer, topBraider Braider, action int) error {
		p := &MoParser{
			Braider: topBraider,
			Lexer:   topLexer,
		}

		p.Lexer = topLexer
		p.Braider = topBraider

		for {
			switch action {
			case NDMORE:
				if p.Lexer.Next().Fail() {
					if ERROR_EOF != p.Lexer.Err() {
						p.Err = p.Lexer.Err()
						return p.Err
					}
				}

				action = p.Braider.Accept(p, p.Lexer.Token())
			case REPLAY:
				action = p.Braider.Accept(p, p.Lexer.Token())
			case ERROR:
				return p.Err
			case EOP:
				return nil
			}
		}
	}
}
