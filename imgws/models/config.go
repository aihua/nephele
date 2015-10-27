package models

import (
	"database/sql"
	"errors"
	"github.com/astaxie/beego/orm"
	cat "github.com/ctripcorp/cat.go"
	"github.com/ctripcorp/nephele/util"
	"strconv"
	"strings"
	"time"
)

var (
	KeySizes = "sizes"
)

var (
	ErrDuplicatedRecord = errors.New("duplicated records exist")
	ErrSizeExisted      = errors.New("size already existed")
	ErrNoRecord         = errors.New("no record found")
	ErrHasBeenModified  = errors.New("Has been modified")

	ERRORTYPE_GETCONFIGS     = "GetConfig"
	ERRORTYPE_CONFIGNOEXISTS = "ConfigNoExists"

	CONFIG_DEFAULTCHANNEL = "00"
	CONFIG_FDFSGROUPS     = "fdfsgroups"
	CONFIG_FDFSPORT       = "fdfsport"
	CONFIG_FDFSDOMAIN     = "fdfsdomain"
	CONFIG_NFS            = "nfs"
	CONFIG_NFST1          = "nfst1"

	SPILT_1 = ","
	SPILT_2 = "|"

	config        map[string]map[string]string = make(map[string]map[string]string)
	getConfigTime time.Time                    = time.Now()
)

type Config struct {
	//Consistent Fields
	Id          int64 `orm:"auto"`
	ChannelCode string
	Key         string

	//Inconsistent Fields
	Value      string
	Recordtime string
	Cat        cat.Cat
}

func getNow() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
func (this *Config) Insert() error {
	this.Recordtime = getNow()
	var (
		err error
		id  int64
		res sql.Result
	)
	if this.Cat != nil {
		tran := this.Cat.NewTransaction(DBTITLE, "Config.Insert")
		defer func() {
			if err != nil {
				tran.SetStatus(err)
			} else {
				tran.SetStatus("0")
			}
			tran.Complete()
		}()
	}
	o := orm.NewOrm()
	res, err = o.Raw("INSERT INTO config(channelCode,`key`,value,recordTime)VALUES(?,?,?,?)", this.ChannelCode, this.Key, this.Value, this.Recordtime).Exec()
	if err != nil {
		return err
	}
	id, err = res.LastInsertId()
	if err != nil {
		return err
	}
	this.Id = id
	return nil
}

