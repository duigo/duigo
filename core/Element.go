package core

type Element struct {
}

var NewElement func(name string) (Element, error)
