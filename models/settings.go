package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	//Sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

type Settings struct {
	Id      int64
	Profile string `orm:"size(64);unique" form:"Profile" valid:"Required;"`

	MIAddress string `orm:"size(64);unique" form:"MIAddress" valid:"Required;"`
	MINetwork string `orm:"size(64);unique" form:"MINetwork" valid:"Required;"`

	OVConfigPath string `orm:"size(64);unique" form:"OVConfigPath" valid:"Required;"`

	OVPkiPath string `orm:"size(64);unique" form:"OVPkiPath" valid:"Required;"`

	ServerAddress string `orm:"size(64);unique" form:"ServerAddress" valid:"Required;"`

	ServerName string `orm:"size(64);unique" form:"ServerName" valid:"Required;"`

	OVEasyRsaPath string `orm:"size(64);unique" form:"OVEasyRsaPath" valid:"Required;"`

	OVExtraFiles string `orm:"size(256);unique" form:"OVExtraFiles" valid:"Optional;"`

	SysAdmin string `orm:"size(64);unique" form:"SysAdmin" valid:"Required;"`

	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
}

//Insert wrapper
func (s *Settings) Insert() error {
	if _, err := orm.NewOrm().Insert(s); err != nil {
		return err
	}
	return nil
}

//Read wrapper
func (s *Settings) Read(fields ...string) error {
	if err := orm.NewOrm().Read(s, fields...); err != nil {
		return err
	}
	return nil
}

//Update wrapper
func (s *Settings) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(s, fields...); err != nil {
		return err
	}
	return nil
}

//Delete wrapper
func (s *Settings) Delete() error {
	if _, err := orm.NewOrm().Delete(s); err != nil {
		return err
	}
	return nil
}
