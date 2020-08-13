package MysqlOrm

import "strconv"

const (
	INSERT = iota
	DELETE
	UPDATE
	FETCH
	FETCHRAW
	WHERERAW
	HAVINGRAW
	ORDERBYRAW
	INNERJOIN
	LEFTJOIN
	RIGHTJOIN
	COUNT
	MAX
	MIN
	AVG
	SUM
	QUERY_SCENE
	EXEC_SCENE
	TRANSACTION
	NORMAL
)

//查询条件
type conditions map[string]interface{}
type OrQuery func(o *OrWhere)

func ToInt(r string) int {
	if i, ok := strconv.Atoi(r); ok == nil {
		return i
	}
	return 0
}
