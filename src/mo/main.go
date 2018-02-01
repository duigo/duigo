package main

import (
	"io/ioutil"
	"fmt"
	"bytes"
)

func main() {
	////创建一个 app
	//app := MoApp{
	//	Types:  make(map[string]Type),
	//	Models: make(map[string]Model),
	//}
	//
	//s, err := ioutil.ReadFile("test.mo")
	//if nil != err {
	//	fmt.Println(err.error())
	//	return
	//}
	//
	//lex := OpenLexer(bytes.NewReader(s))
	//for lex.Next() {
	//	fmt.Println(string(lex.tokenText))
	//}
	//fmt.Println(lex.error.error())

	s, err := ioutil.ReadFile("test.mo")
	if nil != err {
		fmt.Printf("read file error: %s", err.Error())
		return
	}

	lex := NewLexer(bytes.NewReader(s))
	//for lex.Next().OK() {
	//	fmt.Println("--- ", lex.Token().Type, " ", lex.Token().Text)
	//}
	//
	//if TOKEN_EOF != lex.Token().Type {
	//	if err := lex.Err(); nil != err {
	//		fmt.Printf("parse error: %s\n", err.Error())
	//		return
	//	}
	//}

	if err := Parse(lex, new(BraiderProgram).init(nil), NDMORE); nil != err {
		fmt.Println(err.Error())
	}
}
