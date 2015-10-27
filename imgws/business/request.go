package business

import (
	"bytes"
	"errors"
	cat "github.com/ctripcorp/cat.go"
	"github.com/ctripcorp/nephele/imgws/models"
	"github.com/ctripcorp/nephele/util"
	"github.com/ctripcorp/nephele/util/soapparse"
	"github.com/ctripcorp/nephele/util/soapparse/request"
	"github.com/ctripcorp/nephele/util/soapparse/response"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"path"
	"strconv"
	"strings"
	"time"
)

var (
	ERRORTYPE_MARSHALJSON               = "MarshalJsonErr"
	ERRORTYPE_STORAGETYPENOSUPPORTE     = "StorageTypeNoSupporte"
	SVG                                 = 6063
	NEWIMAGENAMELENGTH                  = 21
	TableCount                      int = 6 // 64
)

type ImageRequest struct {
	Cat cat.Cat
}

func (this ImageRequest) Save(r *request.SaveRequest) (response.SaveResponse, util.Error) {
	r.TargetFormat = strings.ToLower(r.TargetFormat)
	r.Channel = strings.ToLower(r.Channel)
	util.LogEvent(this.Cat, "Save", r.Channel, nil)

	this.checkPlanID(r)
	if err := this.checkSaveRequest(r); err.Err != nil {
		return response.SaveResponse{}, err
	}
	if err := this.checkSaveCheckItem(r); err.Err != nil {
		return response.SaveResponse{}, err
	}
	storage, storageType := NewStorage(this.Cat)
	path, e := storage.Upload(r.FileBytes, r.TargetFormat)
	if e.Err != nil {
		return response.SaveResponse{}, e
	}
	tableZone := sharding()
	channel := models.Channel{Cat: this.Cat}
	imgIndex := models.ImageIndex{ChannelCode: channel.GetChannelCode(r.Channel), StoragePath: path, StorageType: storageType, TableZone: tableZone, Cat: this.Cat}
	plan := ""
	if r.Process.AnyTypes != nil && len(r.Process.AnyTypes) > 0 {
		bts, err := r.Process.MarshalJSON()
		if err != nil {
			util.LogErrorEvent(this.Cat, ERRORTYPE_MARSHALJSON, err.Error())
			return response.SaveResponse{}, util.Error{IsNormal: false, Err: err, Type: ERRORTYPE_MARSHALJSON}
		}
		plan = string(bts)
	}

	if e := imgIndex.SaveToDB(plan); e.Err != nil {
		return response.SaveResponse{}, e
	}
	uri := imgIndex.GetImageName()
	return response.SaveResponse{CheckPass: true, OriginalPath: uri, TargetPath: uri}, util.Error{}
}

var shardingCount int = rand.New(rand.NewSource(time.Now().UnixNano())).Intn(TableCount)

func sharding() int {
	if shardingCount == 99999999 {
		shardingCount = 0
	}
	shardingCount = shardingCount + 1
	return shardingCount%TableCount + 1
}

var jpg string = "jpg"

func (this ImageRequest) checkPlanID(r *request.SaveRequest) {
	if len(r.PlanID) > 0 {
		if r.PlanID == "tg_o1.5" || r.PlanID == "tg_n2.0" {
			r.Channel = "tg"
		}
		if r.PlanID == "tg_test1" {
			r.Channel = "test"
		}
		r.OriginalFormat = jpg
		r.TargetFormat = jpg
		r.Process.AnyTypes = make([]request.AnyType, 0)
	}
}
func (this ImageRequest) checkSaveRequest(r *request.SaveRequest) util.Error {
	t := "ParamInvalid"
	if r.Channel == "" {
		util.LogEvent(this.Cat, t, "ChannelEmpty", nil)
		return util.Error{IsNormal: true, Err: errors.New("Channel can't be empty"), Type: t}
	}
	if r.FileBytes == nil || len(r.FileBytes) < 1 {
		util.LogEvent(this.Cat, t, util.JoinString(r.Channel, "-FileBytesEmpty"), nil)
		return util.Error{IsNormal: true, Err: errors.New("FileBytes can't be empty"), Type: t}
	}
	channel := models.Channel{Cat: this.Cat}
	channels, err := channel.GetAll()
	if err.Err != nil {
		return err
	}
	_, exists := channels[r.Channel]
	if !exists {
		util.LogEvent(this.Cat, t, util.JoinString(r.Channel, "-NoRegister"), nil)
		return util.Error{IsNormal: true, Err: errors.New(util.JoinString("channel[", r.Channel, "] isn't register!")), Type: t}
	}
	return util.Error{}
}

