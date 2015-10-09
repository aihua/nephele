package business

import (
	cat "github.com/ctripcorp/cat.go"

	"github.com/ctripcorp/nephele/util"
	"github.com/ctripcorp/nephele/util/soapparse/request"
	"github.com/ctripcorp/nephele/util/soapparse/response"
)

type LogoRequest struct {
	Cat cat.Cat
}


func (this LogoRequest) Save(r *request.SaveRequest) (response.SaveResponse, util.Error) {
	return response.SaveResponse{}, util.Error{}
}
