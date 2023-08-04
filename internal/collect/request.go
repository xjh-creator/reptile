package collect

import "time"

// Request ParseFunc 函数会解析从网站获取到的网站信息，并返回 Requests 数组用于进一步获取数据。而 Items 表示获取到的数据。
type Request struct {
	Url       string
	Cookie    string
	WaitTime  time.Duration
	ParseFunc func([]byte, *Request) ParseResult
}

// ParseResult 爬取后获取的数据
type ParseResult struct {
	Requests []*Request
	Items     []interface{}
}
