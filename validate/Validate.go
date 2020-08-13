package validate

import (
	"EasyGo/tools/helper"
	"strconv"
	"strings"
)

type Validate struct {
	process map[string]verifyFunc
	ruleSet map[string]map[string]map[string]string
}

func NewValidate() *Validate {
	v := &Validate{}
	v.setProcess()
	return v
}

func (v *Validate) reset() {
	v.ruleSet = make(map[string]map[string]map[string]string)
}

func (v *Validate) setProcess() {
	v.process = make(map[string]verifyFunc)
	v.process = map[string]verifyFunc{
		"required":  v.isRequire,
		"max":       v.max,
		"min":       v.min,
		"maxLength": v.maxLength,
		"numeric":   v.isNumeric,
		"int":       v.isInt,
		"float":     v.isFloat,
		"email":     v.isEmail,
		"ip":        v.isIp,
		"phone":     v.isPhone,
	}
}

//验证required规则,返回是否验证通过和错误信息
func (v *Validate) isRequire(field string, val string, rules map[string]string) (bool, string) {
	msg, res := "", true
	if val == "" {
		res = false
		if rules["msg"] != "" {
			msg = rules["msg"]
		} else {
			//没有自定义错误提示
			msg = "param " + field + " is required"
		}
	}
	return res, msg
}

//验证数字最大值
func (v *Validate) max(field string, val string, rules map[string]string) (bool, string) {
	//转换类型
	i, err := helper.ToInt(val)
	if err != nil {
		return false, "param " + field + " is not a number"
	}
	msg, res := "", true
	max, _ := helper.ToInt(rules["ruleParam"])
	if i > max {
		res = false
		if rules["msg"] != "" {
			msg = rules["msg"]
		} else {
			//没有自定义错误提示
			msg = "param " + field + " is too large"
		}
	}
	return res, msg
}

//验证数字最小值
func (v *Validate) min(field string, val string, rules map[string]string) (bool, string) {
	//转换类型
	i, err := helper.ToInt(val)
	if err != nil {
		return false, "param " + field + " is not a number"
	}
	msg, res := "", true
	min, _ := helper.ToInt(rules["ruleParam"])
	if i < min {
		res = false
		if rules["msg"] != "" {
			msg = rules["msg"]
		} else {
			//没有自定义错误提示
			msg = "param " + field + " is too small"
		}
	}
	return res, msg
}

//验证字符最大长度
func (v *Validate) maxLength(field string, val string, rules map[string]string) (bool, string) {
	msg, res := "", true
	maxLen, _ := helper.ToInt(rules["ruleParam"])
	if helper.StrLength(val) > maxLen {
		res = false
		if rules["msg"] != "" {
			msg = rules["msg"]
		} else {
			//没有自定义错误提示
			msg = "param " + field + " is too long"
		}
	}
	return res, msg
}

func (v *Validate) isNumeric(field string, val string, rules map[string]string) (bool, string) {
	msg, res := "", false
	//转int
	if _, err := strconv.Atoi(val); err == nil {
		res = true
	}
	//转float
	if !res {
		if _, err := strconv.ParseFloat(val, 64); err == nil {
			res = true
		}
	}
	if !res {
		if _, err := strconv.ParseFloat(val, 32); err == nil {
			res = true
		}
	}
	if !res {
		if rules["msg"] != "" {
			msg = rules["msg"]
		} else {
			//没有自定义错误提示
			msg = "param " + field + " not numeric"
		}
	}
	return res, msg
}

func (v *Validate) isInt(field string, val string, rules map[string]string) (bool, string) {
	msg, res := "", false
	if _, err := strconv.Atoi(val); err == nil {
		res = true
	}
	if !res {
		if rules["msg"] != "" {
			msg = rules["msg"]
		} else {
			//没有自定义错误提示
			msg = "param " + field + " not int"
		}
	}
	return res, msg
}

func (v *Validate) isFloat(field string, val string, rules map[string]string) (bool, string) {
	msg, res := "", false
	if _, err := strconv.ParseFloat(val, 64); err == nil {
		res = true
	}
	if !res {
		if _, err := strconv.ParseFloat(val, 32); err == nil {
			res = true
		}
	}
	if !res {
		if rules["msg"] != "" {
			msg = rules["msg"]
		} else {
			//没有自定义错误提示
			msg = "param " + field + " not float"
		}
	}
	return res, msg
}

