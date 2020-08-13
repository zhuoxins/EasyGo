package model

import (
	"EasyGo/config"
	"EasyGo/orm"
	"EasyGo/orm/MysqlOrm"
)

type Model struct {
	TableName  string
	Pk         string
	Connection string
	AutoTime   bool
	CreateAt   string
	UpdateAt   string
	PageLimit  int
	orm        orm.Orm
}

func (m *Model) Init() {
	if m.Connection == "" {
		m.Connection = "mysql"
	}
	conf := config.Get("database." + m.Connection)
	m.orm = orm.NewOrm(m.Connection, conf)
	if m.Pk == "" {
		m.Pk = "id"
	}
	if m.Connection == "mysql" {
		m.orm.QueryMS().SetPk(m.Pk)
		m.orm.QueryMS().Table(m.TableName)
		prefix := conf.GetField("prefix", "").String()
		if prefix != "" {
			m.orm.QueryMS().SetTabPrefix(prefix)
		}
	}
	if m.AutoTime == true && m.CreateAt == "" {
		m.CreateAt = "create_time"
	}
	if m.AutoTime == true && m.UpdateAt == "" {
		m.CreateAt = "update_time"
	}
	if m.PageLimit == 0 {
		m.PageLimit = 10
	}
}

func (m *Model) Insert(data interface{}) (int, error) {
	return m.orm.Insert(data)
}

func (m *Model) Delete(condition ...interface{}) (int, error) {
	return m.orm.Delete(condition...)
}

func (m *Model) Update(data interface{}, condition ...interface{}) (int, error) {
	return m.orm.Update(data, condition...)
}

func (m *Model) Get(condition ...interface{}) (interface{}, error) {
	return m.orm.Fetch(condition...)
}

func (m *Model) Query(cmd interface{}, params ...interface{}) (interface{}, error) {
	return m.orm.Query(cmd, params...)
}

func (m *Model) QueryMS() *MysqlOrm.MysqlOrm {
	return m.orm.QueryMS()
}

//func (m *Model) SetTable(tab string) {
//	m.tableName = tab
//}
