package MysqlOrm

import (
	"EasyGo/config"
	"EasyGo/database/query"
	"EasyGo/tools/helper"
	"errors"
	"math"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type condition map[string]interface{}

type MysqlOrm struct {
	conf            config.Result   //连接配置
	provider        query.Processor //处理器
	tableName       string          //表名
	prefix          string          //表前缀
	pk              string          //主键
	allFields       []string        //表所有字段
	fields          []string        //查询字段信息
	fieldRaw        string          //查询字段
	distincts       []string        //不重复字段
	aggField        string          //聚合字段
	condition       condition       //where and条件
	orCondition     []*OrWhere      //where or 条件
	bindParams      []interface{}   //预处理绑定值
	aliasName       string          //表别名
	joinRaw         []string        //连表信息
	groupRaw        string          //分组信息
	havConditions   [][]interface{} //分组查询条件
	orHavConditions []*OrWhere
	orderRaw        string        //排序信息
	limitRaw        string        //截取条数
	querySql        string        //最后执行sql
	status          int           //事务状态
	whereRaw        []interface{} //where原生语句
	havingRaw       []interface{} //having原生语句
	orderByRaw      []interface{} //orderBy原生语句
	selectRaw       string        //查询原生sql语句
	isLock          bool          //是否开启锁
	noClear         bool
}

func NewMysqlOrm(conf config.Result) *MysqlOrm {
	m := &MysqlOrm{
		condition: make(condition),
		pk:        "id",
		status:    NORMAL,
		conf:      conf,
	}
	return m
}

func (this *MysqlOrm) ResetConnection(conf config.Result) {
	this.conf = conf
	this.condition = make(condition)
	this.pk = "id"
	this.status = NORMAL
}

func (this *MysqlOrm) connect(conf config.Result) (err error) {
	provider := query.NewProcessor("mysql")
	err = provider.Connect(conf)
	if err == nil {
		this.provider = provider
	}
	return
}

func (this *MysqlOrm) SetTabPrefix(prefix string) *MysqlOrm {
	this.prefix = prefix
	return this
}

func (this *MysqlOrm) SetPk(pk string) *MysqlOrm {
	this.pk = pk
	return this
}

func (this *MysqlOrm) Table(tableName string) *MysqlOrm {
	this.tableName = tableName
	return this
}

func (this *MysqlOrm) table() string {
	return this.prefix + this.tableName
}

/**
可变where条件
field, value
field, condition, value
*/
func (this *MysqlOrm) Where(field string, vars ...interface{}) *MysqlOrm {
	if len(vars) > 0 {
		this.condition[field] = vars
	}
	return this
}

/**
where条件集合
field => value
field => []{condition, value}
*/
func (this *MysqlOrm) WhereMap(maps map[string]interface{}) *MysqlOrm {
	if len(maps) > 0 {
		for field, condition := range maps {
			this.condition[field] = condition
		}
	}
	return this
}

func (this *MysqlOrm) OrWhere(query OrQuery) *MysqlOrm {
	orWhere := &OrWhere{}
	query(orWhere)
	this.orCondition = append(this.orCondition, orWhere)
	return this
}

func (this *MysqlOrm) WhereRaw(raw string, args ...interface{}) *MysqlOrm {
	return this.saveRaw(WHERERAW, raw, args...)
}

//in 条件查询
func (this *MysqlOrm) WhereIn(field string, vars ...interface{}) *MysqlOrm {
	if length := len(vars); length > 1 {
		condition := []interface{}{"IN", vars}
		this.condition[field] = condition
	} else if length == 1 {
		condition := []interface{}{"IN", vars[0]}
		this.condition[field] = condition
	}
	return this
}

func (this *MysqlOrm) Field(fields ...string) *MysqlOrm {
	if len(fields) > 0 {
		for _, field := range fields {
			this.fields = append(this.fields, strings.TrimSpace(field))
		}
	}
	return this
}

func (this *MysqlOrm) FieldRaw(raw string) *MysqlOrm {
	this.fieldRaw = raw
	return this
}

func (this *MysqlOrm) Distinct(fields ...string) *MysqlOrm {
	if len(fields) > 0 {
		for _, v := range fields {
			this.distincts = append(this.distincts, strings.TrimSpace(v))
		}
	}
	return this
}

//原生语句查询
func (this *MysqlOrm) SelectRaw(SQL string, args ...interface{}) *MysqlOrm {
	if len(args) > 0 {
		for _, vars := range args {
			this.bindParams = append(this.bindParams, vars)
		}
	}
	this.selectRaw = SQL
	return this
}

//关联查询
func (this *MysqlOrm) Join(joinTab, leftHand, condition, rightHand string) *MysqlOrm {
	return this.pushJoinRaw(INNERJOIN, joinTab, leftHand, condition, rightHand)
}

func (this *MysqlOrm) LeftJoin(joinTab, leftHand, condition, rightHand string) *MysqlOrm {
	return this.pushJoinRaw(LEFTJOIN, joinTab, leftHand, condition, rightHand)
}

func (this *MysqlOrm) RightJoin(joinTab, leftHand, condition, rightHand string) *MysqlOrm {
	return this.pushJoinRaw(RIGHTJOIN, joinTab, leftHand, condition, rightHand)
}

func (this *MysqlOrm) pushJoinRaw(scene int, joinTab string, arg ...string) *MysqlOrm {
	joinTab = strings.TrimSpace(joinTab)
	if this.prefix != "" && !strings.Contains(joinTab, this.prefix) {
		joinTab = this.prefix + joinTab
	}
	var joinCond, raw string
	for _, vars := range arg {
		joinCond += strings.TrimSpace(vars) + " "
	}
	joinCond = strings.TrimRight(joinCond, " ")

	switch scene {
	case INNERJOIN:
		raw = " INNER JOIN " + joinTab + " ON " + joinCond
	case LEFTJOIN:
		raw = " LEFT JOIN " + joinTab + " ON " + joinCond
	case RIGHTJOIN:
		raw = " RIGHT JOIN " + joinTab + " ON " + joinCond
	}
	this.joinRaw = append(this.joinRaw, raw)
	return this
}

//排序
func (this *MysqlOrm) OrderBy(field string, rule ...string) *MysqlOrm {
	var sortRule string
	if len(rule) > 0 {
		sortRule = strings.ToUpper(strings.TrimSpace(rule[0]))
	} else {
		sortRule = "ASC"
	}
	if sortRule != "ASC" && sortRule != "DESC" {
		return this
	}
	this.orderRaw = " ORDER BY " + field + " " + sortRule
	return this
}

func (this *MysqlOrm) OrderByRaw(raw string, args ...interface{}) *MysqlOrm {
	return this.saveRaw(ORDERBYRAW, raw, args...)
}

//分组查询
func (this *MysqlOrm) GroupBy(fields ...string) *MysqlOrm {
	var raw string
	if len(fields) > 0 {
		for _, vars := range fields {
			raw += strings.TrimSpace(vars) + ","
		}
		raw = strings.TrimRight(raw, ", ")
		this.groupRaw = " GROUP BY " + raw
	}
	return this
}

//分组排序
func (this *MysqlOrm) Having(field, condition string, val interface{}) *MysqlOrm {
	cond := []interface{}{field, condition, val}
	this.havConditions = append(this.havConditions, cond)
	return this
}

func (this *MysqlOrm) OrHaving(query OrQuery) *MysqlOrm {
	orWhere := &OrWhere{}
	query(orWhere)
	this.orHavConditions = append(this.orHavConditions, orWhere)
	return this
}

func (this *MysqlOrm) HavingRaw(raw string, args ...interface{}) *MysqlOrm {
	return this.saveRaw(HAVINGRAW, raw, args...)
}

func (this *MysqlOrm) saveRaw(scene int, raw string, args ...interface{}) *MysqlOrm {
	switch scene {
	case WHERERAW:
		this.whereRaw = []interface{}{raw, args}
	case HAVINGRAW:
		this.havingRaw = []interface{}{raw, args}
	case ORDERBYRAW:
		this.orderByRaw = []interface{}{raw, args}
	}
	return this
}

//条数限制
func (this *MysqlOrm) Limit(start, end int) *MysqlOrm {
	this.limitRaw = " LIMIT " + strconv.Itoa(start) + "," + strconv.Itoa(end)
	return this
}

/**
聚合查询
*/

//查询数量
func (this *MysqlOrm) Count(field ...string) (string, error) {
	if len(field) > 0 {
		return this.aggregate(field[0], COUNT)
	} else {
		return this.aggregate("*", COUNT)
	}
}

//查询最大值
func (this *MysqlOrm) Max(field string) (string, error) {
	return this.aggregate(field, MAX)
}

//查询数量
func (this *MysqlOrm) Min(field string) (string, error) {
	return this.aggregate(field, MIN)
}

//查询数量
func (this *MysqlOrm) Avg(field string) (string, error) {
	return this.aggregate(field, AVG)
}

//查询数量
func (this *MysqlOrm) Sum(field string) (string, error) {
	return this.aggregate(field, SUM)
}

//执行聚合函数
func (this *MysqlOrm) aggregate(field string, scene int) (string, error) {
	switch scene {
	case COUNT:
		this.aggField = " count(" + strings.TrimSpace(field) + ") as aggres"
	case AVG:
		this.aggField = " avg(" + strings.TrimSpace(field) + ") as aggres"
	case MAX:
		this.aggField = " max(" + strings.TrimSpace(field) + ") as aggres"
	case MIN:
		this.aggField = " min(" + strings.TrimSpace(field) + ") as aggres"
	case SUM:
		this.aggField = " sum(" + strings.TrimSpace(field) + ") as aggres"
	}
	result, err := this.query(FETCH, this.buildingSearch())
	if err == nil {
		res, _ := result[0]["aggres"]
		return res, err
	} else {
		return "", err
	}
}

func (this *MysqlOrm) Paginate(page, pageSize int) (map[string]interface{}, error) {
	//计算总条数
	this.noClear = true
	total, err := this.Count()
	if ToInt(total) == 0 || err != nil {
		return nil, err
	}
	this.aggField = ""
	this.bindParams = []interface{}{}
	result := make(map[string]interface{})
	maxPage := math.Ceil(float64(ToInt(total)) / float64(pageSize))
	page = int(math.Max(1, float64(page)))
	page = int(math.Min(float64(page), maxPage))
	//获取偏移量
	offset := (page - 1) * pageSize
	list, err := this.Limit(offset, pageSize).Get()
	if err != nil {
		return nil, err
	}
	result = map[string]interface{}{
		"total":   total,
		"page":    page,
		"maxPage": maxPage,
		"data":    list,
	}
	return result, err
}

func (this *MysqlOrm) LockForUpdate() *MysqlOrm {
	this.isLock = true
	return this
}

//获取数据表全部字段
func (this *MysqlOrm) getFields() error {
	SQL := "DESC " + this.table()
	result, err := this.query(FETCHRAW, SQL)
	if err != nil {
		return err
	}
	for _, fieldInfo := range result {
		this.allFields = append(this.allFields, fieldInfo["Field"])
	}
	return nil
}

//构建查询语句
func (this *MysqlOrm) buildingSearch() string {
	SQL := "SELECT"
	if this.aggField != "" {
		//拼接聚合查询
		SQL = this.disposeAggregate(SQL)
	} else if this.fieldRaw != "" {
		SQL += " " + this.fieldRaw
	} else if len(this.distincts) > 0 {
		//拼接不重复字段
		SQL = this.disposeDistinct(SQL)
	} else {
		//拼接查询字段
		SQL = this.disposeFields(SQL)
	}
	//拼接表名
	SQL = this.jointTab(SQL)
	//拼接关联查询
	SQL = this.disposeJoin(SQL)
	//拼接查询条件
	SQL = this.disposeWhereRaw(SQL)
	//拼接分组条件
	SQL = this.disposeGroup(SQL)
	//分组查询条件
	SQL = this.disposeHaving(SQL)
	//排序
	SQL = this.disposeOrder(SQL)

	SQL = this.disposeLimit(SQL)

	if this.isLock {
		SQL += " for update"
	}

	return SQL
}

func (this *MysqlOrm) disposeDistinct(SQL string) string {
	if len(this.distincts) > 0 {
		SQL += " DISTINCT "
		for _, field := range this.distincts {
			SQL += strings.TrimSpace(field) + ", "
		}
		SQL = strings.TrimRight(SQL, ", ")
	}
	return SQL
}

func (this *MysqlOrm) disposeFields(SQL string) string {
	if len(this.fields) > 0 {
		for _, field := range this.fields {
			SQL += " " + field + ","
		}
		SQL = strings.TrimRight(SQL, ",")
	} else if len(this.distincts) == 0 && len(this.fields) == 0 {
		SQL += " *"
	}
	SQL += " FROM"
	return SQL
}

func (this *MysqlOrm) disposeAggregate(SQL string) string {
	if this.aggField != "" {
		SQL += this.aggField
	}
	return SQL
}

func (this *MysqlOrm) jointTab(SQL string) string {
	if strings.Contains(SQL, "FROM") || strings.Contains(SQL, "from") {
		SQL += " " + this.table()
	} else {
		SQL += " FROM " + this.table()
	}
	if this.aliasName != "" {
		SQL += " AS " + this.aliasName
	}
	return SQL
}

func (this *MysqlOrm) disposeJoin(SQL string) string {
	if len(this.joinRaw) > 0 {
		for _, raw := range this.joinRaw {
			SQL += raw
		}
	}
	return SQL
}

//拼接where条件
func (this *MysqlOrm) disposeWhereRaw(SQL string) string {
	if length := len(this.whereRaw); length > 0 {
		SQL += " where " + this.whereRaw[0].(string)
		if length > 1 {
			args := this.whereRaw[1].([]interface{})
			for _, v := range args {
				this.bindParams = append(this.bindParams, v)
			}
		}
		return SQL
	}
	var whereRaw string
	if len(this.condition) > 0 {
		whereRaw = " WHERE "
		for k, val := range this.condition {
			s := this.jointCondition(k, val)
			whereRaw += s + " AND "
		}
		whereRaw = strings.TrimRight(whereRaw, "AND ")
	}
	if len(this.orCondition) > 0 {
		for _, query := range this.orCondition {
			if whereRaw == "" {
				whereRaw = " where "
			} else {
				whereRaw += " and ("
			}
			itemRaw, binds := query.getRaw()
			if len(binds) > 0 {
				for _, vars := range binds {
					this.bindParams = append(this.bindParams, vars)
				}
			}
			whereRaw += whereRaw + itemRaw + ")"
		}
	}
	if whereRaw != "" {
		SQL += whereRaw
	}
	return SQL
}

//处理where子句
func (this *MysqlOrm) jointCondition(field string, val interface{}) string {
	var raw string
	var params interface{}
	switch val.(type) {
	case []interface{}:
		conds := val.([]interface{})
		params = conds[1]
		condition := strings.TrimSpace(conds[0].(string))
		if condition == "in" {
			raw = field + " " + condition + " ("
			for _, param := range params.([]interface{}) {
				raw += "?, "
				this.bindParams = append(this.bindParams, param)
			}
			raw = strings.TrimRight(raw, ", ") + ")"
		} else {
			raw = strings.TrimSpace(field) + " " + condition + " ?"
			this.bindParams = append(this.bindParams, params)
		}
	default:
		params = val
		raw = strings.TrimSpace(field) + " = ?"
		this.bindParams = append(this.bindParams, params)
	}
	return raw
}

func (this *MysqlOrm) disposeGroup(SQL string) string {
	if this.groupRaw != "" {
		SQL += this.groupRaw
	}
	return SQL
}

func (this *MysqlOrm) disposeHaving(SQL string) string {
	if length := len(this.havingRaw); length > 0 {
		SQL += " HAVING " + this.havingRaw[0].(string)
		if length > 1 {
			args := this.havingRaw[1].([]interface{})
			for _, v := range args {
				this.bindParams = append(this.bindParams, v)
			}
		}
		return SQL
	}
	var havingRaw string
	if len(this.havConditions) > 0 {
		havingRaw = " HAVING "
		for _, val := range this.havConditions {
			s := this.jointCondition(val[0].(string), val[1:])
			havingRaw += s + " AND "
		}
		havingRaw = strings.TrimRight(havingRaw, "AND ")
	}
	if len(this.orHavConditions) > 0 {
		for _, query := range this.orHavConditions {
			if havingRaw == "" {
				havingRaw = " HAVING "
			} else {
				havingRaw += " AND ("
			}
			itemRaw, binds := query.getRaw()
			if len(binds) > 0 {
				for _, vars := range binds {
					this.bindParams = append(this.bindParams, vars)
				}
			}
			havingRaw += havingRaw + itemRaw + ")"
		}
	}
	if havingRaw != "" {
		SQL += havingRaw
	}
	return SQL
}

func (this *MysqlOrm) disposeOrder(SQL string) string {
	if length := len(this.orderByRaw); length > 0 {
		SQL += " ORDER BY " + this.orderByRaw[0].(string)
		if length > 1 {
			args := this.orderByRaw[1].([]interface{})
			for _, v := range args {
				this.bindParams = append(this.bindParams, v)
			}
		}
		return SQL
	}
	if this.orderRaw != "" {
		SQL += this.orderRaw
	}
	return SQL
}

func (this *MysqlOrm) disposeLimit(SQL string) string {
	if this.limitRaw != "" {
		SQL += this.limitRaw
	}
	return SQL
}

//查询
func (this *MysqlOrm) Get() ([]map[string]string, error) {
	result, err := this.query(FETCH, this.buildingSearch())
	return result, err
}

func (this *MysqlOrm) Fetch(maps ...interface{}) (interface{}, error) {
	if len(maps) > 0 {
		this.WhereMap(maps[0].(map[string]interface{}))
	}
	return this.Get()
}

//查询单条
func (this *MysqlOrm) Find(pks ...int) (map[string]string, error) {
	if len(pks) > 0 {
		this.condition[this.pk] = pks[0]
	}
	SQL := this.buildingSearch()
	if this.limitRaw == "" && !this.isLock {
		SQL += " Limit 0,1"
	}
	result, err := this.query(FETCH, SQL)
	data := make(map[string]string)
	if err == nil && len(result) > 0 {
		data = result[0]
		return data, err
	} else {
		return data, err
	}
}

//返回指定字段值
func (this *MysqlOrm) Value(column string) (string, error) {
	data, err := this.Find()
	if err == nil {
		if res, ok := data[column]; ok {
			return res, err
		}
	}
	return "", err
}

//添加单条数据
func (this *MysqlOrm) Insert(data interface{}) (int, error) {
	insertData := data.(map[string]interface{})
	err := this.getFields()
	if err != nil {
		return 0, err
	}
	SQL := "INSERT INTO " + this.table()
	valRaw := "VALUES("
	fields := ""
	for k, v := range insertData {
		if !helper.InStringArray(k, this.allFields) {
			continue
		} else {
			fields += k + ","
			valRaw += "?,"
			this.bindParams = append(this.bindParams, v)
		}
	}
	fields = strings.TrimRight(fields, ",")
	valRaw = strings.TrimRight(valRaw, ",")
	SQL += "(" + fields + ") " + valRaw + ")"
	result, err := this.exec(INSERT, SQL)
	return result, err
}

func (this *MysqlOrm) Delete(maps ...interface{}) (int, error) {
	if len(maps) > 0 {
		this.WhereMap(maps[0].(map[string]interface{}))
	}
	SQL := "DELETE FROM" + this.table()
	SQL = this.disposeWhereRaw(SQL)
	return this.exec(DELETE, SQL)
}

//修改
func (this *MysqlOrm) Update(changeData interface{}, maps ...interface{}) (int, error) {
	data := changeData.(map[string]interface{})
	if len(data) <= 0 {
		return 0, nil
	}
	if len(maps) > 0 {
		this.WhereMap(maps[0].(map[string]interface{}))
	}
	SQL := "UPDATE " + this.table() + " SET "
	var raw string
	for k, v := range data {
		this.bindParams = append(this.bindParams, v)
		raw += k + "=?,"
	}
	raw = strings.TrimRight(raw, ",")
	SQL = this.disposeWhereRaw(SQL + raw)
	result, err := this.exec(UPDATE, SQL)
	return result, err
}

//开启事务
func (this *MysqlOrm) Begin() error {
	_, err := this.queryRaw("SET AUTOCOMMIT = 0")
	if err != nil {
		return err
	}
	_, err = this.queryRaw("BEGIN")
	return err
}

func (this *MysqlOrm) Rollback() error {
	_, err := this.queryRaw("ROLLBACK")
	return err
}

func (this *MysqlOrm) Commit() error {
	_, err := this.queryRaw("COMMIT")
	return err
}

//事务操作 自动提交回滚
func (this *MysqlOrm) Transaction(exec func()) error {
	defer this.dealTrans()
	this.status = TRANSACTION
	err := this.Begin()
	if err != nil {
		return err
	}
	exec()
	return nil
}

func (this *MysqlOrm) dealTrans() {
	this.status = NORMAL
	err := recover()
	msg := ""
	switch err.(type) {
	case runtime.Error:
		msg = err.(runtime.Error).Error()
	default:
		if reflect.ValueOf(err).String() != reflect.ValueOf(nil).String() {
			msg = err.(string)
		}
	}
	if msg != "" {
		this.Rollback()
	} else {
		this.Commit()
	}
}

func (this *MysqlOrm) Query(cmd interface{}, params ...interface{}) (interface{}, error) {
	return this.queryRaw(cmd.(string), params...)
}

func (this *MysqlOrm) queryRaw(SQL string, params ...interface{}) (interface{}, error) {
	err := this.prepare()
	if err != nil {
		return nil, err
	}
	return this.provider.Exec(FETCHRAW, SQL, params...)
}

//执行查询操作
func (this *MysqlOrm) query(scene int, SQL string) ([]map[string]string, error) {
	if !this.noClear {
		defer this.resetCondition()
	} else {
		this.noClear = false
	}
	err := this.prepare()
	if err != nil {
		return nil, err
	}
	if this.isLock {
		this.isLock = false
		_, err := this.queryRaw(this.jointSql(SQL))
		if err != nil && this.status == TRANSACTION {
			this.throw(err)
		}
		return nil, err
	}
	this.querySql = SQL
	result, err := this.provider.Exec(scene, SQL, this.bindParams)
	if err == nil {
		return result.([]map[string]string), err
	}
	return nil, err
}

func (this *MysqlOrm) jointSql(SQL string) string {
	if len(this.bindParams) > 0 {
		parseSql := strings.Split(SQL, "?")
		var newSql string
		for k, param := range this.bindParams {
			if val, ok := param.(string); ok {
				newSql += parseSql[k] + " " + val + " "
			} else if val, ok := param.(int); ok {
				s := strconv.Itoa(val)
				newSql += parseSql[k] + " " + s + " "
			} else if val, ok := param.(float64); ok {
				s := strconv.FormatFloat(val, 'f', -1, 64)
				newSql += parseSql[k] + " " + s + " "
			} else if val, ok := param.(int64); ok {
				s := strconv.FormatInt(val, 10)
				newSql += parseSql[k] + " " + s + " "
			}
		}
		newSql += parseSql[len(parseSql)-1]
		return newSql
	}
	return SQL
}

//执行增删改操作
func (this *MysqlOrm) exec(scene int, SQL string) (int, error) {
	defer this.resetCondition()
	err := this.prepare()
	if err != nil {
		return 0, err
	}
	this.querySql = SQL
	result, err := this.provider.Exec(scene, SQL, this.bindParams)
	if err == nil {
		return result.(int), err
	}
	if this.status == TRANSACTION {
		this.throw(err)
	}
	return 0, err
}

func (this *MysqlOrm) prepare() error {
	if this.conf == nil && this.provider == nil {
		return errors.New("mysql client error")
	}
	if this.provider == nil {
		err := this.connect(this.conf)
		if err != nil {
			return err
		}
	}
	return nil
}

//手动宕机
func (this *MysqlOrm) throw(err error) {
	panic(err.Error())
}

func (this *MysqlOrm) GetLastSql() string {
	return this.querySql
}

//重置查询条件
func (this *MysqlOrm) resetCondition() {
	this.aliasName = ""
	this.joinRaw = make([]string, 0)
	this.orderRaw = ""
	this.limitRaw = ""
	this.condition = make(condition)
	this.bindParams = make([]interface{}, 0)
	this.fields = make([]string, 0)
	this.orCondition = make([]*OrWhere, 0)
	this.havConditions = make([][]interface{}, 0)
	this.orHavConditions = make([]*OrWhere, 0)
	this.groupRaw = ""
	this.distincts = make([]string, 0)
	this.aggField = ""
	this.orderByRaw = make([]interface{}, 0)
	this.whereRaw = make([]interface{}, 0)
	this.havingRaw = make([]interface{}, 0)
	this.fieldRaw = ""
}

func (this *MysqlOrm) Close() error {
	return this.provider.Close()
}
