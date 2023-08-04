package engine

import (
	"github.com/xjh-creator/reptile/internal/collect"
	"go.uber.org/zap"
)

type ScheduleEngine struct {
	requestCh chan *collect.Request // requestCh 负责接收请求，并将请求存储到 reqQueue 队列中
	workerCh  chan *collect.Request // workerCh 通道负责传送任务，负责分配任务给 worker
	WorkCount int
	Fetcher   collect.Fetcher
	Logger    *zap.Logger
	out       chan collect.ParseResult // out 负责处理爬取后的数据，
	Seeds     []*collect.Request // 请求集
}

// Run 初始化三个通道，并完成下一步的存储操作
func (s *ScheduleEngine) Run() {
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
func (s *ScheduleEngine) Schedule() {
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

func (s *ScheduleEngine) CreateWork() {
	for {
		r := <-s.workerCh
		body, err := s.Fetcher.Get(r)
		if err != nil {
			s.Logger.Error("can't fetch ",
				zap.Error(err),
			)
			continue
		}
		result := r.ParseFunc(body, r)
		s.out <- result
	}
}

// HandleResult 处理爬取并解析后的数据结构
func (s *ScheduleEngine) HandleResult() {
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
