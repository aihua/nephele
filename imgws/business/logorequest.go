package business

import (
	"fmt"
	"github.com/ctripcorp/nephele/util"
	"github.com/ctripcorp/nephele/util/soapparse/request"
	"github.com/ctripcorp/nephele/util/soapparse/response"
	"github.com/ctripcorp/nephele/imgws/models"
)

type LogoRequest struct {
	ImageRequest
	FontSizes []int
}


func (this LogoRequest) Save(r *request.SaveRequest) (response.SaveResponse, util.Error) {
	s := r.SourcePath
	li := models.LogoInfo{}
	li.UnmarshalJSON([]byte(s))

	storage, e := this.getStorage(li.Path)
	if e.Err != nil {
		return response.SaveResponse{}, e
	}

	for _, size := range this.FontSizes {
		bts := getTextImage(li.Name, size)
		storage.UploadSlave(bts, fmt.Sprintf("_logo_%d", size), "PNG")
	}
	return response.SaveResponse{}, util.Error{}
}
