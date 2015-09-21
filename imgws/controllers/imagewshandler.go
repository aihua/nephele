package controllers

import (
	"errors"
	"fmt"
	"github.com/ctripcorp/nephele/imgws/business"
	"github.com/ctripcorp/nephele/util"
	"github.com/ctripcorp/nephele/util/soapparse"
	"github.com/ctripcorp/nephele/util/soapparse/request"
	"github.com/ctripcorp/nephele/util/soapparse/response"
	"io/ioutil"
	"net/http"
)

type ImageWS struct{}

var (
	RESULTCODE_SUCCESS      = "Sucess"
	RESULTCODE_FALI         = "Fail"
	REQUESTTYPE_SAVEIMAGE   = "Arch.Base.ImageWS.SaveImage"
	REQUESTTYPE_DELETEIMAGE = "Arch.Base.ImageWS.DeleteImage"
	REQUESTTYPE_LOADZIP     = "Arch.Base.ImageWS.LoadZip"
	REQUESTTYPE_LOADIMAGE   = "Arch.Base.ImageWS.LoadImage"
)

func (handler *ImageWS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")

	bts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := []byte(err.Error())
		writeResponse(w, msg)
		return
	}

	req := request.Request{}
	if err := soapparse.EncReq(bts, &req); err != nil {
		content := []byte("request invalid")
		writeResponse(w, content)
		return
	}
	var (
		resp       interface{}
		header     *response.Header
		e          util.Error
		imgRequest = business.ImageRequest{}
	)

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
		e = util.Error{IsNormal: true, Err: errors.New("requesttype is invalid!"), Type: "RequestTypeInvalid"}
	}
	if e.Err != nil {
		header = createFailHeader(req.Header, fmt.Sprintf("%v", e.Err))
	} else {
		header = transformHeader(req.Header, RESULTCODE_SUCCESS, "")
	}
	content, err := soapparse.DecResp(header, resp)
	if err != nil {
		msg := []byte(err.Error())
		writeResponse(w, msg)
		return
	}
	writeResponse(w, content)
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
	//TODO
	return &header
}
