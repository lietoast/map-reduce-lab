package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/lietoast/map-reduce-lab/input"
	map_workder "github.com/lietoast/map-reduce-lab/map-workder"
	pub_sub "github.com/lietoast/map-reduce-lab/pub-sub"
)

const basePath = "/tmp/map-reduce/intermediate"

// 执行Map任务
//
//	0           1             2           3   4     5          6           7        8         9
//
// 用法: map-worker <master-addr> <worker id> <R> <map> <reader> <inputfile> <offset> <block size> <read mod>
func main() {
	offset, _ := strconv.ParseInt(os.Args[7], 10, 64)
	blockSize, _ := strconv.ParseInt(os.Args[8], 10, 64)
	workerID, _ := strconv.ParseInt(os.Args[2], 10, 64)
	r, _ := strconv.ParseInt(os.Args[3], 10, 64)

	// 创建reader
	var reader input.InputReader
	switch os.Args[5] {
	case "LocalFileReader":
		reader = &input.LocalFileReader{
			Path:      os.Args[6],
			Offset:    offset,
			BlockSize: blockSize,
			Result:    nil,
		}
	default:
		return
	}

	// 创建MapWorker对象
	var worker map_workder.MapWorker
	switch os.Args[4] {
	case "WordFreqCounter":
		worker = &map_workder.WordFreqCounter{}
	}

	// 创建Map任务的输出
	output := pub_sub.GetMapWriter()
	defer output.Close()
	var wg sync.WaitGroup
	for i := 0; i < int(r); i++ {
		// 创建收取中间产出的通道
		rid := i
		sub := output.Subscribe(func(topic interface{}) bool {
			if t, ok := topic.(string); ok {
				// 计算哈希值
				hashValue := sha256.New().Sum([]byte(t))
				// 哈希值取模
				return (int(hashValue[len(hashValue)-1]) % int(r)) == rid
			}
			return false
		})

		wg.Add(1)
		// 将最终产出写入文件
		go func() {
			defer wg.Done()
			// 路径
			absPath := fmt.Sprintf("%s/%d-%d", basePath, workerID, rid)
			// 创建文件
			f, _ := os.OpenFile(absPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
			// 写入操作
			for intermediate := range sub {
				if kv, ok := intermediate.(pub_sub.KV); ok {
					fmt.Fprintf(f, "%s %s\n", kv.Key, kv.Value)
				}
			}
		}()
	}

	// 执行Map
	err := map_workder.RunMap(worker, reader, os.Args[9])
	switch err {
	case fmt.Errorf("faulty input"): // 输入无法处理错误应当当场处理, 这里不需要额外处理
	case nil:
		wg.Wait()
		map_workder.ReportComplete(int(workerID)) // 向master报告任务完成
	default:
		map_workder.ReportError(os.Args[1], err.Error(), 0, int(workerID))
	}
}
