package main

import (
	"fmt"
	. "github.com/ctripcorp/nephele/util/soapparse"
	"github.com/ctripcorp/nephele/util/soapparse/request"
	"io/ioutil"
)

func main() {
	content, err := ioutil.ReadFile("source.xml")
	if err != nil {
	}
	var req request.Request
	err = EncReq(content, &req)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(req)
}
