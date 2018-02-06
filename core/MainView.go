package core

type MainView interface {
	Handle(msg uint32, wParam uintptr, lParam uintptr) uintptr
	RootElement() Element
}

var NewDefaultMainView func(windowName string, exStyle uint32, style uint32) MainView