func (this *Config) UpdateValue() error {
	var err error
	if this.Cat != nil {
		tran := this.Cat.NewTransaction(DBTITLE, "Config.Update")
		defer func() {
			if err != nil {
				tran.SetStatus(err)
			} else {
				tran.SetStatus("0")
			}
			tran.Complete()
		}()
	}

	o := orm.NewOrm()
	now := getNow()
	_, err = o.Raw("UPDATE config SET value=?, recordtime=? WHERE channelCode=? AND `key`=?", this.Value, now, this.ChannelCode, this.Key).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (this *Config) GetSizes() (string, error) {
	this.Key = "sizes"
	var (
		num int64
		err error
	)

	if this.Cat != nil {
		tran := this.Cat.NewTransaction(DBTITLE, "Config.GetSizes")
		defer func() {
			if err != nil {
				tran.SetStatus(err)
			} else {
				tran.SetStatus("0")
			}
			tran.Complete()
		}()
	}
	o := orm.NewOrm()
	var maps []orm.Params
	num, err = o.Raw("SELECT value,recordTime FROM config WHERE `channelCode` = ? AND `key` = ?", this.ChannelCode, this.Key).Values(&maps)
	if err != nil {
		return EmptyString, err
	}
	if num == 0 {
		return EmptyString, ErrNoRecord
	} else if num == 1 {
		this.Value = maps[0]["value"].(string)
		this.Recordtime = maps[0]["recordTime"].(string)
		return this.Value, nil
	} else {
		return EmptyString, ErrDuplicatedRecord
	}

}

func (this *Config) AddSize(size string) error {
	this.Key = "sizes"
	var (
		result sql.Result
		err    error
		num    int64
	)
	if this.Cat != nil {
		tran := this.Cat.NewTransaction(DBTITLE, "Config.AddSizes")
		defer func() {
			if err != nil {
				tran.SetStatus(err)
			} else {
				tran.SetStatus("0")
			}
			tran.Complete()
		}()
	}
	o := orm.NewOrm()
	var maps []orm.Params
	num, err = o.Raw("SELECT * FROM config WHERE `channelCode` = ? AND `key` = ?", this.ChannelCode, this.Key).Values(&maps)
	if err != nil {
		return err
	}
	if num == 0 {
		this.Value = size
		this.Recordtime = getNow()
		err = this.Insert()
		return err
	} else if num == 1 {
		sizes := strings.Split(maps[0]["value"].(string), ",")
		isExisted := false
		for _, t := range sizes {
			if t == size {
				isExisted = true
			}
		}
		if isExisted {
			err = ErrSizeExisted
			return err
		} else {
			newValue := maps[0]["value"].(string) + "," + size
			newRecordtime := getNow()

			result, err = o.Raw("UPDATE config set `value` = ?, `recordtime` = ? WHERE `channelCode` = ? AND `key` = ? AND `recordtime` = ?", newValue, newRecordtime, this.ChannelCode, this.Key, maps[0]["recordTime"]).Exec()
			if err != nil {
				return err
			}
			num, err = result.RowsAffected()
			if err != nil {
				return err
			}
			if num < 1 {
				err = ErrHasBeenModified
				return err
			}
			this.Value = newValue
			this.Recordtime = newRecordtime
			return nil
		}
	} else {
		return ErrDuplicatedRecord
	}
}

func (this *Config) GetConfigs() (map[string]map[string]string, util.Error) {
	if len(config) < 1 || isRefresh(getConfigTime) {
		var err error
		if this.Cat != nil {
			tran := this.Cat.NewTransaction(DBTITLE, "Config.GetConfigs")
			defer func() {
				if err != nil {
					tran.SetStatus(err)
				} else {
					tran.SetStatus("0")
				}
				tran.Complete()
			}()
		}

		o := orm.NewOrm()
		var list []orm.Params
		_, err = o.Raw("SELECT channelCode,`key`,value FROM config").Values(&list)
		if err != nil {
			util.LogErrorEvent(this.Cat, ERRORTYPE_GETCONFIGS, err.Error())
			return nil, util.Error{IsNormal: false, Err: err, Type: ERRORTYPE_GETCONFIGS}
		}
		m := make(map[string]map[string]string)
		for _, v := range list {
			channelCode := v["channelCode"].(string)
			key := v["key"].(string)
			value := v["value"].(string)
			_, exists := m[channelCode]
			if exists {
				m[channelCode][key] = value
			} else {
				childmap := make(map[string]string)
				childmap[key] = value
				m[channelCode] = childmap
			}
		}
		config = m
		getConfigTime = time.Now()
	}
	return config, util.Error{}
}

func isRefresh(t time.Time) bool {
	return t.Add(1 * time.Minute).Before(time.Now())
}

func (this *Config) GetChannelConfigs(channelCode string) (map[string]string, util.Error) {
	configs, e := this.GetConfigs()
	if e.Err != nil {
		return nil, e
	}
	config, exists := configs[channelCode]
	if !exists {
		return nil, util.Error{IsNormal: false, Err: errors.New(util.JoinString("channelCode[", channelCode, "] Config is't exists!")), Type: ERRORTYPE_CONFIGNOEXISTS}
	}
	return config, util.Error{}
}
func (this *Config) GetGroups() ([]string, util.Error) {
	configs, e := this.GetChannelConfigs(CONFIG_DEFAULTCHANNEL)
	if e.Err != nil {
		return []string{}, e
	}
	groups, exists := configs[CONFIG_FDFSGROUPS]
	if !exists {
		return nil, util.Error{IsNormal: false, Err: errors.New(util.JoinString("channel[", CONFIG_DEFAULTCHANNEL, "] Config[", CONFIG_FDFSGROUPS, "] is't exists!")), Type: ERRORTYPE_CONFIGNOEXISTS}
	}

	if len(groups) < 1 {
		return nil, util.Error{IsNormal: false, Err: errors.New(util.JoinString("channel[", CONFIG_DEFAULTCHANNEL, "] Config[", CONFIG_FDFSGROUPS, "] is't exists!")), Type: ERRORTYPE_CONFIGNOEXISTS}
	}
	return strings.Split(groups, SPILT_1), util.Error{}
}

func (this *Config) GetFdfsDomain() (string, util.Error) {
	configs, e := this.GetChannelConfigs(CONFIG_DEFAULTCHANNEL)
	if e.Err != nil {
		return "", e
	}
	domain, exists := configs[CONFIG_FDFSDOMAIN]
	if !exists {
		return "", util.Error{IsNormal: false, Err: errors.New(util.JoinString("channel[", CONFIG_DEFAULTCHANNEL, "] Config[", CONFIG_FDFSDOMAIN, "] is't exists!")), Type: ERRORTYPE_CONFIGNOEXISTS}
	}
	return domain, util.Error{}
}

func (this *Config) GetFdfsPort() (string, util.Error) {
	configs, e := this.GetChannelConfigs(CONFIG_DEFAULTCHANNEL)
	if e.Err != nil {
		return "", e
	}
	port, exists := configs[CONFIG_FDFSPORT]
	if !exists {
		return "", util.Error{IsNormal: false, Err: errors.New(util.JoinString("channel[", CONFIG_DEFAULTCHANNEL, "] Config[", CONFIG_FDFSPORT, "] is't exists!")), Type: ERRORTYPE_CONFIGNOEXISTS}
	}
	return port, util.Error{}
}

func (this *Config) GetNfsPath(channel string) ([]string, util.Error) {
	value, e := this.getValue(channel, CONFIG_NFS)
	if e.Err != nil {
		return []string{}, e
	}
	return strings.Split(value, SPILT_1), util.Error{}
}

func (this *Config) GetNfsT1Path(channel string) ([]string, util.Error) {
	value, e := this.getValue(channel, CONFIG_NFST1)
	if e.Err != nil {
		return []string{}, e
	}
	return strings.Split(value, SPILT_1), util.Error{}
}

func (this *Config) getValue(channel, key string) (string, util.Error) {
	configs, e := this.GetConfigs()
	if e.Err != nil {
		return "", e
	}

	var (
		value  string
		exists bool
	)
	channelConfigs, exists := configs[channel]
	if exists {
		value, exists = channelConfigs[key]
		if exists {
			return value, util.Error{}
		}
	}
	if channel != CONFIG_DEFAULTCHANNEL {
		defaultConfigs, exists := configs[CONFIG_DEFAULTCHANNEL]
		if exists {
			value, exists = defaultConfigs[key]
			if exists {
				return value, util.Error{}
			}
		}
	}
	return "", util.Error{IsNormal: false, Err: errors.New(util.JoinString("Channel[", channel, "] Key[", key, "] is't exists!")), Type: ERRORTYPE_CONFIGNOEXISTS}
}
