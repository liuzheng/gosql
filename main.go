package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
)

var (
	migrate        = flag.Bool("migrate", false, "migrate, which is responsible for applying and unapplying migrations.")
	makemigrations = flag.String("makemigrations", "", "makemigrations, which is responsible for creating new migrations based on the changes you have made to your models.")
)

type Colume struct {
	Name       string
	Tag        string
	Mysql      string
	PrimaryKey bool
}
type Table struct {
	Name   string
	Colume []Colume
}

func main() {
	flag.Parse()
	if *makemigrations != "" {
		_, err := os.Stat(*makemigrations)
		if os.IsNotExist(err) {
			fmt.Println("Nothing to do")
			return
		}
		fmt.Println("makemigrations")
		for _, gofile := range GetAllFiles(*makemigrations) {
			db := ReadGOFile(gofile)
			fmt.Println(db)
			if len(db) == 0 {
				//fmt.Println(gofile)
			}
		}

	} else if *migrate {
		fmt.Println("migrate")

	} else {
		fmt.Println("Nothing to do")
	}
}

//获取指定目录下的所有文件,包含子目录下的文件
func GetAllFiles(dirPath string) (files []string) {
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	//PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix)

	for _, fi := range dir {
		if fi.IsDir() {
			files = append(files, GetAllFiles(dirPath+fi.Name())...)
		} else {
			ok := strings.HasSuffix(fi.Name(), ".go")
			if ok {
				//fmt.Println(fi.Name())
				files = append(files, dirPath+fi.Name())
			}
		}
	}

	return files
}
func ReadGOFile(gofile string) (DB []Table) {
	aA_Compile, _ := regexp.Compile("([a-z])([A-Z])")

	content, err := ioutil.ReadFile(gofile)
	if err != nil {
		log.Fatal(err)
	}

	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, gofile, content, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	for _, Decls := range f.Decls {
		GenDecl, ok := Decls.(*ast.GenDecl)
		if ok && GenDecl.Tok == token.TYPE {
			for _, Specs := range GenDecl.Specs {
				//fmt.Println(Specs.(*ast.TypeSpec).Name.Name)
				table := Table{
					Name:   Specs.(*ast.TypeSpec).Name.Name,
					Colume: []Colume{},
				}
				for _, list := range Specs.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List {
					//ast.Print(fset, list)
					if list.Tag != nil {
						//ast.Print(fset, list)
						colume := Colume{
							Name: strings.ToLower(aA_Compile.ReplaceAllString(list.Names[0].Name, "${1}_${2}")),
						}

						colume.Tag = strings.Replace(list.Tag.Value, "`", "", -1)
						if tag := reflect.StructTag(colume.Tag).Get("mysql"); tag != "" {
							// e.g.: 	id         uint16    `mysql:"SMALLINT,NOT_NULL,AUTO_INCREMENT,PRIMARY_KEY"`
							// Todo: define the sql with golang type not in the tag, like django. e.g.: if `id` is uint16,
							//  so the will be `smallint` in mysql, no need to add it in tag
							if strings.Contains(tag, "PRIMARY_KEY") {
								colume.PrimaryKey = true
								tag = strings.Replace(tag, ",PRIMARY_KEY", "", -1)
							}
							tag = strings.Replace(tag, "_", " ", -1)
							colume.Mysql = strings.Replace(tag, ",", " ", -1)
						}
						// Todo: not only mysql
						table.Colume = append(table.Colume, colume)
						//fmt.Println(reflect.StructTag(strings.Replace(list.Tag.Value, "`", "", -1)).Get("mysql"))
						//fmt.Println(list.Names[0].Name)
						//fmt.Println(list.Tag.Value)
					}
				}
				if len(table.Colume) > 0 {
					DB = append(DB, table)
				}
				//ast.Print(fset, Specs)
			}
		}
	}
	return
}
