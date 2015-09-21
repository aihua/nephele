package models

import (
	"github.com/astaxie/beego/orm"
)

type Channel struct {
	Name string
	Code string
}

func (this *Channel) Insert() error {
	o := orm.NewOrm()
	_, err := o.Raw("INSERT INTO channel(name,code) VALUES(?,?)", this.Name, this.Code).Exec()
	return err
}

func (this *Channel) Upload() error {
	o := orm.NewOrm()
	_, err := o.Raw("UPDATE FROM channel SET code=? WHERE name=?", this.Code, this.Name).Exec()
	return err
}

func (this *Channel) Get() ([]Channel, error) {
	var channels []Channel
	o := orm.NewOrm()

	_, err := o.Raw("SELECT name,code FROM channel").QueryRows(&channels)
	return channels, err
}
