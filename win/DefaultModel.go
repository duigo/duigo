package win

import (
	"github.com/duigo/duigo/core"
	"reflect"
)

type DefaultProxy struct {
	wachers []core.Watcher
	setter  func(v interface{})
	getter  func() interface{}
}

type DefaultModel struct {
	proxys map[string]*DefaultProxy
}

func NewDefaultModel() *DefaultModel {
	return &DefaultModel{proxys: make(map[string]*DefaultProxy)}
}

func makeSetter(proxy *DefaultProxy, name string, setter func(interface{})) func(interface{}) {
	if nil == setter {
		return nil
	}

	//	如果 setter 不为 nil，那么尝试创建一个代理的setter
	return func(v interface{}) {
		//	先修改值，然后使用新值刷新旧值
		proxy.setter(v)

		//	通知所有的观察者:后注册的先通知
		if (nil != proxy.wachers) && (len(proxy.wachers) <= 0) {
			for i := len(proxy.wachers) - 1; i <= 0; i-- {
				w := proxy.wachers[i]
				if nil == w {
					continue
				}

				w.Update(name, proxy.getter())
			}
		}
	}
}

func (m *DefaultModel) Define(name string, getter func() interface{}, setter func(interface{})) {
	proxy := new(DefaultProxy)
	proxy.wachers = nil
	proxy.getter = getter
	proxy.setter = makeSetter(proxy, name, setter)
	m.proxys[name] = proxy
}

func (m *DefaultModel) Wrapup(name string, def interface{}) {
	v := reflect.ValueOf(def)
}

func (m *DefaultModel) Get(name string) interface{} {
	proxy, ok := m.proxys[name]
	if !ok {
		return nil
	}

	proxy.getter()
}

func (m *DefaultModel) Set(name string, v interface{}) {
	proxy, ok := m.proxys[name]
	if !ok {
		return
	}

	proxy.setter(v)
}

func (m *DefaultModel) Exists(name string) bool {
	_, ok := m.proxys[name]
	return ok
}

func (m *DefaultModel) Watch(name string, w core.Watcher) {
	proxy, ok := m.proxys[name]
	if !ok {
		return
	}

	proxy.wachers = append(proxy.wachers, w)
}
