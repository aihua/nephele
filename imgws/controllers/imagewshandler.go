package controllers

import (
	"errors"
	"fmt"
	cat "github.com/ctripcorp/cat.go"
	"github.com/ctripcorp/nephele/imgws/business"
	"github.com/ctripcorp/nephele/util"
	"github.com/ctripcorp/nephele/util/soapparse"
	"github.com/ctripcorp/nephele/util/soapparse/request"
	"github.com/ctripcorp/nephele/util/soapparse/response"
	"io/ioutil"
	"net/http"
	"strings"
)

type ImageWS struct{}

var (
	RESULTCODE_SUCCESS      = "Success"
	RESULTCODE_FALI         = "Fail"
	REQUESTTYPE_SAVEIMAGE   = "Arch.Base.ImageWS.SaveImage"
	REQUESTTYPE_DELETEIMAGE = "Arch.Base.ImageWS.DeleteImage"
	REQUESTTYPE_LOADZIP     = "Arch.Base.ImageWS.LoadZip"
	REQUESTTYPE_LOADIMAGE   = "Arch.Base.ImageWS.LoadImage"
	CATTRANSUCCESS          = "0"
)

func (handler *ImageWS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/xml")

	Cat := cat.Instance()
	tran := Cat.NewTransaction("URL", "ImageWs.asmx")
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
		msg := []byte(err.Error())
		writeResponse(w, msg)
		return
	}
	req := request.Request{}
	result = EncodeRequest(Cat, bts, &req)
	if result.Err != nil && !result.IsNormal {
		writeResponse(w, []byte(result.Type))
		return
	}
	var (
		resp       interface{}
		header     *response.Header
		e          util.Error
		imgRequest = business.ImageRequest{Cat: Cat}
	)

	requestTran := Cat.NewTransaction("Request", strings.Replace(req.Header.RequestType, "Arch.Base.ImageWS.", "", -1))
	defer func() {
		if p := recover(); p != nil {
			Cat.LogPanic(p)
		}
		if result.Err != nil && !result.IsNormal {
			requestTran.SetStatus(result.Err)
		} else {
			requestTran.SetStatus(CATTRANSUCCESS)
		}
		requestTran.Complete()
	}()

	switch req.Header.RequestType {
	case REQUESTTYPE_SAVEIMAGE:
		resp, e = imgRequest.Save(&req.SaveRequest)
	case REQUESTTYPE_LOADIMAGE:
		resp, e = imgRequest.Download(&req.LoadImgRequest)
	case REQUESTTYPE_LOADZIP:
		resp, e = imgRequest.DownloadZip(&req.LoadZipRequest)
	case REQUESTTYPE_DELETEIMAGE:
		resp, e = imgRequest.Delete(&req.DeleteRequest)
	default:
		util.LogErrorEvent(Cat, "RequestTypeInvalid", req.Header.RequestType)
		e = util.Error{IsNormal: true, Err: errors.New("requesttype is invalid!"), Type: "RequestTypeInvalid"}
	}
	if e.Err != nil {
		result = e
		header = createFailHeader(req.Header, fmt.Sprintf("%v", e.Err))
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
}

func DecodeResponse(Cat cat.Cat, header *response.Header, resp interface{}) ([]byte, util.Error) {
	soapTran := Cat.NewTransaction("SOAP", "DecodeResponse")
	result := util.Error{}
	defer func() {
		if result.Err != nil && !result.IsNormal {
			soapTran.SetStatus(result.Err)
		} else {
			soapTran.SetStatus(CATTRANSUCCESS)
		}
		soapTran.Complete()
	}()
	content, err := soapparse.DecResp(header, resp) //
	if err != nil {
		util.LogErrorEvent(Cat, "SoapParseResponseError", err.Error())
		result = util.Error{IsNormal: false, Err: err, Type: "SoapParseResponseError"}
	}
	return content, result
}
func EncodeRequest(Cat cat.Cat, bts []byte, req *request.Request) util.Error {
	soapTran := Cat.NewTransaction("SOAP", "EncodeRequest")
	result := util.Error{}
	defer func() {
		if result.Err != nil && !result.IsNormal {
			soapTran.SetStatus(result.Err)
		} else {
			soapTran.SetStatus(CATTRANSUCCESS)
		}
		soapTran.Complete()
	}()
	if err := soapparse.EncReq(bts, req); err != nil {
		util.LogErrorEvent(Cat, "SoapParseRequestError", err.Error())
		result = util.Error{IsNormal: false, Err: err, Type: "SoapParseRequestError"}
	}
	return result
}

func writeResponse(w http.ResponseWriter, msg []byte) {
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(msg)))
	w.Write(msg)
}

func createFailHeader(r request.Header, resultmsg string) *response.Header {
	return transformHeader(r, RESULTCODE_FALI, resultmsg)
}

func transformHeader(r request.Header, resultcode string, resultmsg string) *response.Header {
	header := response.Header{}
	header.AssemblyVersion = r.AssemblyVersion
	header.Environment = "fws" //todo
	header.RequestBodySize = r.RequestBodySize
	header.RequestID = r.RequestID
	header.ResultCode = resultcode
	header.RouteStep = r.RouteStep
	header.SerializeMode = r.SerializeMode
	header.ServerIP = util.GetIP()
	header.ShouldRecordPerformanceTime = false //todo
	header.UserID = r.UserID
	header.ResultMsg = resultmsg
	//TODO
	return &header
}
