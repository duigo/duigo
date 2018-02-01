package core

type View interface {
	RootElement() Element
}

var NewView func() View
