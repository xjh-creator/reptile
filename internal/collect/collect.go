package collect

import (
	"bufio"
	"fmt"
	"github.com/xjh-creator/reptile/internal/proxy"
	"go.uber.org/zap"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"net/http"
	"time"
)

/*
	采取引擎，之后所有和爬取相关的代码都会放在这个目录下。
*/

type Fetcher interface {
	Get(req *Request) ([]byte, error)
}

type BaseFetch struct {
}

func (BaseFetch) Get(r *Request) ([]byte, error) {
	resp, err := http.Get(r.Url)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error status code:%d\n", resp.StatusCode)
		return nil, err
	}
	bodyReader := bufio.NewReader(resp.Body)
	e := DetermineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())

	return io.ReadAll(utf8Reader)
}

type BrowserFetch struct {
	Timeout time.Duration
	Proxy   proxy.ProxyFunc
	Logger  *zap.Logger
}

//模拟浏览器访问
//有一些反爬机制阻止了我们对服务器的访问。如果我们使用浏览器的开发者
//工具（一般在 windows 下为 F12 快捷键），或者通过 wireshark
//等抓包工具查看数据包，会看到浏览器自动在 HTTP Header 中设置了很多内容，
//其中比较重要的一个就是 User-Agent 字段，它可以表明当前正在使用的应用程序、
//设备类型和操作系统的类型与版本。
func (b BrowserFetch) Get(r *Request) ([]byte, error) {
	client := &http.Client{
		Timeout: b.Timeout,
	}

	if b.Proxy != nil {
		transport := http.DefaultTransport.(*http.Transport)
		transport.Proxy = b.Proxy
		client.Transport = transport
	}

	req, err := http.NewRequest("GET", r.Url, nil)
	if err != nil {
		return nil, fmt.Errorf("get url failed:%v", err)
	}

	if len(r.Task.Cookie) > 0 {
		req.Header.Set("Cookie", r.Task.Cookie)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	time.Sleep(r.Task.WaitTime)
	if err != nil {
		b.Logger.Error("fetch failed",
			zap.Error(err),
		)
		return nil, err
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := DetermineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())

	return io.ReadAll(utf8Reader)
}

func DetermineEncoding(r *bufio.Reader) encoding.Encoding {

	bytes, err := r.Peek(1024)

	if err != nil {
		fmt.Printf("fetch error:%v\n", err)
		return unicode.UTF8
	}

	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}

