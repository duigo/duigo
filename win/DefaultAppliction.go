package win

import (
	"github.com/duigo/duigo/core"
)

type DefaultApplication struct {
}

func initDefaultApplication() {
	core.NewApplication = NewApplication
}

func NewApplication() core.Application {
	return new(DefaultApplication)
}

//
//func (m *DefaultApplication) GetInstance() interface{} {
//	return nil
//}
//
//func (m *DefaultApplication) Run() int {
//	var bDoIdle bool = true
//	var nIdleCount int = 0
//
//	for {
//		for bDoIdle && !win.PeekMessage(&(m.msg), 0, 0, 0, win.PM_NOREMOVE) {
//			if !m.OnIdle(nIdleCount) {
//				bDoIdle = false
//			}
//			nIdleCount++
//		}
//
//		bRet := win.GetMessage(&m.msg, 0, 0, 0);
//		if bRet == -1 {
//			continue // error, don't process
//		}
//
//		if 0 == bRet {
//			break // WM_QUIT, exit message loop
//		}
//
//		if (!m.PreTranslateMessage(&m.msg)) {
//			win.TranslateMessage(&m.msg);
//			win.DispatchMessage(&m.msg);
//		}
//
//		if (IsIdleMessage(&m.msg)) {
//			bDoIdle = true;
//			nIdleCount = 0;
//		}
//	}
//
//	return (int)(unsafe.Pointer(m.msg.WParam))
//}
//
//func IsIdleMessage(msg *win.MSG) bool {
//	// These messages should NOT cause idle processing
//	switch msg.Message {
//	case win.WM_MOUSEMOVE:
//	case win.WM_PAINT:
//	case 0x0118: // WM_SYSTIMER (caret blink)
//		return false
//	}
//
//	return true
//}
//
//func (m *DefaultApplication) PreTranslateMessage(pMsg *win.MSG) bool {
//	// loop backwards
//	for i := len(m.MessageFilters) - 1; i >= 0; i-- {
//		pMessageFilter := m.MessageFilters[i]
//		if (nil != pMessageFilter) && (pMessageFilter.PreTranslateMessage(pMsg)) {
//			return true
//		}
//	}
//	return false // not translated
//}
//
//// override to change idle processing
//func (m *DefaultApplication) OnIdle(int /*nIdleCount*/) bool {
//	for i := 0; i < len(m.IdleHandlers); i++ {
//		pIdleHandler := m.IdleHandlers[i]
//		if pIdleHandler != nil {
//			pIdleHandler.OnIdle()
//		}
//	}
//	return false // don't continue
//}
