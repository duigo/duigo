package main

import "github.com/duigo/duigo/core"

func main() {

	app := core.NewApplication()
	m := app.DefineModel("mytable")
	m.Define()
}
