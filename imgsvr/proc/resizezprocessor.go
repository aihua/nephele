package proc

import (
	log "github.com/Sirupsen/logrus"
	cat "github.com/ctripcorp/cat.go"
	"github.com/ctripcorp/nephele/imgsvr/img4g"
	"math"
)

type ResizeZProcessor struct {
	Width  int64
	Height int64
	Cat    cat.Cat
}

//高固定，宽（原图比例计算），宽固定，高（原图比例计算） （压缩）
func (this *ResizeZProcessor) Process(img *img4g.Image) error {
	log.Debug("process resize z")
	var err error
	tran := this.Cat.NewTransaction("Command", "ResizeW")
	defer func() {
		tran.SetStatus(err)
		tran.Complete()
	}()

	width, height, err1 := img.Size()
	if err1 != nil {
		err = err1
		return err1
	}

	w, h := this.Width, this.Height
	if w == 0 {
		w = width * h / height
		err = img.Resize(w, h)
		return err
	}
	if h == 0 {
		h = height * w / width
		err = img.Resize(w, h)
		return err
	}

	p1 := float64(this.Width) / float64(this.Height)
	p2 := float64(width) / float64(height)

	if p2 > p1 {
		h = int64(math.Floor(float64(this.Width) / p2))
		if int64(math.Abs(float64(h-this.Height))) < 3 {
			h = this.Height
		}
	} else {
		w = int64(math.Floor(float64(this.Height) * p2))
		if int64(math.Abs(float64(w-this.Width))) < 3 {
			w = this.Width
		}
	}
	err = img.Resize(w, h)
	return err
}
