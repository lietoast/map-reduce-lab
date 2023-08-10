package pub_sub

import "sync"

// subscriber 订阅者
type subsriber chan interface{}

// filter 频道过滤器
type filter func(interface{}) bool

// 发布-订阅模型
type Radio struct {
	mtx      *sync.RWMutex        // 读写锁
	channels map[subsriber]filter // 订阅者 -> 频道过滤器
	csize    int                  // 订阅者处最多积压的消息个数
}

// NewRadio 新建发布-订阅模型
// params:
// @csize: 每个订阅者处最多积压的消息个数, 定义域为[1,100]
// 返回一个新的发布-订阅模型
func NewRadio(csize int) *Radio {
	if csize < 1 || csize > 100 {
		csize = 100
	}
	return &Radio{
		mtx:      new(sync.RWMutex),
		channels: map[subsriber]filter{},
		csize:    csize,
	}
}

// Subscribe 订阅频道, 返回接收该频道消息的Go通道
// params:
// @topicMatcher: 当Radio接收到新消息时, 会判断消息的主题与该频道所关注的主题是否匹配, 如匹配, 则发送到频道里
func (r *Radio) Subscribe(topicMatcher filter) chan interface{} {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	sub := make(chan interface{}, r.csize)
	r.channels[sub] = topicMatcher

	return sub
}

// Publish 发布消息
// params:
// @topic: 消息主题
// @content: 消息内容
func (r *Radio) Publish(topic interface{}, content interface{}) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	for sub, filter := range r.channels {
		if !filter(topic) {
			continue
		}
		sub <- content
	}
}

// Close 关闭所有频道
func (r *Radio) Close() {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	for sub := range r.channels {
		close(sub)
	}
}
