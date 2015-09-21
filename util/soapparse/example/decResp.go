package main

import (
	"fmt"
	. "github.com/ctripcorp/nephele/util/soapparse"
	"github.com/ctripcorp/nephele/util/soapparse/response"
)

func main() {
	var header response.Header
	var resp response.SaveResponse
	resp.Process.ProcessResponses = []response.ProcessResponse{
		response.ProcessResponse{
			"1", "PATH", "20,20",
		},
	}
	header.ServerIP = "10.2.25.0"
	content, _ := DecResp(&header, &resp)
	fmt.Println(string(content))
}
