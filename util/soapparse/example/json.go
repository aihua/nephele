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
	_, req, err := GetRequestTypeAndData(content)
	bytes, _ := req.SaveRequest.Process.MarshalJSON()
	fmt.Println(string(bytes))
}
