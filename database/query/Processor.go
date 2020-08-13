package query

func NewProcessor(driver string) (processor Processor) {
	switch driver {
	case "mysql":
		processor = &mysqlProcessor{}
	default:
		processor = nil
	}
	return
}
