package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/ctripcorp/nephele/util"
	"strconv"
	"strings"
	"time"
)

var (
	ERRTYPE_IMAGENAMEINVALID   = "ImageNameInvalid"
	ERRTYPE_GETIMAGEINDEX      = "GetImageIndex"
	ERRORTYPE_INSERTIMAGEINDEX = "InsertImageIndex"
	ERRORTYPE_INSERTIMAGEPLAN  = "InsertImagePlan"
	NEWIMAGENAMELENGTH         = 21
)

type ImageIndex struct {
	Idx          int64
	Channel      string
	StoragePath  string
	StorageType  string
	Profile      string
	PartitionKey int16
	TableZone    int
}

func getDBString(tableZone int) string {
	return "default"
}
func getExt(path string) string {
	arr := strings.Split(path, ".")
	return arr[len(arr)-1]
}

func (this *ImageIndex) SaveToDB(plan string) util.Error {
	o := orm.NewOrm()
	o.Using(getDBString(this.TableZone))
	o.Begin()
	partitionKey := util.GetPartitionKey(time.Now())
	res, err := o.Raw("INSERT INTO `imageindex_"+strconv.Itoa(this.TableZone)+"`(`channel`,`storagePath`,`storageType`,`profile`,`createtime`,`partitionKey`)VALUES(?,?,?,?,NOW(),?)", this.Channel, this.StoragePath, this.StorageType, this.Profile, partitionKey).Exec()
	if err != nil {
		o.Rollback()
		return util.Error{IsNormal: false, Err: err, Type: ERRORTYPE_INSERTIMAGEINDEX}
	}
	id, err := res.LastInsertId()
	if err != nil {
		o.Rollback()
		return util.Error{IsNormal: false, Err: err, Type: ERRORTYPE_INSERTIMAGEINDEX}
	}
	if plan == "" {
		res, err := o.Raw("INSERT INTO 'imageplan_"+strconv.Itoa(this.TableZone)+"'(imgIdx,plan,partitionKey)VALUES(?,?,?)", id, plan, partitionKey).Exec()
		if err != nil {
			o.Rollback()
			return util.Error{IsNormal: false, Err: err, Type: ERRORTYPE_INSERTIMAGEPLAN}
		}
		if _, err := res.RowsAffected(); err != nil {
			o.Rollback()
			return util.Error{IsNormal: false, Err: err, Type: ERRORTYPE_INSERTIMAGEPLAN}
		}
	}
	if err := o.Commit(); err != nil {
		return util.Error{IsNormal: false, Err: err, Type: ERRORTYPE_INSERTIMAGEINDEX}
	}
	this.Idx = id
	this.PartitionKey = partitionKey
	return util.Error{}
}

func (this ImageIndex) GetImageName() string {
	ext := getExt(this.StoragePath)                              //扩展名
	zone := strconv.FormatInt(int64(this.TableZone), 36)         //转36进制
	partition := strconv.FormatInt(int64(this.PartitionKey), 36) //转36进制
	//1~2 频道 3~4 分区 5~6 时间 7 版本号 8~17 索引 18~21 检验码
	tmp := util.JoinString(this.Channel, util.Cover(zone, "0", 2), util.Cover(partition, "0", 2), "0", util.Cover(strconv.FormatInt(this.Idx, 10), "0", 10))
	return util.JoinString("\\", tmp, util.Compute(tmp), ".", ext)
}

func (this *ImageIndex) Parse(imageName string) util.Error {
	imgName := this.DropExtension(imageName)
	if e := this.ParseName(imgName); e.Err != nil {
		return e
	}
	if e := this.GetStorage(); e.Err != nil {
		return e
	}
	return util.Error{}
}

func (this *ImageIndex) DropExtension(imageName string) string {
	return strings.Split(imageName, ".")[0]
}

func (this *ImageIndex) ParseName(imageName string) util.Error {
	if len(imageName) != NEWIMAGENAMELENGTH {
		return util.Error{IsNormal: false, Err: errors.New("imagename length is invalid"), Type: ERRTYPE_IMAGENAMEINVALID}
	}
	swithoutcompute := util.Substr(imageName, 0, 17)
	scompute := util.Substr(imageName, 18, 4)
	if util.Compute(swithoutcompute) != scompute {
		return util.Error{IsNormal: true, Err: errors.New("Compute check faile."), Type: ERRTYPE_IMAGENAMEINVALID}
	}
	channel := util.Substr(imageName, 0, 2)

	tableZone, err := strconv.ParseInt(util.Substr(imageName, 2, 2), 36, 10)
	if err != nil {
		return util.Error{IsNormal: true, Err: errors.New("tablezone is invalid."), Type: ERRTYPE_IMAGENAMEINVALID}
	}
	partitionKey, err := strconv.ParseInt(util.Substr(imageName, 4, 2), 36, 10)
	if err != nil {
		return util.Error{IsNormal: true, Err: errors.New("partition is invalid."), Type: ERRTYPE_IMAGENAMEINVALID}
	}
	idx, err := strconv.ParseInt(util.Substr(imageName, 7, 10), 10, 10)
	if err != nil {
		return util.Error{IsNormal: true, Err: errors.New("index is invalid."), Type: ERRTYPE_IMAGENAMEINVALID}
	}
	this.Channel = channel
	this.Idx = idx
	this.PartitionKey = int16(partitionKey)
	this.TableZone = int(tableZone)
	return util.Error{}
}

func (this *ImageIndex) GetStorage() util.Error {
	if this.Idx < 1 || this.TableZone < 1 || this.PartitionKey < 1 {
		return util.Error{IsNormal: true, Err: errors.New("getimageindex parameters is invalid"), Type: ERRTYPE_GETIMAGEINDEX}
	}
	o := orm.NewOrm()
	o.Using(getDBString(this.TableZone))
	res := make(orm.Params)
	nums, err := o.Raw(util.JoinString("SELECT storagepath,storagetype FROM imageindex_", strconv.Itoa(this.TableZone), " WHERE idx=? AND partitionKey=?"), this.Idx, this.PartitionKey).RowsToMap(&res, "storagepath", "storagetype")
	if err != nil {
		return util.Error{IsNormal: false, Err: err, Type: ERRTYPE_GETIMAGEINDEX}
	}
	if nums < 1 {
		return util.Error{IsNormal: false, Err: errors.New("idx is't exists"), Type: ERRTYPE_GETIMAGEINDEX}
	}
	for k, v := range res {
		this.StoragePath = k
		this.StorageType = fmt.Sprintf("%v", v)
	}
	return util.Error{}
}
