package soapparse

import (
	"github.com/ctripcorp/nephele/util/soapparse/request"
	"strings"
	"errors"
)

var (
	EmptyString = ""
	EmptyReq = request.Request{}
	RequestTypes []string = []string{
		"SaveImage",
		"DeleteRequest",
		"LoadZipRequest",
		"LoadImgRequest",
	}
)

var (
	ErrIllegalType = errors.New("illegel request type")
)

func GetRequestTypeAndData(content []byte) (string, request.Request, error){
	var req request.Request
	err := EncReq(content, &req)
	if err != nil {
		return EmptyString, EmptyReq, err
	}
	for _, requestType := range RequestTypes {
		if strings.HasSuffix(req.Header.RequestType, requestType) {
			return requestType, req, nil
		}
	}
	return EmptyString, EmptyReq, ErrIllegalType
}
