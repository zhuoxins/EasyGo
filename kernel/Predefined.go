package kernel

const (
	ProjectName = "EasyGo"
)

type Service interface {
	Do(handle func())
}

//自定义服务
type ConsumeExecutor interface {
	Exec(handle func())
}
