package core

type MessageHandler interface {
	Handle(msg uint32, wParam uintptr, lParam uintptr) uintptr
}
