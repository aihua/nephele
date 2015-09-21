package main

import (
	"fmt"
	"io/ioutil"
	. "github.com/ctripcorp/nephele/util/soapparse"
)

func main() {
	content, err := ioutil.ReadFile("source.xml")
	if err != nil {
	}
	t, req, _ := GetRequestTypeAndData(content)
	fmt.Println(t)
	fmt.Println(req)
}
