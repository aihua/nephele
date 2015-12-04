package config

import (
	"github.com/Unknwon/goconfig"
	"strconv"
	"strings"
)

type ProcConfig interface {
	GetFdfsDomain() (string, error)
	GetFdfsPort() int
	GetDirPath(channel, storagetype string) (string, error)
	GetResizeTypes(channel string) (string, error)
	GetSizes(channel string) (string, error)
	GetRotates(channel string) (string, error)
	GetQuality(channel string) (string, error)
	GetQualities(channel string) (string, error)
	GetLogodir(channel string) (string, error)
	IsEnableNameLogo(channel string) (bool, error)
	GetDefaultLogo(channel string) (string, error)
	GetLogoNames(channel string) (string, error)
	GetImagelesswidthForLogo(channel string) (int64, error)
	GetImagelessheightForLogo(channel string) (int64, error)
	GetDissolves(channel string) (string, error)
	GetDissolve(channel string) int
	GetNamelogoDissolve(channel string) int
	GetSequenceofoperation(channel string) ([]string, error)
	Reload() error
}

type LocalConfig struct {
	conf *goconfig.ConfigFile
}

func NewLocalConfig(env Env) *LocalConfig {
	var confPath string
	if env == uat {
		confPath = "../conf/uat_conf.ini"
	} else {
		confPath = "../conf/prod_conf.ini"
	}
	conf, _ := goconfig.LoadConfigFile(confPath)
	return &LocalConfig{conf: conf}
}

func (lconf *LocalConfig) GetFdfsDomain() (string, error) {
	return lconf.getValue("", "fdfsdomain")
}

func (lconf *LocalConfig) GetFdfsPort() int {
	return lconf.mustInt("", "fdfsport", 22122)
}

func (lconf *LocalConfig) GetDirPath(channel, storagetype string) (string, error) {
	//storagetype: nfs1,nfs2
	return lconf.getValue(channel, storagetype)
}

func (lconf *LocalConfig) GetResizeTypes(channel string) (string, error) {
	return lconf.getValue(channel, "resizetypes")
}
func (lconf *LocalConfig) GetSizes(channel string) (string, error) {
	return lconf.getValue(channel, "sizes")
}

func (lconf *LocalConfig) GetRotates(channel string) (string, error) {
	return lconf.getValue(channel, "rotates")
}

func (lconf *LocalConfig) GetQuality(channel string) (string, error) {
	return lconf.getValue(channel, "quality")
}
func (lconf *LocalConfig) GetQualities(channel string) (string, error) {
	return lconf.getValue(channel, "qualities")
}
func (lconf *LocalConfig) GetLogodir(channel string) (string, error) {
	return lconf.getValue(channel, "logodir")
}
func (lconf *LocalConfig) IsEnableNameLogo(channel string) (bool, error) {
	isenable, err := lconf.getValue(channel, "isenablenamelogo")
	if err != nil {
		return false, err
	}
	if isenable == "1" {
		return true, nil
	} else {
		return false, nil
	}
}
func (lconf *LocalConfig) GetDefaultLogo(channel string) (string, error) {
	return lconf.getValue(channel, "defaultlogo")
}
func (lconf *LocalConfig) GetLogoNames(channel string) (string, error) {
	return lconf.getValue(channel, "logonames")
}
func (lconf *LocalConfig) GetImagelesswidthForLogo(channel string) (int64, error) {
	width, err := lconf.getValue(channel, "imagelesswidthforlogo")
	if err != nil {
		return 0, err
	}
	if width == "" {
		return 0, nil
	}
	return strconv.ParseInt(width, 10, 64)
}
func (lconf *LocalConfig) GetImagelessheightForLogo(channel string) (int64, error) {
	height, err := lconf.getValue(channel, "imagelessheightforlogo")
	if err != nil {
		return 0, err
	}
	if height == "" {
		return 0, nil
	}
	return strconv.ParseInt(height, 10, 64)
}

func (lconf *LocalConfig) GetDissolves(channel string) (string, error) {
	return lconf.getValue(channel, "dissolves")
}
func (lconf *LocalConfig) GetDissolve(channel string) int {
	return lconf.mustInt(channel, "dissolve", 100)
}

func (lconf *LocalConfig) GetNamelogoDissolve(channel string) int {
	i, err := lconf.conf.Int(channel, "namelogodissolve")
	if err != nil {
		return 0
	} else {
		return i
	}
}
func (lconf *LocalConfig) GetSequenceofoperation(channel string) ([]string, error) {
	v, err := lconf.getValue(channel, "sequenceofoperation")
	if err != nil {
		return nil, err
	}
	if v == "" {
		v = "s,resize,q,m,rotate"
	}
	return strings.Split(v, ","), nil
}

func (lconf *LocalConfig) Reload() error {
	return lconf.conf.Reload()
}

func (lconf *LocalConfig) getValue(channel string, key string) (string, error) {
	v, _ := lconf.conf.GetValue(channel, key)
	if v == "" && channel != "" {
		v, _ = lconf.conf.GetValue("", key)
		return v, nil
	}
	if v == "nil" {
		v = ""
	}
	return v, nil
}

func (lconf *LocalConfig) mustInt(channel string, key string, defaultvalue int) int {
	//if err := loadConfiguration(); err != nil {
	//	return defaultvalue
	//}
	return lconf.conf.MustInt(channel, key, defaultvalue)
}

type RemoteConfig struct{}
