package core

type Filter interface {
	PreTranslateMessage(typo int, msg interface{}) bool
}
