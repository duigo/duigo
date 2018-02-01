package main

import (
	"io"
	"bufio"
	"fmt"
)

/*
type  Auth
{
	username	string
	password	string
	exts		map(string,string)
	tags		list(string)
}

type 	UserInfo
{
	id 			string
	height		float
	email		string
	address()	string
}

model	userinfo	UserInfo
*/

var (
	ERROR_EOF = fmt.Errorf("(end)")

	TOKEN_NULL  = 0
	TOKEN_NAME  = 1
	TOKEN_DELIM = 2
	TOKEN_EOF   = 99
)

type Tokenizer struct {
	prev    *Tokenizer     //	前一个词法分析器
	scanner *bufio.Scanner //	扫描器
	line    []byte         //	一行的内容
	lino    int            //	行号
	index   int            //	已经扫描的位置
}

type MoLexer struct {
	top             *Tokenizer
	error           error
	tokenText       string
	tokenType       int
	randomNameIndex int
}

func init() {
	NewLexer = func(r io.Reader) Lexer {
		return &MoLexer{
			top:       &Tokenizer{scanner: bufio.NewScanner(r)},
			tokenType: TOKEN_NULL,
		}
	}
}

func (lex *MoLexer) SetFail(err error) Lexer {
	lex.error = err
	return lex
}

func (lex *MoLexer) OK() bool {
	return lex.error == nil
}

func (lex *MoLexer) Fail() bool {
	return lex.error != nil
}

func (lex *MoLexer) Err() error {
	return lex.error
}

func (lex *MoLexer) Token() *Token {
	return &Token{
		Text: lex.tokenText,
		Type: lex.tokenType,
	}
}

func (lex *MoLexer) Locate() Lexer {
	x := lex.top
	for {
		//	如果一行已经读取完毕
		if x.index >= len(x.line) {
			//	如果读取下一行失败
			if !x.scanner.Scan() {
				//	如果上次兑取失败是由于文件技术导致的
				err := x.scanner.Err()
				if nil != err {
					return lex.SetFail(fmt.Errorf("read failed: %s", err.Error()))
				}

				//	如果没有下层的词法器,不需要再扫描下去了
				if nil == x.prev {
					return lex.SetFail(ERROR_EOF)
				}

				//	当前词法器退栈,回退到前一层次的词法器
				x = x.prev
				continue
			}

			//	成功读取了下一行
			x.line = x.scanner.Bytes()[:]
			x.lino++
			x.index = 0

			//fmt.Println("=== ", string(x.scanner.Bytes()))
		}

		//	尝试跳过本行的空白行或者注释行
		lineSize := len(x.line)
		for i := x.index; i < lineSize; i++ {
			c := x.line[i]
			if 0 != (cm[c] & (CM_SPACE) ) {
				continue
			}

			if '#' == c {
				x.index = lineSize
				break
			}

			//	遇到正常字符开头了
			x.index = i
			break
		}

		//	跳过空白后,如果已经到行尾,说明遇到了空行或者空白行或者注释行
		if x.index >= lineSize {
			continue
		}

		//	终于定位到了
		return lex
	}
}

func (lex *MoLexer) Next() Lexer {
	//	先定位
	if lex.Locate().Fail() {
		if ERROR_EOF == lex.error {
			lex.tokenType = TOKEN_EOF
			lex.tokenText = ""
		}
		return lex
	}

	x := lex.top
	c := x.line[x.index]

	//	字母开头的,识别为名字
	if 0 != (cm[c] & CM_ALPHA) {
		i := x.index + 1
		for (i < len(x.line)) && (0 != (cm[x.line[i]] & (CM_ALPHA | CM_DEC))) {
			i++
		}

		lex.tokenText = string(x.line[x.index:i])
		lex.tokenType = TOKEN_NAME
		x.index = i
		return lex
	}

	//	如果是符号开头的,字节吸收该符号
	if 0 != (cm[c] & CM_TERMINAL) {
		lex.tokenText = string(x.line[x.index:x.index+1])
		lex.tokenType = TOKEN_DELIM
		x.index++
		return lex
	}

	//	如果遇到不能接受的字符
	return lex.SetFail(fmt.Errorf("unsupported char '%c'", c))
}

//
//func (lex *MoLexer) randomName() string {
//	lex.randomNameIndex++
//	return fmt.Sprintf("(unnamed-%d)", lex.randomNameIndex)
//}

//func (lex *MoLexer) AcceptKeyword(keyword []byte) bool {
//	x := lex.top
//
//	//	先定位
//	if !lex.Locate() {
//		return false
//	}
//
//	//	计算几个关键变量缓存
//	keywordSize := len(keyword)
//	remain := x.line[x.index:]
//	remainLen := len(remain)
//
//	//	剩余长度不够
//	if remainLen < keywordSize {
//		lex.error = fmt.Errorf("expect keyword '%s' but got '%s'", string(keyword), string(remain))
//		return false
//	}
//
//	//	逐个字符比较
//	i := 0
//	for ; i < keywordSize; i++ {
//		if keyword[i] != remain[i] {
//			lex.error = fmt.Errorf("except keyword '%s' but got '%s'", string(keyword), remain[0:i])
//			return false
//		}
//	}
//
//	//	完全匹配,但是是行尾
//	if remainLen == keywordSize {
//		x.index += remainLen
//		lex.tokenText = string(remain[:keywordSize])
//		lex.tokenType = TOKEN_NAME
//		return true
//	}
//
//	//	不能是关键字开头的名字
//	if 0 != (cm[remain[i+1]] & (CM_ALPHA | CM_DEC)) {
//		lex.error = fmt.Errorf("except keyword '%s' but got '%s'", string(keyword), remain[0:i+2])
//		return false
//	}
//
//	x.index += keywordSize
//	lex.tokenText = string(remain[:keywordSize])
//	lex.tokenType = TOKEN_NAME
//	return true
//}
//
//func (lex *MoLexer) AcceptName() bool {
//	x := lex.top
//
//	//	先定位
//	if !lex.Locate() {
//		return false
//	}
//
//	//	计算几个关键变量缓存
//	remain := x.line[x.index:]
//	remainLen := len(remain)
//
//	//	名字的首个字符必须是字母
//	if 0 == (cm[remain[0]] & CM_ALPHA) {
//		lex.error = fmt.Errorf("name should be start with alpha but got '%c'", rune(remain[0]))
//		return false
//	}
//
//	//	名字的结束点
//	i := 1
//	for ; i < remainLen; i++ {
//		if 0 == (cm[remain[i]] & (CM_ALPHA | CM_DEC)) {
//			break
//		}
//	}
//
//	x.index += i
//	lex.tokenText = string(remain[:i])
//	lex.tokenType = TOKEN_NAME
//	return true
//}
//
//func (lex *MoLexer) AcceptTerminal(c byte) bool {
//	//	先定位
//	if !  lex.Locate() {
//		return false
//	}
//
//	x := lex.top
//	if c != x.line[0] {
//		lex.error = fmt.Errorf("expect '%c' but got '%c'", c, x.line[0])
//		return false
//	}
//
//	x.index++
//	lex.tokenText = string(x.line[0:0])
//	lex.tokenType = TOKEN_TERMINAL
//	return true
//}
