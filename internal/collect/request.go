package collect

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

// Task 一个任务实例，
type Task struct {
	Url         string
	Cookie      string
	WaitTime    time.Duration
	MaxDepth    int
	Visited     map[string]bool
	VisitedLock sync.Mutex
	RootReq     *Request
	Fetcher     Fetcher
}

// Request 单个请求, 结构体会在每一次请求时发生变化，但是我们希望有一个字段能够表示一整个网站的爬取任务，
// 因此我们需要抽离出一个新的结构 Task 作为一个爬虫任务，而 Request 则作为单独的请求存在。
// ParseFunc 函数会解析从网站获取到的网站信息，并返回 Requests
// 数组用于进一步获取数据。而 Items 表示获取到的数据。
type Request struct {
	Task      *Task

	unique    string
	Method    string
	Url       string
	Depth     int
	Priority  int
	ParseFunc func([]byte, *Request) ParseResult
}

// ParseResult 爬取后获取的数据
type ParseResult struct {
	Requests []*Request
	Items     []interface{}
}

func (r *Request) Check() error {
	if r.Depth > r.Task.MaxDepth {
		return errors.New("Max depth limit reached")
	}

	return nil
}

// Unique 请求的唯一识别码
func (r *Request) Unique() string {
	block := md5.Sum([]byte(r.Url + r.Method))

	return hex.EncodeToString(block[:])
}
