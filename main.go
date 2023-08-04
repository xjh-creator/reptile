package main

import (
	"fmt"
	"github.com/xjh-creator/reptile/internal/collect"
	"github.com/xjh-creator/reptile/internal/proxy"
	"time"
)

func main()  {
	proxyURLs := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8889"}
	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		fmt.Println("RoundRobinProxySwitcher failed")
	}
	//url := "https://google.com"
	url := "https://www.thepaper.cn/"
	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		Proxy: p,
	}
	body, err := f.Get(url)

	if err != nil {
		fmt.Println("read content failed:%v", err)
		return
	}

	fmt.Println(string(body))
}