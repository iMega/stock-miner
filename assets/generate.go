// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

var fs http.FileSystem = http.Dir("./build")

//go run -tags=dev assets/generate.go
func main() {
	err := vfsgen.Generate(fs, vfsgen.Options{
		PackageName:  "main",
		BuildTags:    "!dev",
		VariableName: "Assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
