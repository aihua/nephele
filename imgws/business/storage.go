package business

import (
	cat "github.com/ctripcorp/cat.go"
	"github.com/ctripcorp/nephele/fdfs"
	"github.com/ctripcorp/nephele/imgws/models"
	"github.com/ctripcorp/nephele/util"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

type Storage interface {
	Upload(bts []byte, fileExt string) (string, util.Error)
	Download() ([]byte, util.Error)
	Delete(isDeleteAll bool) util.Error
	ConvertFilePath(isSource bool) util.Error
	UploadSlave(bts []byte, prefixName string, fileExtName string) (string, util.Error)
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
	CATTITLE        = "Storage"
)

func (this *FdfsStorage) Upload(bts []byte, fileExt string) (string, util.Error) {
	if e := this.initFdfsClient(); e.Err != nil {
		return "", e
	}
	config := models.Config{Cat: this.Cat}
	groups, e := config.GetGroups()
	if e.Err != nil {
		return "", e
	}
	if uploadcount == 99999999 {
		uploadcount = 0
	}
	i := uploadcount % len(groups)
	uploadcount = 0
	g := groups[i]

	var result util.Error = util.Error{}
	if this.Cat != nil {
		tran := this.Cat.NewTransaction(CATTITLE, "Fdfs.Upload")
		util.LogEvent(this.Cat, "Size", util.GetImageSizeDistribution(len(bts)), map[string]string{"size": strconv.Itoa(len(bts))})

		defer func() {
			if result.Err != nil && result.IsNormal {
				tran.SetStatus(result.Err)
			} else {
				tran.SetStatus("0")
			}
			tran.Complete()
		}()
	}
	path, err := fdfsClient.UploadByBuffer(g, bts, fileExt)
	if err != nil {
		util.LogErrorEvent(this.Cat, ERRORTYPE_FDFSUPLOADERR, err.Error())
		result = util.Error{IsNormal: false, Err: err, Type: ERRORTYPE_FDFSUPLOADERR}
		return "", result
	}
	return path, result
}

func (this *FdfsStorage) UploadSlave(bts []byte, prefixName string, fileExtName string) (string, util.Error) {
	if e := this.initFdfsClient(); e.Err != nil {
		return "", e
	}

	var result util.Error = util.Error{}
	if this.Cat != nil {
		tran := this.Cat.NewTransaction(CATTITLE, "Fdfs.UploadSlave")
		defer func() {
			if result.Err != nil && result.IsNormal {
				tran.SetStatus(result.Err)
			} else {
				tran.SetStatus("0")
			}
			tran.Complete()
		}()
	}

	path, err := fdfsClient.UploadSlaveByBuffer(bts, this.Path, prefixName, fileExtName)
	if err != nil {
		result = util.Error{IsNormal: true, Err: err, Type: ERRORTYPE_FDFSUPLOADERR}
		return "", result
	}

	return path, result
}
func (this *FdfsStorage) Download() ([]byte, util.Error) {
	if e := this.initFdfsClient(); e.Err != nil {
		return nil, e
	}
	var (
		bts    []byte
		err    error
		result = util.Error{}
	)
	var messagesize int
	if this.Cat != nil {
		tran := this.Cat.NewTransaction(CATTITLE, "Fdfs.Download")
		tran.AddData("path", this.Path)
		defer func() {
			if result.Err != nil && result.IsNormal {
				tran.SetStatus(result.Err)
			} else {
				util.LogEvent(this.Cat, "Size", util.GetImageSizeDistribution(messagesize), map[string]string{"size": strconv.Itoa(messagesize)})
				tran.SetStatus("0")
			}
			tran.Complete()
		}()
	}
	bts, err = fdfsClient.DownloadToBuffer(this.Path, this.Cat)
	if err != nil {
		util.LogErrorEvent(this.Cat, ERRORTYPE_FDFSDOWNLOADERR, err.Error())
		result = util.Error{IsNormal: false, Err: err, Type: ERRORTYPE_FDFSDOWNLOADERR}
		return []byte{}, result
	}
	messagesize = len(bts)
	return bts, result
}

//fdfs ignore isDeleteAll
func (this *FdfsStorage) Delete(isDeleteAll bool) util.Error {
	if e := this.initFdfsClient(); e.Err != nil {
		return e
	}
	var (
		err    error
		result = util.Error{}
	)
	if this.Cat != nil {
		tran := this.Cat.NewTransaction(CATTITLE, "Fdfs.Delete")
		tran.AddData("path", this.Path)
		defer func() {
			if result.Err != nil && result.IsNormal {
				tran.SetStatus(result.Err)
			} else {
				tran.SetStatus("0")
			}
			tran.Complete()
		}()
	}
	err = fdfsClient.DeleteFile(this.Path)
	if err != nil {
		util.LogErrorEvent(this.Cat, ERRORTYPE_FDFSDELETEERR, err.Error())
		result = util.Error{IsNormal: false, Err: err, Type: ERRORTYPE_FDFSDELETEERR}
		return result
	}
	return result
}

func (this *FdfsStorage) ConvertFilePath(isSource bool) util.Error {
	this.Path = strings.Replace(this.Path, "\\", "/", -1)
	this.Path = util.Substr(this.Path, 4, len(this.Path)-4)
	index := strings.Index(this.Path, "/")
	this.Path = util.Substr(this.Path, index+1, len(this.Path)-index-1)
	if isSource {
		this.Path = strings.Replace(this.Path, ".", "_source.", -1)
	}
	return util.Error{}
}

func (this *FdfsStorage) initFdfsClient() util.Error {
	if fdfsClient == nil {
		lock <- 1
		defer func() {
			<-lock
		}()
		if fdfsClient != nil {
			return util.Error{}
		}
		config := models.Config{Cat: this.Cat}
		fdfsdomain, e := config.GetFdfsDomain()
		if e.Err != nil {
			return e
		}
		fdfsport, e := config.GetFdfsPort()
		if e.Err != nil {
			return e
		}
		var err error
		fdfsClient, err = fdfs.NewFdfsClient([]string{fdfsdomain}, fdfsport)
		if err != nil {
			util.LogErrorEvent(this.Cat, ERRORTYPE_FDFSCONNECTIONERR, err.Error())
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
func (this *NfsStorage) UploadSlave(bts []byte, prefixName string, fileExtName string) (string, util.Error) {
	return "", util.Error{}
}

func (this *NfsStorage) Download() ([]byte, util.Error) {
	var (
		bts []byte
		err error
	)
	var messagesize int
	if this.Cat != nil {
		tran := this.Cat.NewTransaction(CATTITLE, "Nfs.Download")
		tran.AddData("path", this.Path)
		defer func() {
			if err != nil {
				tran.SetStatus(err)
			} else {
				util.LogEvent(this.Cat, "Size", util.GetImageSizeDistribution(messagesize), map[string]string{"size": strconv.Itoa(messagesize)})
				tran.SetStatus("0")
			}
			tran.Complete()
		}()
	}
	bts, err = ioutil.ReadFile(this.Path)
	if err != nil {
		util.LogErrorEvent(this.Cat, ERRORTYPE_NFSDOWNLOADERR, err.Error())
		return []byte{}, util.Error{IsNormal: false, Err: err, Type: ERRORTYPE_NFSDOWNLOADERR}
	}
	messagesize = len(bts)
	return bts, util.Error{}
}

func (this *NfsStorage) Delete(isDeleteAll bool) util.Error {
	var (
		cmd *exec.Cmd
		err error
	)
	if this.Cat != nil {
		tran := this.Cat.NewTransaction(CATTITLE, "Nfs.Delete")
		tran.AddData("path", this.Path)
		defer func() {
			if err != nil {
				tran.SetStatus(err)
			} else {
				tran.SetStatus("0")
			}
			tran.Complete()
		}()
	}
	if !isDeleteAll {
		cmd = exec.Command("rm", this.Path)
	} else {
		cmd = exec.Command("rm", util.Substr(this.Path, 4, len(this.Path)-4)+"_*")
	}
	err = cmd.Run()
	if err != nil {
		util.LogErrorEvent(this.Cat, ERRORTYPE_NFSDELETEERR, err.Error())
		return util.Error{IsNormal: false, Err: err, Type: ERRORTYPE_NFSDELETEERR}
	}
	return util.Error{}
}

func (this *NfsStorage) ConvertFilePath(isSource bool) util.Error {
	this.Path = strings.Replace(this.Path, "/", "\\", -1)

	config := models.Config{Cat: this.Cat}
	if this.isT1() {
		this.Path = util.Substr(this.Path, 4, len(this.Path)-4)
		index := strings.Index(this.Path, "\\")
		channel := util.Substr(this.Path, 0, index)
		nfs, e := config.GetNfsPath(channel)
		if e.Err != nil {
			return e
		}

		this.Path = shading(nfs) + this.Path
	} else {
		index := strings.Index(this.Path, "\\")
		channel := util.Substr(this.Path, 1, index)
		nfs, e := config.GetNfsT1Path(channel)
		if e.Err != nil {
			return e
		}

		this.Path = shading(nfs) + this.Path
	}
	this.Path = strings.Replace(this.Path, "\\", "/", -1)
	if isSource {
		this.Path = strings.Replace(this.Path, "\\target\\", "\\source\\", -1)
	}
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
