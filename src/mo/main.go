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

	result := NewResult()

	context := NewContext(result, NewLexer(result, bytes.NewReader(s)), NewSyntaxer(result, new(BraiderProgram))).Walk(ACCEPT)
	if context.Result.Fail() {
		fmt.Println(result.Error())
	}
}
