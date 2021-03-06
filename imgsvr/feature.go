package imgsvr

import (
	"errors"
	cat "github.com/ctripcorp/cat.go"
	"github.com/ctripcorp/nephele/imgsvr/data"
	"github.com/ctripcorp/nephele/imgsvr/proc"
	"strconv"
	"strings"
)

type feature interface {
	Process() (proc.ImageProcessor, bool, error)
}

type tgresizefeature struct {
	width      int64
	height     int64
	resizetype string
	Cat        cat.Cat
}

func (this *tgresizefeature) Process() (proc.ImageProcessor, bool, error) {
	var W = "w"
	if this.resizetype == W && this.height != 0 {
		return nil, false, errors.New("channel: tg, resizetype: W, reason: height can only be '0' ")
	}

	if this.resizetype != W && this.height == 0 {
		return nil, false, errors.New("channel: tg, resizetype: " + this.resizetype + ", reason: height can't be '0' ")
	}
	if this.width == 10000 && this.height < 10000 {
		return &proc.ResizeWProcessor{this.width, this.height, this.Cat}, false, nil
	}
	if this.width < 10000 && this.height == 10000 {
		return &proc.ResizeWProcessor{this.width, this.height, this.Cat}, false, nil
	}
	if this.resizetype == W && this.height == 0 {
		return &proc.ResizeZProcessor{Width: this.width, Height: 100000, Cat: this.Cat}, false, nil
	}
	return &proc.ResizeRProcessor{this.width, this.height, this.Cat}, false, nil
}

type hotelresizefeature struct {
	width      int64
	height     int64
	resizetype string
	Cat        cat.Cat
}

func (this *hotelresizefeature) Process() (proc.ImageProcessor, bool, error) {
	if this.resizetype == "r" || this.resizetype == "c" {
		return &proc.ResizeRProcessor{this.width, this.height, this.Cat}, false, nil
	} else {
		return nil, true, nil
	}
}

type hotelrotatefeature struct {
	rotate float64
}

func (this *hotelrotatefeature) Process() (proc.ImageProcessor, bool, error) {
	rotateStr, err := data.GetRotates(Hotel)
	if err != nil {
		return nil, true, err
	}

	var checkstr = JoinString(",", strconv.FormatFloat(this.rotate, 'f', -1, 64), ",")
	if strings.Contains(rotateStr, checkstr) {
		return nil, true, nil
	}
	opacitiesStr, err := data.GetDissolves(Hotel)
	if err != nil {
		return nil, true, err
	}
	if strings.Contains(opacitiesStr, checkstr) {
		return nil, false, nil
	}
	return nil, false, errors.New(JoinString("channel: hotel, reason: not support rotate degree", strconv.FormatFloat(this.rotate, 'f', -1, 64)))
}
