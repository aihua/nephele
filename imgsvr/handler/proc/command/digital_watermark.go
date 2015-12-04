package command

import (
	log "github.com/Sirupsen/logrus"
	cat "github.com/ctripcorp/cat.go"
	"github.com/ctripcorp/nephele/imgsvr/img4g"
	"math"
	"strconv"
	"strings"
)

type DigitalWatermarkProcessor struct {
	Copyright *img4g.Image
	Cat       cat.Cat
}

func (this *DigitalWatermarkProcessor) Process(img *img4g.Image) error {
	log.Debug("process digitalwatermark")
	var err error = nil

	format, err := img.GetFormat()
	if err != nil {
		return err
	}
	if strings.ToLower(format) != "jpg" && strings.ToLower(format) != "jpeg" {
		info := make(map[string]string)
		info["format"] = format
		logEvent(this.Cat, "DigitalWatermarkRefuse", "NotSupportFormat", info)
		return nil
	}

	width, err := img.GetWidth()
	if err != nil {
		return err
	}
	if width < 256 {
		info := make(map[string]string)
		info["width"] = strconv.Itoa(int(width))
		logEvent(this.Cat, "DigitalWatermarkRefuse", "NotSupportSize", info)
		return nil
	}

	height, err := img.GetHeight()
	if err != nil {
		return err
	}
	if height < 256 {
		info := make(map[string]string)
		info["height"] = strconv.Itoa(int(height))
		logEvent(this.Cat, "DigitalWatermarkRefuse", "NotSupportSize", info)
		return nil
	}

	upr := ((int(math.Min(float64(width), float64(height))) / 100.0) + 1) * 100
	tran := this.Cat.NewTransaction("DigitalWatermark", "Min(width, height)<"+strconv.Itoa(int(upr)))
	tran.AddData("size", "width: "+strconv.Itoa(int(width))+"height: "+strconv.Itoa(int(height)))
	defer func() {
		this.Copyright.DestoryWand()
		tran.SetStatus(err)
		tran.Complete()
	}()
	if err = this.Copyright.CreateWand(); err != nil {
		return err
	}
	err = img.DigitalWatermark(this.Copyright)

	return err
}

func logEvent(cat cat.Cat, title string, name string, data map[string]string) {
	event := cat.NewEvent(title, name)
	if data != nil {
		for k, v := range data {
			event.AddData(k, v)
		}
	}
	event.SetStatus("0")
	event.Complete()
}
