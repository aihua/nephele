package controllers

import (
	"fmt"
	cat "github.com/ctripcorp/cat.go"
	"github.com/ctripcorp/nephele/imgws/business"
	"github.com/ctripcorp/nephele/util"
	"github.com/ctripcorp/nephele/util/soapparse/request"
	"github.com/ctripcorp/nephele/util/soapparse/response"
	"io/ioutil"
	"net/http"
)

type LogoWS struct{}

func (handler *LogoWS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Cat := cat.Instance()
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/xml")

	tran := Cat.NewTransaction("URL", "LogoWs.asmx")
	util.LogEvent(Cat, "URL", "URL.Client", map[string]string{
		"clientip": util.GetClientIP(r),
		"serverip": util.GetIP(),
		"proto":    r.Proto,
		"referer":  r.Referer(),
		//"agent":    request.UserAgent(),
	})
	var result util.Error
	defer func() {
		if p := recover(); p != nil {
			Cat.LogPanic(p)
		}
		if result.Err != nil && !result.IsNormal {
			tran.SetStatus(result.Err)
		} else {
			tran.SetStatus("0")
		}
		tran.Complete()
	}()

	bts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.LogErrorEvent(Cat, "RequestReadError", err.Error())
		result = util.Error{IsNormal: false, Err: err, Type: "RequestReadError"}

		return
	}
	req := request.Request{}
	result = EncodeRequest(Cat, bts, &req)
	if result.Err != nil && !result.IsNormal {
		//writeResponse(w, []byte(result.Type))
		return
	}

	var (
		resp        interface{}
		header      *response.Header
		e           util.Error
		logoRequest = business.LogoRequest{
			ImageRequest: business.ImageRequest{
				Cat: Cat,
			},
			FontSizes: []int{14, 16, 18, 20},
		}
	)

	requestTran := Cat.NewTransaction("Request", "SaveLogo")
	defer func() {
		if p := recover(); p != nil {
			Cat.LogPanic(p)
		}
		if result.Err != nil && !result.IsNormal {
			requestTran.SetStatus(result.Err)
		} else {
			requestTran.SetStatus("0")
		}
		requestTran.Complete()
	}()
	
	resp, result = logoRequest.Save(&req.SaveRequest)

	if result.Err != nil {
		header = createFailHeader(req.Header, fmt.Sprintf("%v", result.Err))
	} else {
		header = transformHeader(req.Header, RESULTCODE_SUCCESS, "")
	}
	content, e := DecodeResponse(Cat, header, resp)
	if e.Err != nil {
		result = e
		writeResponse(w, []byte(result.Err.(error).Error()))
	} else {
		writeResponse(w, content)
	}

	//content, err := soapparse.DecResp(header, resp) //
	//if err != nil {
	//	return
	//}
	//writeResponse(w, content)
}
