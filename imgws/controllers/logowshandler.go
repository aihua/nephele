package controllers

import (
	"io/ioutil"
	"net/http"

	cat "github.com/ctripcorp/cat.go"

	"github.com/ctripcorp/nephele/imgws/business"
	"github.com/ctripcorp/nephele/util"
	"github.com/ctripcorp/nephele/util/soapparse"
	"github.com/ctripcorp/nephele/util/soapparse/request"
	"github.com/ctripcorp/nephele/util/soapparse/response"
)

type LogoWS struct {}

func (handler *LogoWS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Cat := cat.Instance()
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/xml")
	bts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	req := request.Request{}
	if err := soapparse.EncReq(bts, &req); err != nil {
		return
	}

	var (
		resp interface{}
		header *response.Header
		e util.Error
		logoRequest = business.LogoRequest{
			ImageRequest: business.ImageRequest{
				Cat: Cat,
			},
			FontSizes: []int{14, 16, 18, 20},
		}
	)

	switch req.Header.RequestType {
	case REQUESTTYPE_SAVEIMAGE:
		resp, e = logoRequest.Save(&req.SaveRequest)
	default:
		return
	}
	if e.Err != nil {
	} else {
	}

	content, err := soapparse.DecResp(header, resp) //
	if err != nil {
		return
	}
	writeResponse(w, content)
}
