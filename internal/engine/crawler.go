package engine

import (
	"github.com/xjh-creator/reptile/internal/collect"
	"go.uber.org/zap"
	"sync"
)

type Crawler struct {
	options

	out         chan collect.ParseResult
	Visited     map[string]bool
	VisitedLock sync.Mutex
}

func NewEngine(opts ...Option) *Crawler {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}

	e := &Crawler{}
	e.Visited = make(map[string]bool, 100)
	out := make(chan collect.ParseResult)
	e.out = out
	e.options = options

	return e
}

func (c *Crawler) Run() {
	go c.Schedule()
	for i := 0; i < c.WorkCount; i++ {
		go c.CreateWork()
	}

	c.HandleResult()
}

func (c *Crawler) Schedule() {
	var reqs []*collect.Request
	for _, seed := range c.Seeds {
		seed.RootReq.Task = seed
		seed.RootReq.Url = seed.Url
		reqs = append(reqs, seed.RootReq)
	}

	go c.scheduler.Schedule()
	go c.scheduler.Push(reqs...)
}

func (c *Crawler) CreateWork() {
	for {
		r := c.scheduler.Pull()
		if err := r.Check(); err != nil {
			c.Logger.Error("check failed",
				zap.Error(err),
			)
			continue
		}

		// 判断当前请求是否已被访问
		if c.HasVisited(r) {
			c.Logger.Debug("request has visited",
				zap.String("url:", r.Url),
			)
			continue
		}
		// 设置当前请求已被访问
		c.StoreVisited(r)

		body, err := r.Task.Fetcher.Get(r)
		if len(body) < 6000 {
			c.Logger.Error("can't fetch ",
				zap.Int("length", len(body)),
				zap.String("url", r.Url),
			)
			continue
		}
		if err != nil {
			c.Logger.Error("can't fetch ",
				zap.Error(err),
				zap.String("url", r.Url),
			)
			continue
		}
		result := r.ParseFunc(body, r)

		if len(result.Requests) > 0 {
			go c.scheduler.Push(result.Requests...)
		}

		c.out <- result
	}
}

func (c *Crawler) HandleResult() {
	for {
		select {
		case result := <-c.out:
			for _, item := range result.Items {
				// todo: store
				c.Logger.Sugar().Info("get result: ", item)
			}
		}
	}
}

func (c *Crawler) HasVisited(r *collect.Request) bool {
	c.VisitedLock.Lock()
	defer c.VisitedLock.Unlock()
	unique := r.Unique()

	return c.Visited[unique]
}

func (c *Crawler) StoreVisited(reqs ...*collect.Request) {
	c.VisitedLock.Lock()
	defer c.VisitedLock.Unlock()

	for _, r := range reqs {
		unique := r.Unique()
		c.Visited[unique] = true
	}
}


