package main

import (
	"io/ioutil"

	im4g "../"
	cat "github.com/ctripcorp/nephele/Godeps/_workspace/src/github.com/ctripcorp/cat.go"
)

func main() {
	CAT := cat.Instance()
	i, _ := im4g.NewImageAsPNG(100, 100, CAT)
	i.CreateWand()
	i.SetFormat("PNG")
	err := i.AnnotateImage("hello")
	if err != nil {
		println(err.Error())
	}
	ioutil.WriteFile("1.png", i.Blob, 0644)
	i.DestoryWand()
}
