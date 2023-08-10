package map_workder

import (
	"fmt"

	"github.com/lietoast/map-reduce-lab/input"
)

// MapWorker
type MapWorker interface {
	// Map函数定义
	Map(in input.InputReader) error
}

// RunMap 执行Map任务
// params:
// @w: 执行Map任务的Worker
// @in: 读取器组件
// @readMod: 读取模式
func RunMap(w MapWorker, in input.InputReader, readMod string) error {
	// 根据读取模式, 选择读取模式
	switch readMod {
	default: // 默认的模式为直接读取一整个数据块
		if in.Read() < 0 {
			return fmt.Errorf("err occur while reading data source")
		}
	}

	// 执行Map任务
	return w.Map(in)
}
