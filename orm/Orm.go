package orm

import (
	"EasyGo/config"
	"EasyGo/orm/MysqlOrm"
)

type Ormer struct {
	provider driver
}

func NewOrm(driver string, conf config.Result) Orm {
	o := &Ormer{}
	o.connect(driver, conf)
	return o
}

func (o *Ormer) connect(driver string, conf config.Result) {
	switch driver {
	case "mysql":
		o.provider = MysqlOrm.NewMysqlOrm(conf)
	}
}

func (o *Ormer) Insert(data interface{}) (int, error) {
	return o.provider.Insert(data)
}

func (o *Ormer) Update(data interface{}, condition ...interface{}) (int, error) {
	return o.provider.Update(data, condition...)
}

func (o *Ormer) Delete(condition ...interface{}) (int, error) {
	return o.provider.Delete(condition...)
}

func (o *Ormer) Fetch(condition ...interface{}) (interface{}, error) {
	return o.provider.Fetch(condition...)
}

func (o *Ormer) Query(cmd interface{}, params ...interface{}) (interface{}, error) {
	return o.provider.Query(cmd, params...)
}

func (o *Ormer) QueryMS() *MysqlOrm.MysqlOrm {
	if orm, ok := o.provider.(*MysqlOrm.MysqlOrm); ok {
		return orm
	}
	return nil
}
