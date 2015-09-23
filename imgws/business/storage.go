package business

import (
	"errors"
	cat "github.com/ctripcorp/cat.go"
	"github.com/ctripcorp/nephele/fdfs"
	"github.com/ctripcorp/nephele/imgws/models"
	"github.com/ctripcorp/nephele/util"
	"io/ioutil"
	"strings"
	"os/exec"
)

type Storage interface {
	Upload(bts []byte, fileExt string) (string, util.Error)
	Download() ([]byte, util.Error)
	Delete(isDeleteAll bool) (util.Error)
	ConvertFilePath() util.Error
}

var (
	ERRORTYPE_FDFSUPLOADERR     = "fdfs.uploaderr"
	ERRORTYPE_FDFSCONNECTIONERR = "fdfs.connectionerr"
	ERRORTYPE_FDFSDOWNLOADERR   = "fdfs.downloaderr"
	ERRORTYPE_NFSDOWNLOADERR    = "nfs.downloaderr"
	ERRORTYPE_FDFSDELETEERR     = "fdfs.deleteerr"
	ERRORTYPE_NFSDELETEERR      = "nfs.deleteerr"

	STORAGETYPE_FDFS = "fdfs"
	STORAGETYPE_NFS  = "nfs"

	fdfsClient fdfs.FdfsClient
)

func NewStorage(c cat.Cat) (Storage, string) {
	return &FdfsStorage{Path: "", Cat: c}, STORAGETYPE_FDFS
}

func CreateStorage(path, storageType string, c cat.Cat) Storage {
	switch storageType {
	case STORAGETYPE_FDFS:
		return &FdfsStorage{Path: path, Cat: c}
	case STORAGETYPE_NFS:
		return &NfsStorage{Path: path, Cat: c}
	default:
		return nil
	}
}

type FdfsStorage struct {
	Path string
	Cat  cat.Cat
}

var (
	uploadcount int = 0
	count           = 0
	lock            = make(chan int, 1)
)

func (this *FdfsStorage) Upload(bts []byte, fileExt string) (string, util.Error) {
	if e := initFdfsClient(); e.Err != nil {
		return "", e
	}
	groups, e := models.GetGroups()
	if e.Err != nil {
		return "", e
	}
	if uploadcount == 99999999 {
		uploadcount = 0
	}
	i := uploadcount % len(groups)
	uploadcount = 0
	g := groups[i]
	path, err := fdfsClient.UploadByBuffer(g, bts, fileExt)
	if err != nil {
		return "", util.Error{IsNormal: true, Err: err, Type: ERRORTYPE_FDFSUPLOADERR}
	}
	return path, util.Error{}
}

func (this *FdfsStorage) Download() ([]byte, util.Error) {
	if e := initFdfsClient(); e.Err != nil {
		return nil, e
	}
	var (
		bts []byte
		err error
	)

	bts, err = fdfsClient.DownloadToBuffer(this.Path, this.Cat)
	if err != nil {
		return []byte{}, util.Error{IsNormal: false, Err: err, Type: ERRORTYPE_FDFSDOWNLOADERR}
	}
	return bts, util.Error{}
}

//fdfs ignore isDeleteAll
func (this *FdfsStorage) Delete(isDeleteAll bool) (util.Error) {
	if e := initFdfsClient(); e.Err != nil {
		return e
	}
	var err error
	err = fdfsClient.DeleteFile(this.Path)
	if err != nil {
		return util.Error{IsNormal: false, Err: err, Type: ERRORTYPE_FDFSDELETEERR}
	}
	return util.Error{}
}


func (this *FdfsStorage) ConvertFilePath() util.Error {
	this.Path = strings.Replace(this.Path, "\\", "/", -1)
	this.Path = util.Substr(this.Path, 4, len(this.Path)-4)
	index := strings.Index(this.Path, "/")
	this.Path = util.Substr(this.Path, index+1, len(this.Path)-index-1)
	return util.Error{}
}

func initFdfsClient() util.Error {
	if fdfsClient == nil {
		lock <- 1
		defer func() {
			<-lock
		}()
		if fdfsClient != nil {
			return util.Error{}
		}
		fdfsdomain, e := models.GetFdfsDomain()
		if e.Err != nil {
			return e
		}
		fdfsport, e := models.GetFdfsPort()
		if e.Err != nil {
			return e
		}
		var err error
		fdfsClient, err = fdfs.NewFdfsClient([]string{fdfsdomain}, fdfsport)
		if err != nil {
			return util.Error{IsNormal: false, Err: err, Type: ERRORTYPE_FDFSCONNECTIONERR}
		}
	}
	return util.Error{}
}

type NfsStorage struct {
	Path string
	Cat  cat.Cat
}

func (this *NfsStorage) Upload(bts []byte, fileExt string) (string, util.Error) {
	return "", util.Error{}
}

func (this *NfsStorage) Download() ([]byte, util.Error) {
	var (
		bts []byte
		err error
	)
	bts, err = ioutil.ReadFile(this.Path)
	if err != nil {
		return []byte{}, util.Error{IsNormal: false, Err: errors.New("download file failed!"), Type: ERRORTYPE_NFSDOWNLOADERR}
	}
	return bts, util.Error{}
}

func (this *NfsStorage) Delete(isDeleteAll bool) (util.Error) {
	var (
		cmd *exec.Cmd
		err error
	)
	if !isDeleteAll {
		cmd = exec.Command("rm", this.Path)
	} else {
		cmd = exec.Command("rm", util.Substr(this.Path, 4, len(this.Path)-4)+"_*")
	}
	err = cmd.Run()
	if err != nil {
		return util.Error{IsNormal: false, Err: err, Type: ERRORTYPE_NFSDELETEERR}
	}
	return util.Error{}
}


func (this *NfsStorage) ConvertFilePath() util.Error {
	this.Path = strings.Replace(this.Path, "/", "\\", -1)

	if this.isT1() {
		this.Path = util.Substr(this.Path, 4, len(this.Path)-4)
		index := strings.Index(this.Path, "\\")
		channel := util.Substr(this.Path, 0, index)
		nfs, e := models.GetNfsPath(channel)
		if e.Err != nil {
			return e
		}

		this.Path = shading(nfs) + this.Path
	} else {
		index := strings.Index(this.Path, "\\")
		channel := util.Substr(this.Path, 1, index)
		nfs, e := models.GetNfsT1Path(channel)
		if e.Err != nil {
			return e
		}

		this.Path = shading(nfs) + this.Path
	}
	this.Path = strings.Replace(this.Path, "\\", "/", -1)
	return util.Error{}
}

func (this NfsStorage) isT1() bool {
	return util.Substr(this.Path, 0, 4) == "\\t1\\"
}

func shading(arr []string) string {
	if count == 9999999 {
		count = 0
	}
	i := count % len(arr)
	count = count + 1
	return arr[i]
}
