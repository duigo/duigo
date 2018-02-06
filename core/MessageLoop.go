package core

type Filter interface {
	PreTranslateMessage(msg interface{}) bool
}

type IdleHandler interface {
	OnIdle() bool
}

type MessageLoop interface {
	InsertFilter(f Filter)
	InsertFilterFunc(f func(msg interface{}) bool)
	RemoveFilter(f interface{})

	InsertIdleHandle(f IdleHandler)
	InsertIdleHandleFunc(f func() bool)
	RemoveIdleHandle(f interface{})

	Walk(view MainView) int
}

var NewDefaultMessageLoop func() MessageLoop
