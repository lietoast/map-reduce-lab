package map_workder

// 向master进程汇报Map任务执行的情况

type ErrID int

// 错误代码
const (
	// 输入数据有Map函数无法处理的奇点, 当errID为FaultyInput时, errMsg必须为发生错误处的偏移量
	FaultyInput = ErrID(0)
	Unknown     = ErrID(100)
)

// ReportError 汇报执行错误
// 调用一次gRPC, 向master汇报当前进程发生的错误信息、错误代码、当前进程的ID等信息
func ReportError(masterAddr, errMsg string, errID ErrID, workerID int) {
	// 待完成
}

func ReportComplete(workerID int) {
	// 待完成
}
