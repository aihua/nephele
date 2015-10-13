package main

import (
	"io/ioutil"

	im4g "../"
	cat "github.com/ctripcorp/nephele/Godeps/_workspace/src/github.com/ctripcorp/cat.go"
)

func main() {
	CAT := cat.Instance()
	i, _ := im4g.NewImageAsPNG(100, 100, CAT)
	ioutil.WriteFile("1.png", i.Blob, 0644)
}