func (this ImageRequest) checkSaveCheckItem(r *request.SaveRequest) util.Error {
	t := "CheckFail"
	//if r.CheckItem == nil {
	//	return util.Error{}
	//}
	if r.CheckItem.IsOtherImage {
		if !isSvg(r.FileBytes) {
			util.LogEvent(this.Cat, t, "FormatInvalid", map[string]string{"detail": "image isn't svg!"})
			return util.Error{IsNormal: true, Err: errors.New("image isn't svg!"), Type: t}
		}
	} else {
		img, _, err := image.Decode(bytes.NewReader(r.FileBytes))
		if err != nil {
			util.LogEvent(this.Cat, t, "FormatInvalid", map[string]string{"detail": err.Error()})
			return util.Error{IsNormal: true, Err: err, Type: t}
		}
		//todo check img format
		if r.CheckItem.MinWidth > 0 && r.CheckItem.MinWidth > img.Bounds().Dx() {
			util.LogEvent(this.Cat, t, "LessMinWidth", map[string]string{"detail": util.JoinString("MinWidth:"+strconv.Itoa(r.CheckItem.MinWidth), " ImageWidth:", strconv.Itoa(img.Bounds().Dx()))})
			return util.Error{IsNormal: true, Err: errors.New("image width is less minwidth!"), Type: t}
		}
		if r.CheckItem.MinHeight > 0 && r.CheckItem.MinHeight > img.Bounds().Dy() {
			util.LogEvent(this.Cat, t, "LessMinHeight", map[string]string{"detail": util.JoinString("MinHeight:"+strconv.Itoa(r.CheckItem.MinHeight), " ImageHeight:", strconv.Itoa(img.Bounds().Dy()))})
			return util.Error{IsNormal: true, Err: errors.New("image heigth is less minheight!"), Type: t}
		}
	}
	if r.CheckItem.MaxBytes > 0 && int(r.CheckItem.MaxBytes) < len(r.FileBytes) {
		util.LogEvent(this.Cat, t, "BeyondMaxSize", map[string]string{"detail": util.JoinString("MaxSize:"+strconv.Itoa(int(r.CheckItem.MaxBytes)), " ImageSize:", strconv.Itoa(len(r.FileBytes)))})
		return util.Error{IsNormal: true, Err: errors.New("image size beyond max size"), Type: t}
	}
	return util.Error{Err: nil, IsNormal: true}
}

func isSvg(bts []byte) bool {
	i, _ := strconv.Atoi(strconv.Itoa(int(bts[0])) + strconv.Itoa(int(bts[1])))
	return i == SVG
}

func (this ImageRequest) Download(r *request.LoadImgRequest) (response.LoadImgResponse, util.Error) {
	storage, e := this.getStorageBySource(r.FilePath, r.IsSource)
	if e.Err != nil {
		return response.LoadImgResponse{}, e
	}

	util.LogEvent(this.Cat, "Download", this.GetChannel(r.FilePath), map[string]string{"uri": r.FilePath})

	bts, e := storage.Download()
	if e.Err != nil {
		return response.LoadImgResponse{}, e
	}
	bts = []byte(soapparse.B64.EncodeToString(bts))
	return response.LoadImgResponse{FileBytes: bts}, util.Error{}
}

func (this ImageRequest) DownloadZip(r *request.LoadZipRequest) (response.LoadZipResponse, util.Error) {
	t := "ParamInvalid"
	if len(r.Files.LoadFiles) < 1 {
		util.LogEvent(this.Cat, t, "NoLoadFiles", nil)
		return response.LoadZipResponse{}, util.Error{IsNormal: true, Err: errors.New("No files in request"), Type: "NoFilesInRequest"}
	}
	files := make(map[string][]byte)
	for _, file := range r.Files.LoadFiles {
		name := path.Base(strings.Replace(file.FilePath, "\\", "/", -1))
		if len(file.Rename) > 0 {
			name = file.Rename + path.Ext(file.FilePath)
		}
		storage, e := this.getStorageBySource(file.FilePath, file.IsSource)
		if e.Err != nil {
			return response.LoadZipResponse{}, e
		}
		util.LogEvent(this.Cat, "DownloadZip", this.GetChannel(file.FilePath), map[string]string{"uri": file.FilePath})
		bts, e := storage.Download()
		if e.Err != nil {
			return response.LoadZipResponse{}, e
		}
		files[name] = bts
	}
	bts, e := util.Zip(files)
	if e.Err != nil {
		return response.LoadZipResponse{}, e
	}
	bts = []byte(soapparse.B64.EncodeToString(bts))
	return response.LoadZipResponse{FileBytes: bts}, util.Error{}
}

