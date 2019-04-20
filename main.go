package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
)

var (
	migrate        *bool = flag.Bool("migrate", false, "migrate, which is responsible for applying and unapplying migrations.")
	makemigrations *bool = flag.Bool("makemigrations", false, "makemigrations, which is responsible for creating new migrations based on the changes you have made to your models.")
)

func main() {
	flag.Parse()
	if *makemigrations {
		fmt.Println("makemigrations")
		content, err := ioutil.ReadFile("main.go")
		if err != nil {
			log.Fatal(err)
		}

		fset := token.NewFileSet()

		f, err := parser.ParseFile(fset, "main.go", content, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		for _, Decls := range f.Decls {
			GenDecl, ok := Decls.(*ast.GenDecl)
			if ok && GenDecl.Tok == token.TYPE {
				for _, Specs := range GenDecl.Specs {
					fmt.Println(Specs.(*ast.TypeSpec).Name)
					for _,list := range Specs.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List{
						fmt.Println(list.Tag.Value)
					}
					//ast.Print(fset, Specs)
				}
			}
		}
	} else if *migrate {
		fmt.Println("migrate")

	} else {
		fmt.Println("Nothing to do")
	}
}
