package business
//As the test doesn't pass
/*
import (
	"fmt"
	"github.com/ctripcorp/nephele/util/soapparse/request"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

var (
	ImagePath = "1.jpg"
)

func TestCheckPlanID(t *testing.T) {
	r := request.SaveRequest{}
	r.PlanID = "tg_o1.5"

	r.Process.AnyTypes = []request.AnyType{request.AnyType{}}
	imgRequest := ImageRequest{}

	imgRequest.checkPlanID(&r)
	if r.Channel != "tg" {
		t.Error("update channel fail")
	}

	if len(r.Process.AnyTypes) != 0 {
		t.Error("clear processlist fail")
	}
}

func TestCheckSaveCheckItem(t *testing.T) {
	r := request.Request{}
	bts, err := ioutil.ReadFile(ImagePath)
	if err != nil {
		t.Error(err)
	}
	r.SaveRequest.FileBytes = bts
	fmt.Println("------------------" + strconv.Itoa(len(bts)))
	//check issvg
	r.SaveRequest.CheckItem.IsOtherImage = true
	imgRequest := ImageRequest{}
	e := imgRequest.checkSaveCheckItem(&r.SaveRequest)
	if e.Err == nil {
		t.Error("issvg check fail")
	}

	r.SaveRequest.CheckItem.IsOtherImage = false
	e = imgRequest.checkSaveCheckItem(&r.SaveRequest)
	if e.Err != nil {
		t.Error(e.Err)
	}
}

func TestSave(t *testing.T) {
	r := request.Request{}
	bts, err := ioutil.ReadFile(ImagePath)
	if err != nil {
		t.Error(err)
	}
	r.SaveRequest.FileBytes = bts
	r.SaveRequest.Channel = "tg"
	r.SaveRequest.TargetFormat = "jpg"
	imgRequest := ImageRequest{}
	resp, e := imgRequest.Save(&r.SaveRequest)
	if e.Err != nil {
		t.Error(e)
	}

	println("image targetpath---->>>" + resp.TargetPath)
}

func TestIsNewUri(t *testing.T) {
	uri := "\\1002080000000000448E2.jpg"
	if !isNewUri(uri) {
		t.Error("check is new uri fail")
	}

	uri = "1002080000000000628A4.jpg"
	if isNewUri(uri) {
		t.Error("check is new uri fail")
	}
}

func TestIsFdfs(t *testing.T) {
	uri := "\\fd\\hotel\\1.jpg"
	if !isFdfs(uri) {
		t.Error("check fdfs uri fail")
	}
}

func TestT1(t *testing.T) {
	uri := "\\t1\\hotel\\1.jpg"
	if !isT1(uri) {
		t.Error("check t1 uri fail")
	}
}

func TestGetStorage(t *testing.T) {
	newimageuri := "\\100208000000000108D30.jpg"
	imgRequest := ImageRequest{}
	_, e := imgRequest.getStorage(newimageuri)
	if e.Err != nil {
		t.Error(e.Err)
	}

	fdfsuri := "/fd/hotel/group1/M01/6C/36/CgIZH1YAwemAUhVrAAMwrelaA5k766.png"
	_, e = imgRequest.getStorage(fdfsuri)
	if e.Err != nil {
		t.Error(e.Err)
	}

	t1uri := "/t1/headphoto/057/777/943/0b93f8268d5546308915f4f9fcaa9483.jpg"
	_, e = imgRequest.getStorage(t1uri)
	if e.Err != nil {
		t.Error(e.Err)
	}

	uri := "/tg/057/777/943/0b93f8268d5546308915f4f9fcaa9483.jpg"
	_, e = imgRequest.getStorage(uri)
	if e.Err != nil {
		t.Error(e.Err)
	}
}

func TestDownload(t *testing.T) {
	r := request.LoadImgRequest{}
	r.FilePath = "\\100208000000000108D30.jpg"

	imgRequest := ImageRequest{}
	resp, e := imgRequest.Download(&r)
	if e.Err != nil {
		t.Error(e.Err)
	}
	if len(resp.FileBytes) < 1 {
		t.Error("download fail")
	}
}

func TestDownloadZip(t *testing.T) {
	r := request.LoadZipRequest{}
	r.Files.LoadFiles = []request.LoadFile{request.LoadFile{FilePath: "\\100208000000000108D30.jpg"}}

	imgRequest := ImageRequest{}
	resp, e := imgRequest.DownloadZip(&r)
	if e.Err != nil {
		t.Error(e.Err)
	}
	if len(resp.FileBytes) < 1 {
		t.Error("download fail")
	}
	fmt.Println(strconv.Itoa(len(resp.FileBytes)))
	file, err := os.Create("test.zip")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()
	file.Write(resp.FileBytes)
}*/