func (this ImageRequest) Delete(r *request.DeleteRequest) (response.DeleteResponse, util.Error) {
	storage, e := this.getStorage(r.FilePath)
	if e.Err != nil {
		return response.DeleteResponse{}, e
	}
	util.LogEvent(this.Cat, "Delete", this.GetChannel(r.FilePath), map[string]string{"uri": r.FilePath, "IsDeleteAll": strconv.FormatBool(r.IsDeleteAll)})

	if isNewUri(r.FilePath) {
		imgIndex := models.ImageIndex{Cat: this.Cat}
		e = imgIndex.ParseName(r.FilePath)
		if e.Err != nil {
			return response.DeleteResponse{}, e
		}
		e = imgIndex.Delete()
		if e.Err != nil {
			return response.DeleteResponse{}, e
		}
	}
	e = storage.Delete(r.IsDeleteAll)
	if e.Err != nil {
		return response.DeleteResponse{}, e
	}
	return response.DeleteResponse{}, util.Error{}
}

func (this ImageRequest) getStorageBySource(path string, isSource bool) (Storage, util.Error) {
	path = strings.Replace(path, "/", "\\", -1)
	var (
		storagePath string
		storageType string
	)
	var storage Storage
	if isNewUri(path) {
		imagename := util.Substr(path, 1, NEWIMAGENAMELENGTH)
		imgIndex := models.ImageIndex{Cat: this.Cat}
		if e := imgIndex.Parse(imagename); e.Err != nil {
			return nil, e
		}
		storagePath = imgIndex.StoragePath

		storageType = imgIndex.StorageType
		storage = CreateStorage(storagePath, storageType, this.Cat)
		if storage == nil {
			util.LogErrorEvent(this.Cat, ERRORTYPE_STORAGETYPENOSUPPORTE, util.JoinString("Can't supporte storagetype[", storageType, "]"))
			return nil, util.Error{IsNormal: false, Err: errors.New(util.JoinString("Can't supporte storagetype[", storageType, "]")), Type: ERRORTYPE_STORAGETYPENOSUPPORTE}
		}
	} else {
		storageType = STORAGETYPE_NFS
		if isFdfs(path) {
			storageType = STORAGETYPE_FDFS
		}
		storage = CreateStorage(path, storageType, this.Cat)
		if e := storage.ConvertFilePath(isSource); e.Err != nil {
			return nil, e
		}
	}
	return storage, util.Error{}
}

func (this ImageRequest) getStorage(path string) (Storage, util.Error) {
	return this.getStorageBySource(path, false)
}

func (this ImageRequest) GetChannel(path string) string {
	path = strings.Replace(path, "/", "\\", -1)
	if isNewUri(path) {
		channelCode := util.Substr(path, 1, 2)
		channel := models.Channel{Cat: this.Cat}
		channels, _ := channel.GetAll()
		for k, v := range channels {
			if v == channelCode {
				return k
			}
		}
		return ""
	}
	if isFdfs(path) {
		s := util.Substr(path, 4, len(path)-4)
		i := strings.Index(s, "\\")
		return util.Substr(s, 0, i)
	}
	if isT1(path) {
		s := util.Substr(path, 4, len(path)-4)
		i := strings.Index(s, "\\")
		return util.Substr(s, 0, i)
	}
	s := util.Substr(path, 1, len(path)-4)
	i := strings.Index(s, "\\")
	return util.Substr(s, 0, i)
}
func isNewUri(uri string) bool {
	index := strings.Index(uri, ".")
	if index < 1 {
		return false
	}
	s := util.Substr(uri, 1, index-1)
	return len(s) == NEWIMAGENAMELENGTH
}

func isFdfs(uri string) bool {
	return util.Substr(uri, 0, 4) == "\\fd\\"
}

func isT1(uri string) bool {
	return util.Substr(uri, 0, 4) == "\\t1\\"
}
