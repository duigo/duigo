package core

type Application interface {

	InsertFilter(typo int, filter Filter)
	RemoveFilter(typo int, filter Filter)

	Run(view View) int
}

var NewApplication func() Application
