package proxy

import (
	"errors"
	"net/http"
	"net/url"
	"sync/atomic"
)

type ProxyFunc func(*http.Request) (*url.URL, error)

// RoundRobinProxySwitcher 函数会接收代理服务器地址列表，将其字符串地址解析为 url.URL，
// 并放入到 roundRobinSwitcher 结构中，该结构中还包含了一个自增的序号 index。
func RoundRobinProxySwitcher(proxyURLs ...string) (ProxyFunc, error) {
	if len(proxyURLs) < 1 {
		return nil, errors.New("Proxy URL list is empty")
	}

	urls := make([]*url.URL, len(proxyURLs))
	for i, u := range proxyURLs {
		parsedUrl, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		urls[i] = parsedUrl
	}

	return (&roundRobinSwitcher{urls, 0}).GetProxy, nil
}

type roundRobinSwitcher struct {
	proxyURLs []*url.URL
	index     uint32
}

// GetProxy atomic.AddUint32 会将 index 加 1，并通过取余操作实现对代理地址的轮询。
func (r *roundRobinSwitcher) GetProxy(pr *http.Request) (*url.URL, error) {
	index := atomic.AddUint32(&r.index, 1) - 1
	u := r.proxyURLs[index%uint32(len(r.proxyURLs))]

	return u, nil
}


