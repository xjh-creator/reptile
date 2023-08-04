package main

import (
	"github.com/xjh-creator/reptile/internal/collect"
	"github.com/xjh-creator/reptile/internal/log"
	"github.com/xjh-creator/reptile/internal/proxy"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func main()  {
	plugin, c := log.NewFilePlugin("./log/log.txt", zapcore.InfoLevel)
	defer c.Close()
	logger := log.NewLogger(plugin)
	logger.Info("log init end")

	proxyURLs := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8888"}
	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		logger.Error("RoundRobinProxySwitcher failed")
	}
	//url := "https://google.com"
	url := "https://www.thepaper.cn/"
	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		Proxy: p,
	}
	body, err := f.Get(url)

	if err != nil {
		logger.Error("read content failed",
			zap.Error(err),
		)
		return
	}

	logger.Info("get content", zap.Int("len", len(body)))
}