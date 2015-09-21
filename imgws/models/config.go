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
	KeySizes = "sizes"
)

var (
	ErrDuplicatedRecord = errors.New("duplicated records exist")
	ErrSizeExisted      = errors.New("size already existed")
	ErrNoRecord         = errors.New("no record found")

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
	Id      int64 `orm:"auto"`
	Channel string
	Key     string

	//Inconsistent Fields
	Value      string
	Recordtime string
}

func (this *Config) GetSizes() (string, error) {
	this.Key = "sizes"
	var (
		configs []Config
		sql     string
	)

	o := orm.NewOrm()
	o.Using("default")

	sql = fmt.Sprintf("SELECT * FROM config WHERE `channel` = '%s' AND `key` = '%s'", this.Channel, this.Key)
	num, err := o.Raw(sql).QueryRows(&configs)
	if err != nil {
		return EmptyString, err
	}
	if num == 0 {
		return EmptyString, ErrNoRecord
	} else if num == 1 {
		this.Value = configs[0].Value
		this.Recordtime = configs[0].Recordtime
		return configs[0].Value, nil
	} else {
		return EmptyString, ErrDuplicatedRecord
	}

}

func (this *Config) AddSize(size string) error {
	this.Key = "sizes"
	var (
		configs []Config
		sql     string
	)

	o := orm.NewOrm()
	o.Using("default")

	sql = fmt.Sprintf("SELECT * FROM config WHERE `channel` = '%s' AND `key` = '%s'", this.Channel, this.Key)
	num, err := o.Raw(sql).QueryRows(&configs)
	if err != nil {
		return err
	}
	if num == 0 {
		this.Value = size
		this.Recordtime = strconv.FormatInt(time.Now().UnixNano(), 10)
		id, err := o.Insert(this)
		if err != nil {
			return err
		}
		this.Id = id
		return nil
	} else if num == 1 {
		sizes := strings.Split(configs[0].Value, ",")
		isExisted := false
		for _, t := range sizes {
			if t == size {
				isExisted = true
			}
		}
		if isExisted {
			return ErrSizeExisted
		} else {
			newValue := configs[0].Value + "," + size
			newRecordtime := strconv.FormatInt(time.Now().UnixNano(), 10)

			sql = fmt.Sprintf("UPDATE config set `value` = '%s', `recordtime` = '%s' "+
				"WHERE `channel` = '%s' AND `key` = '%s' AND `recordtime` = '%s'",
				newValue, newRecordtime,
				configs[0].Channel, configs[0].Key, configs[0].Recordtime)

			_, err := o.Raw(sql).Exec()
			if err != nil {
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

func GetConfigs() (map[string]map[string]string, util.Error) {
	if len(config) < 1 || isRefresh(getConfigTime) {
		o := orm.NewOrm()
		var configs []Config
		_, err := o.Raw("SELECT channel,key,value FROM config").QueryRows(&configs)
		if err != nil {
			return nil, util.Error{IsNormal: false, Err: err, Type: ERRORTYPE_GETCONFIGS}
		}
		m := make(map[string]map[string]string)
		for _, config := range configs {
			_, exists := m[config.Channel]
			if exists {
				m[config.Channel][config.Key] = config.Value
			} else {
				childmap := make(map[string]string)
				childmap[config.Key] = config.Value
				m[config.Channel] = childmap
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

func GetChannelConfigs(channel string) (map[string]string, util.Error) {
	configs, e := GetConfigs()
	if e.Err != nil {
		return nil, e
	}
	config, exists := configs[channel]
	if !exists {
		return nil, util.Error{IsNormal: false, Err: errors.New(util.JoinString("channel[", channel, "] Config is't exists!")), Type: ERRORTYPE_CONFIGNOEXISTS}
	}
	return config, util.Error{}
}
func GetGroups() ([]string, util.Error) {
	configs, e := GetChannelConfigs(CONFIG_DEFAULTCHANNEL)
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

func GetFdfsDomain() (string, util.Error) {
	configs, e := GetChannelConfigs(CONFIG_DEFAULTCHANNEL)
	if e.Err != nil {
		return "", e
	}
	domain, exists := configs[CONFIG_FDFSDOMAIN]
	if !exists {
		return "", util.Error{IsNormal: false, Err: errors.New(util.JoinString("channel[", CONFIG_DEFAULTCHANNEL, "] Config[", CONFIG_FDFSDOMAIN, "] is't exists!")), Type: ERRORTYPE_CONFIGNOEXISTS}
	}
	return domain, util.Error{}
}

func GetFdfsPort() (string, util.Error) {
	configs, e := GetChannelConfigs(CONFIG_DEFAULTCHANNEL)
	if e.Err != nil {
		return "", e
	}
	port, exists := configs[CONFIG_FDFSPORT]
	if !exists {
		return "", util.Error{IsNormal: false, Err: errors.New(util.JoinString("channel[", CONFIG_DEFAULTCHANNEL, "] Config[", CONFIG_FDFSPORT, "] is't exists!")), Type: ERRORTYPE_CONFIGNOEXISTS}
	}
	return port, util.Error{}
}

func GetNfsPath(channel string) ([]string, util.Error) {
	value, e := getValue(channel, CONFIG_NFS)
	if e.Err != nil {
		return []string{}, e
	}
	return strings.Split(value, SPILT_1), util.Error{}
}

func GetNfsT1Path(channel string) ([]string, util.Error) {
	value, e := getValue(channel, CONFIG_NFST1)
	if e.Err != nil {
		return []string{}, e
	}
	return strings.Split(value, SPILT_1), util.Error{}
}

func getValue(channel, key string) (string, util.Error) {
	configs, e := GetChannelConfigs(channel)
	if e.Err != nil {
		return "", e
	}
	var (
		value  string
		exists bool
	)
	value, exists = configs[key]
	if exists {
		return value, util.Error{}
	}
	defaultConfigs, e := GetChannelConfigs(CONFIG_DEFAULTCHANNEL)
	if e.Err != nil {
		return "", e
	}
	value, exists = defaultConfigs[key]
	if exists {
		return value, util.Error{}
	}
	return "", util.Error{IsNormal: false, Err: errors.New(util.JoinString("Channel[", channel, "] Key[", key, "] is't exists!")), Type: ERRORTYPE_CONFIGNOEXISTS}
}
