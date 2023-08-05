package collect

//采集规则树
type RuleTree struct {
	Root  func() []*Request // 根节点(执行入口),用于生成爬虫的种子网站
	Trunk map[string]*Rule  // 规则哈希表,用于存储当前任务所有的规则,哈希表的 Key 为规则名，Value 为具体的规则。每一个规则就是一个 ParseFunc 解析函数。
}

// 采集规则节点
type Rule struct {
	ParseFunc func(*Context) ParseResult // 内容解析函数
}


type Context struct {
	Body []byte
	Req  *Request
}
