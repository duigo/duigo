package core

type Model interface {
	Define(name string, getter func() interface{}, setter func(interface{}))
	Depict(name string, v string)
	Exists(name string) bool

	Get(name string) interface{}
	Set(name string, def interface{})

	Watch(name string, watcher Watcher)
}

type Watcher interface {
	Update(name string, v interface{})
}
