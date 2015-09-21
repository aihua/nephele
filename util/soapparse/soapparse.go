package soapparse

import (
	"encoding/xml"
	"errors"
	"strings"
	"github.com/ctripcorp/nephele/util/soapparse/request"
	"github.com/ctripcorp/nephele/util/soapparse/response"
)

var (
	ReqPrefix    = []byte(`<?xml version="1.0" encoding="utf-8"?><soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema"><soap:Body><Request xmlns="http://tempuri.org/"><requestXML>`)
	ReqSuffix    = []byte(`</requestXML></Request></soap:Body></soap:Envelope>`)
	RespPrefix   = []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?><soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\"><soap:Body><RequestResponse xmlns=\"http://tempuri.org/\"><RequestResult><?xml version=\"1.0\"?><Response>")
	RespSuffix   = []byte("</Response></RequestResult></RequestResponse></soap:Body></soap:Envelope>")
	ReqPrefixLen = len(ReqPrefix)
	ReqSuffixLen = len(ReqSuffix)
)

var (
	ErrCorruptedSoapStream = errors.New("corrupted soap stream")
)

func EncReq(content []byte, req *request.Request) (err error) {
	l := len(content)
	if l < ReqPrefixLen+ReqSuffixLen {
		return ErrCorruptedSoapStream
	}
	content = content[ReqPrefixLen : l-ReqSuffixLen-2]
	s := string(content)
	s = strings.Replace(s, "&lt;", "<", -1)
	s = strings.Replace(s, "&gt;", ">", -1)
	return xml.Unmarshal([]byte(s), &req)
}

func DecResp(header *response.Header, resp interface{}) ([]byte, error) {
	var content []byte
	headerContent, err := xml.MarshalIndent(&header, "", "\r")
	respContent, err := xml.MarshalIndent(&resp, "", "\r")
	content = append(RespPrefix, headerContent...)
	content = append(content, respContent...)
	content = append(content, RespSuffix...)
	return content, err
}
