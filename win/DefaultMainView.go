package win

import (
	"github.com/lxn/win"
	"unsafe"
	"syscall"
	"fmt"
	"github.com/duigo/duigo/core"
)

type DefaultMainView struct {
	HWnd       win.HWND
	OldWndProc uintptr
}

const (
	DefaultMainViewClassName = "duigo.MainWindow"
)

var (
	RegisteredWindowClasses = make(map[string]*win.WNDCLASSEX)

	defaultMainViewProc = syscall.NewCallback(DefaultMainViewProc)
)

func initDefaultMainView() {
	core.NewDefaultMainView = NewDefaultMainView
}

func DefaultMainViewProc(hWnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	var mainView *DefaultMainView

	//	如果是创建窗口的第一个消息,获取mainView的方式比较特别,需要专门设置HWND的golang对象的指针
	if msg == win.WM_NCCREATE {
		mainView = (*DefaultMainView)(unsafe.Pointer((*win.CREATESTRUCT)(unsafe.Pointer(lParam)).CreateParams))

		//	设置本地窗口的钩子
		mainView.HWnd = hWnd
		win.SetWindowLongPtr(hWnd, win.GWLP_USERDATA, (uintptr)(unsafe.Pointer(mainView)))
	} else {
		mainView = (*DefaultMainView)(unsafe.Pointer(win.GetWindowLongPtr(hWnd, win.GWLP_USERDATA)))
	}

	if nil != mainView {
		//	将消息委托给窗口的HandleMessage进行处理
		lResult := mainView.Handle(msg, wParam, lParam)

		//	对于窗口过程的最后一个消息,执行特殊处理
		if (msg == win.WM_NCDESTROY) && (mainView != nil) {
			win.SetWindowLongPtr(mainView.HWnd, win.GWLP_USERDATA, 0)
			mainView.HWnd = 0
		}

		return lResult
	}

	return win.DefWindowProc(hWnd, msg, wParam, lParam)
}

func NewDefaultMainView(windowName string, exStyle uint32, style uint32) core.MainView {

	//	注册
	err := RegisterDefaultMainViewClass(DefaultMainViewClassName, defaultMainViewProc)
	if nil != err {
		panic(fmt.Errorf("register window class for the default-main-view failed: %s", err.Error()))
	}

	//	创建一个视图
	view := &DefaultMainView{}

	//	获取主窗口的窗口类名
	u16ClassName, _ := syscall.UTF16PtrFromString(DefaultMainViewClassName)

	//	获取主窗口的窗口名
	u16WindowName, _ := syscall.UTF16PtrFromString(windowName)

	//	创建一个窗口
	wnd := win.CreateWindowEx(
		exStyle,
		u16ClassName,
		u16WindowName,
		style|win.WS_CLIPSIBLINGS,
		win.CW_USEDEFAULT,
		win.CW_USEDEFAULT,
		win.CW_USEDEFAULT,
		win.CW_USEDEFAULT,
		0,
		0,
		0,
		unsafe.Pointer(view))

	if 0 == wnd {
		panic(fmt.Errorf("create main-view failed"))
	}

	return view
}

func (v *DefaultMainView) RootElement() core.Element {
	return nil
}

func (v *DefaultMainView) ShowWindow(nCmdShow int32) {
	win.ShowWindow(v.HWnd, nCmdShow)
}

func (v *DefaultMainView) UpdateWindow() {
	win.UpdateWindow(v.HWnd)
}

func (v *DefaultMainView) Handle(msg uint32, wParam uintptr, lParam uintptr) uintptr {
	return win.DefWindowProc(v.HWnd, msg, wParam, lParam)
}

// MustRegisterWindowClass registers the specified window class.
func RegisterDefaultMainViewClass(className string, wndProcPtr uintptr) error {
	if _, ok := RegisteredWindowClasses[className]; ok {
		return fmt.Errorf("window class '%s' already registered", className)
	}

	hInst := win.GetModuleHandle(nil)
	if hInst == 0 {
		return fmt.Errorf("retrieve the module handle of this process failed")
	}

	hIcon := win.LoadIcon(0, win.MAKEINTRESOURCE(win.IDI_APPLICATION))
	if hIcon == 0 {
		return fmt.Errorf("load icon of '%d' failed", win.IDI_APPLICATION)
	}

	hCursor := win.LoadCursor(0, win.MAKEINTRESOURCE(win.IDC_ARROW))
	if hCursor == 0 {
		return fmt.Errorf("load cursor of '%d' failed", win.IDC_ARROW)
	}

	uclassName, err := syscall.UTF16PtrFromString(className)
	if nil != err {
		return fmt.Errorf("convert string to utf-16 format failed")
	}

	//	生成注册窗口的对象
	var wc win.WNDCLASSEX
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.LpfnWndProc = wndProcPtr
	wc.HInstance = hInst
	wc.HIcon = hIcon
	wc.HCursor = hCursor
	wc.HbrBackground = win.COLOR_WINDOW
	wc.LpszClassName = uclassName
	wc.Style = 0
	if atom := win.RegisterClassEx(&wc); atom == 0 {
		return fmt.Errorf("register window class '%s' failed", className)
	}

	RegisteredWindowClasses[className] = &wc
	return nil
}
