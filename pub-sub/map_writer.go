package pub_sub

import "sync"

var (
	mapWriter *Radio
	mpo       sync.Once
)

type KV struct {
	Key   string
	Value string
}

func EmitIntermediate(key string, value string) {
	mapWriter.Publish(key, KV{Key: key, Value: value})
}

func GetMapWriter() *Radio {
	mpo.Do(func() {
		mapWriter = NewRadio(100) // 默认的消息序列大小为100
	})
	return mapWriter
}