func (v *Validate) isEmail(field string, val string, rules map[string]string) (bool, string) {
	msg, res := "", false
	if emailPattern.MatchString(val) {
		res = true
	}
	if !res {
		if rules["msg"] != "" {
			msg = rules["msg"]
		} else {
			//没有自定义错误提示
			msg = "param " + field + " not email"
		}
	}
	return res, msg
}

func (v *Validate) isIp(field string, val string, rules map[string]string) (bool, string) {
	msg, res := "", false
	if ipPattern.MatchString(val) {
		res = true
	}
	if !res {
		if rules["msg"] != "" {
			msg = rules["msg"]
		} else {
			//没有自定义错误提示
			msg = "param " + field + " not ip"
		}
	}
	return res, msg
}

func (v *Validate) isPhone(field string, val string, rules map[string]string) (bool, string) {
	msg, res := "", false
	if mobilePattern.MatchString(val) {
		res = true
	}
	if !res {
		if telPattern.MatchString(val) {
			res = true
		}
	}
	if !res {
		if rules["msg"] != "" {
			msg = rules["msg"]
		} else {
			//没有自定义错误提示
			msg = "param " + field + " not phone"
		}
	}
	return res, msg
}

func (v *Validate) Verify(params map[string]string, ruleSet map[string][]string) (bool, string) {
	v.reset()
	v.extract(ruleSet)
	for field, rules := range v.ruleSet {
		param, _ := params[field]
		//依次验证规则
		res, msg := v.factory(field, param, rules)
		if !res && msg != "" {
			return res, msg
		}
	}
	return true, ""
}

func (v *Validate) factory(field string, val string, rules map[string]map[string]string) (bool, string) {
	for ruleName, exec := range v.process {
		if rule, ok := rules[ruleName]; ok {
			//交给相应验证器进行验证
			res, msg := exec(field, val, rule)
			if !res && msg != "" {
				//验证失败
				return res, msg
			}
		}
	}
	return true, ""
}

func (v *Validate) PushRule(field, ruleName, ruleParam, msg string) {
	if info, ok := v.ruleSet[field]; ok {
		info[ruleName] = map[string]string{
			"ruleParam": ruleParam,
			"msg":       msg,
		}
	} else {
		v.ruleSet[field] = map[string]map[string]string{
			ruleName: {
				"ruleParam": ruleParam,
				"msg":       msg,
			},
		}
	}
}

//处理路由规则
func (v *Validate) extract(ruleSet map[string][]string) {
	for field, rule := range ruleSet {
		//规则信息
		parseRule := strings.Split(rule[0], "|")
		//错误提示
		if len(rule) > 1 {
			parseMsg := strings.Split(rule[1], "|")
			//有错误提示
			for _, r := range parseRule {
				ruleName := r
				var ruleParam string
				if v.isParamRule(r) {
					ruleName, ruleParam = v.extractRule(r)
				}
				msg := v.extractMsg(parseMsg, ruleName)
				v.PushRule(field, ruleName, ruleParam, msg)
			}
		} else {
			for _, r := range parseRule {
				ruleName := r
				var ruleParam string
				if v.isParamRule(r) {
					ruleName, ruleParam = v.extractRule(r)
				}
				v.PushRule(field, ruleName, ruleParam, "")
			}
		}
	}
}

//提取验证规则
func (v Validate) extractRule(rule string) (string, string) {
	rules := strings.Split(rule, ":")
	return rules[0], rules[1]
}

//提取错误信息
func (v Validate) extractMsg(msg []string, ruleName string) string {
	for _, v := range msg {
		if strings.Contains(v, ruleName) {
			return strings.TrimSpace(strings.Split(v, ":")[1])
		}
	}
	return ""
}

//判断是否带参数规则
func (v Validate) isParamRule(ruleName string) bool {
	if strings.Contains(ruleName, "max:") {
		return true
	}
	if strings.Contains(ruleName, "min:") {
		return true
	}
	if strings.Contains(ruleName, "maxLength:") {
		return true
	}
	return false
}
