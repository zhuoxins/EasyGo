package kernel

import "log"

type consumerServe struct {
	serveName string
}

func (c *consumerServe) Exec(do func()) {
	log.Println("service " + c.serveName + " is ready")
	do()
}

//注册自定义服务
func BindServe(serveName string) ConsumeExecutor {
	return &consumerServe{serveName: serveName}
}
