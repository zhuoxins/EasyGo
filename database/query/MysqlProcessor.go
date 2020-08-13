package query

import (
	"EasyGo/config"
	"EasyGo/database"
	"database/sql"
)

type mysqlProcessor struct {
	clients *sql.DB
}

func (m *mysqlProcessor) Connect(conf config.Result) (err error) {
	client, err := database.Connect("mysql", conf)
	if err == nil {
		m.clients = client.(*sql.DB)
	}
	return
}

func (m *mysqlProcessor) Exec(scene int, SQL string, params ...interface{}) (result interface{}, err error) {
	var bindParam []interface{}
	if len(params) > 0 {
		bindParam = params[0].([]interface{})
	}
	switch scene {
	case FETCH:
		result, err = m.query(SQL, bindParam...)
	case FETCHRAW:
		result, err = m.queryRaw(SQL, bindParam...)
	default:
		result, err = m.exec(scene, SQL, bindParam...)
	}
	return
}

//执行sql语句(增删改)
func (m *mysqlProcessor) exec(scene int, SQL string, bindParams ...interface{}) (int, error) {
	res, err := m.clients.Exec(SQL, bindParams...)
	if err != nil {
		return 0, err
	}
	var result int64
	if scene == INSERT {
		result, err = res.LastInsertId()
	} else {
		result, err = res.RowsAffected()
	}
	if err != nil {
		return 0, err
	}
	return int(result), nil
}

func (m *mysqlProcessor) query(SQL string, bindParams ...interface{}) ([]map[string]string, error) {
	stmtOut, err := m.clients.Prepare(SQL)
	if err != nil {
		return nil, err
	}
	rows, err := stmtOut.Query(bindParams...)
	if err != nil {
		return nil, err
	}

	defer stmtOut.Close()
	defer rows.Close()
	return m.collectResult(rows)
}

//普通sql查询
func (m *mysqlProcessor) queryRaw(SQL string, args ...interface{}) ([]map[string]string, error) {
	rows, err := m.clients.Query(SQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return m.collectResult(rows)
}

//收集查询结果
func (m *mysqlProcessor) collectResult(rows *sql.Rows) (result []map[string]string, err error) {
	result = make([]map[string]string, 0)
	for rows.Next() {
		//获取字段
		cols, err := rows.Columns()
		if err != nil {
			return result, err
		}
		//定义查询字段切片
		fieldSlice := make([][]byte, len(cols))
		scanSlice := make([]interface{}, 0)
		for k := range fieldSlice {
			scanSlice = append(scanSlice, &fieldSlice[k])
		}
		err = rows.Scan(scanSlice...)
		if err != nil {
			return result, err
		}
		//查询结果存到map
		row := make(map[string]string)
		for k, v := range fieldSlice {
			key := cols[k]
			row[key] = string(v)
		}
		result = append(result, row)
	}
	return result, err
}

func (m *mysqlProcessor) Close() error {
	return m.clients.Close()
}
