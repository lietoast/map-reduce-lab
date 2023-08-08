package input

// InputReader 读取器组件
// 负责从本地文件/网络文件/数据库等数据源读取数据, 以供后续处理
type InputReader interface {
	// Read 读取数据
	// 返回读取到的数据量(Byte), 如发生错误, 返回-1
	Read() int64
}
