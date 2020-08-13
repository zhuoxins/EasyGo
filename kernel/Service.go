package kernel

import "log"

func init() {
	log.Println("EasyGo Loading ...")
}

type service struct {
	serveName string
}

func (this *service) Do(handle func()) {
	log.Println("service " + this.serveName + " is ready")
	handle()
}

func Register(serveName string) Service {
	return &service{serveName: serveName}
}
