package core

import "reflect"

type Element interface {
	Name() string

	Parent() Element

	FirstChild(name string) Element
	LastChild(name string) Element

	ChildCount() int

	Prev(name string) Element
	Next(name string) Element

	Attributes() Attributes
	Styles() Styles
	Status() Status
}

type Attribute struct {
	Name  string      //	属性名
	Value interface{} //	属性值
	Unit  string      //	属性单位
}

//	属性表
type Attributes interface {
	Set(name string) Attribute
	Get(name string) (error, Attribute)
}

//	样式表
type Styles interface {
}

//	状态表
type Status interface {
}

var ControlClasses map[string]reflect.Type

var NewElement func(name string) (Element, error)
