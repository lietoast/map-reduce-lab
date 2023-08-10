package input

import (
	"io"
	"os"
)

// LocalFileReader 本地文件读取器
// 从本地文件中读取数据, 以供后续处理
type LocalFileReader struct {
	Path      string   // 文件路径
	Offset    int64    // 文件偏移量
	BlockSize int64    // 需要读取的数据量
	Result    []string // 读取到的最终数据
}

func (lfr *LocalFileReader) Read() int64 {
	// 打开文件
	f, err := os.OpenFile(lfr.Path, os.O_RDONLY, 0600)
	if err != nil {
		return -1
	}
	defer f.Close()

	// 文件指针偏移
	f.Seek(lfr.Offset, os.SEEK_SET)

	// 读入原始数据
	chunkSize := 1024               // 读缓冲大小为1KB
	buff := make([]byte, chunkSize) // 创建读缓冲

	count := int64(0)
	data := []byte{}
	for count < lfr.BlockSize {
		n, err := f.Read(buff)
		// 读取时发生错误
		if err != nil && err != io.EOF {
			return -1
		}
		// 读取完毕
		if err == io.EOF {
			break
		}
		// 去尾
		if count+int64(n) > lfr.BlockSize {
			n = int(lfr.BlockSize - count)
		}
		// 追加数据
		data = append(data, buff[:n]...)
	}
	lfr.Result = append(lfr.Result, string(data))

	// 返回实际读取到的数据量
	return count
}

func (lfr *LocalFileReader) Next() string {
	if lfr.Result == nil || len(lfr.Result) <= 0 {
		return ""
	}

	result := lfr.Result[0]
	lfr.Result = lfr.Result[1:]
	return result
}
