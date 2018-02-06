package win

import (
	"github.com/lxn/win"
	"github.com/duigo/duigo/core"
	"container/list"
)

const (
	asFunc = 0
	asObj  = 1
)

type hook struct {
	Type int
	Call interface{}
}

type DefaultMessageLoop struct {
	filters *list.List
	idlehds *list.List
	msg     win.MSG
}

func initDefaultMessageLoop() {
	core.NewDefaultMessageLoop = NewDefaultMessageLoop
}

func NewDefaultMessageLoop() core.MessageLoop {
	return &DefaultMessageLoop{
		filters: list.New(),
		idlehds: list.New(),
	}
}

func (loop *DefaultMessageLoop) InsertFilter(f core.Filter) {
	loop.filters.PushBack(&hook{Type: asObj, Call: f})
}

func (loop *DefaultMessageLoop) InsertFilterFunc(f func(msg interface{}) bool) {
	loop.filters.PushBack(&hook{Type: asFunc, Call: f})
}

func (loop *DefaultMessageLoop) RemoveFilter(f interface{}) {
	e := loop.filters.Front()
	for ; nil != e; e = e.Next() {
		hook := e.Value.(*hook)
		if f == hook.Call {
			break
		}
	}

	if nil != e {
		loop.filters.Remove(e)
	}
}

func (loop *DefaultMessageLoop) InsertIdleHandle(f core.IdleHandler) {
	loop.idlehds.PushBack(&hook{Type: asObj, Call: f})
}

func (loop *DefaultMessageLoop) InsertIdleHandleFunc(f func() bool) {
	loop.idlehds.PushBack(&hook{Type: asFunc, Call: f})
}

func (loop *DefaultMessageLoop) RemoveIdleHandle(f interface{}) {
	e := loop.idlehds.Front()
	for ; nil != e; e = e.Next() {
		hook := e.Value.(*hook)
		if hook.Call == f {
			break
		}
	}

	if nil != e {
		loop.idlehds.Remove(e)
	}
}

func (loop *DefaultMessageLoop) Walk(view core.MainView) int {

	winView := view.(*DefaultMainView)

	winView.ShowWindow(win.SW_SHOW)
	winView.UpdateWindow()

	bDoIdle := true

	for {
		for bDoIdle && !win.PeekMessage(&loop.msg, 0, 0, 0, win.PM_NOREMOVE) {
			if !loop.OnIdle() {
				bDoIdle = false
			}
		}

		iRet := win.GetMessage(&loop.msg, 0, 0, 0)
		if iRet < 0 {
			continue // error, don't process
		}

		if 0 == iRet {
			break // WM_QUIT, exit message loop
		}

		if !loop.PreTranslateMessage(&loop.msg) {
			win.TranslateMessage(&loop.msg)
			win.DispatchMessage(&loop.msg)
		}

		//if IsIdleMessage(&loop.msg) {
		//	bDoIdle = true
		//}
	}

	//return loop.msg.WParam
	return 0
}

func (loop *DefaultMessageLoop) OnIdle() bool {
	for e := loop.idlehds.Front(); e != nil; e = e.Next() {
		f := e.Value.(func() bool)
		if nil != f {
			f()
		}
	}

	return false // don't continue
}

func (loop *DefaultMessageLoop) PreTranslateMessage(msg interface{}) bool {
	// loop backwards
	for e := loop.filters.Front(); e != nil; e = e.Next() {
		f := e.Value.(func(msg interface{}) bool)
		if (nil != f) && f(msg) {
			return true
		}
	}

	return false // not translated
}
