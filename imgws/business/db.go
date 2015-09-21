package business

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/ctripcorp/nephele/util"
	"time"
)

var (
	getChannelsTime time.Time         = time.Now()
	channels        map[string]string = make(map[string]string)

	ERRORTYPE_GETCHANNEL = "GetChannel"
)

func GetChannels() (map[string]string, util.Error) {
	if len(channels) < 1 || IsRefresh(getChannelsTime) {
		o := orm.NewOrm()
		res := make(orm.Params)
		_, err := o.Raw("SELECT name,code FROM channel").RowsToMap(&res, "name", "code")
		if err != nil {
			return nil, util.Error{IsNormal: false, Err: err, Type: ERRORTYPE_GETCHANNEL}
		}
		m := make(map[string]string)
		for k, v := range res {
			m[k] = fmt.Sprintf("%v", v)
		}
		channels = m
		getChannelsTime = time.Now()
	}
	return channels, util.Error{}
}

func GetChannelCode(channel string) string {
	channels, e := GetChannels()
	if e.Err != nil {
		return ""
	}
	code, _ := channels[channel]
	return code
}

func IsRefresh(t time.Time) bool {
	return t.Add(1 * time.Minute).Before(time.Now())
}
