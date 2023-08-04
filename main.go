package main

import (
	"fmt"
	"github.com/xjh-creator/reptile/internal/collect"

)

func main()  {
	url := "https://www.thepaper.cn/"
	var f collect.Fetcher = collect.BrowserFetch{}
	body, err := f.Get(url)

	if err != nil {
		fmt.Println("read content failed:%v", err)
		return
	}

	fmt.Println(string(body))
}