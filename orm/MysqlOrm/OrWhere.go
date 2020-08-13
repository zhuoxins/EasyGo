package MysqlOrm

type OrWhere struct {
	condition [][]interface{}
	bindParam []interface{}
}

func (o *OrWhere) where(field string, vars ...interface{}) *OrWhere {
	o.condition = append(o.condition, []interface{}{"where", field, vars})
	return o
}

func (o *OrWhere) OrWhere(field string, vars ...interface{}) *OrWhere {
	o.condition = append(o.condition, []interface{}{"or", field, vars})
	return o
}

func (o *OrWhere) getRaw() (string, []interface{}) {
	var raw, scene string
	for _, item := range o.condition {
		scenes := item[0].(string)
		if scene == "" {
			raw = "("
		} else if scene == scenes {
			raw += " AND "
		} else if scene != scenes {
			if scenes == "where" {
				raw += ") AND ("
			} else {
				raw += ") OR ("
			}
		}
		scene = scenes
		itemField := item[1].(string)
		itemCond := item[2].([]interface{})
		itemRaw := o.jointWhere(itemField, itemCond)
		raw += itemRaw
	}
	raw += ")"
	return raw, o.bindParam
}

func (o *OrWhere) jointWhere(field string, whereCond []interface{}) string {
	var raw string
	if length := len(whereCond); length == 1 {
		raw += field + " = ?"
		o.bindParam = append(o.bindParam, whereCond[0])
	} else if length == 2 {
		raw += field + " " + whereCond[0].(string) + " ?"
		o.bindParam = append(o.bindParam, whereCond[1])
	}
	return raw
}
