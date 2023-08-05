package engine

import (
	"github.com/xjh-creator/reptile/internal/collect"
	"go.uber.org/zap"
)

type Schedule struct {
	options

	requestCh chan *collect.Request // requestCh 负责接收请求，并将请求存储到 reqQueue 队列中
	workerCh  chan *collect.Request // workerCh 通道负责传送任务，负责分配任务给 worker
	out       chan collect.ParseResult // out 负责处理爬取后的数据，
}

func NewSchedule(opts ...Option) *Schedule {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	s := &Schedule{}
	s.options = options

	return s
}

// Run 初始化三个通道，并完成下一步的存储操作
func (s *Schedule) Run() {
	requestCh := make(chan *collect.Request)
	workerCh := make(chan *collect.Request)
	out := make(chan collect.ParseResult)
	s.requestCh = requestCh
	s.workerCh = workerCh
	s.out = out
	go s.Schedule()

	for i := 0; i < s.WorkCount; i++ {
		go s.CreateWork()
	}

	s.HandleResult()
}

// Schedule 创建调度程序，负责的是调度的核心逻辑。
func (s *Schedule) Schedule() {
	var reqQueue = s.Seeds
	go func() {
		for {
			var req *collect.Request
			var ch chan *collect.Request

			if len(reqQueue) > 0 {
				req = reqQueue[0]
				reqQueue = reqQueue[1:]
				ch = s.workerCh
			}
			select {
			case r := <-s.requestCh:
				reqQueue = append(reqQueue, r)

			case ch <- req:
			}
		}
	}()
}

func (s *Schedule) CreateWork() {
	for {
		r := <-s.workerCh
		if err := r.Check(); err != nil {
			s.Logger.Error("check failed",
				zap.Error(err),
			)
			continue
		}

		body, err := s.Fetcher.Get(r)
		if len(body) < 6000 {
			s.Logger.Error("can't fetch ",
				zap.Int("length", len(body)),
				zap.String("url", r.Url),
			)
			continue
		}
		if err != nil {
			s.Logger.Error("can't fetch ",
				zap.Error(err),
				zap.String("url", r.Url),
			)
			continue
		}

		result := r.ParseFunc(body, r)
		s.out <- result
	}
}

// HandleResult 处理爬取并解析后的数据结构
func (s *Schedule) HandleResult() {
	for {
		select {
		case result := <-s.out:
			for _, req := range result.Requests {
				s.requestCh <- req
			}
			for _, item := range result.Items {
				// todo: store
				s.Logger.Sugar().Info("get result: ", item)
			}
		}
	}
}
