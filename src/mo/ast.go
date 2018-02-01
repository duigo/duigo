package main

import "fmt"

type Field struct {
	Name    string //	字段名称
	Virtual bool   //	是否虚拟字段
	Type    *Type  //	类型名
}

type Type struct {
	Name   string  //	类型名称
	Params []*Type //	类型参数
	Fields []Field //	字段列表
}

type Model struct {
	Name string //	模型名称
	Type *Type  //	类型引用
}

type MoApp struct {
	Types  map[string]Type
	Models map[string]Model
}



func PrintApp(app *MoApp) {
	if nil != app {
		for name, value := range app.Types {
			fmt.Println("--- ", name)
			PrintType(&value)
		}

		for name, value := range app.Models {
			fmt.Println("--- ", name)
			PrintModel(&value)
		}
	}
}

func PrintType(typo *Type) {
	//	说明是基于模板的类型
	if len(typo.Params) > 0 {
		fmt.Printf("%s", typo.Name)
	} else {
		fmt.Printf("type %s\n", typo.Name)
		fmt.Printf("{\n")
		for _, f := range typo.Fields {
			fmt.Printf("\t%s", f.Name)
			if f.Virtual {
				fmt.Printf("()")
			}
			fmt.Printf("\t%s\n", f.Type.Name)
		}
		fmt.Printf("}\n")
	}
}

func PrintModel(model *Model) {
	fmt.Printf("model\t%s\t%s\n", model.Name, model.Type.Name)
}
