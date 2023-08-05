package doubangroup

import (
	"github.com/xjh-creator/reptile/internal/collect"
	"regexp"
)

const urlListRe = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`

// ParseURL 对于首页样式的页面，需要获取所有帖子的 URL，这里选择使用正则表达式的方式来实现。
// 匹配到符合帖子格式的 URL 后，我们把它组装到一个新的 Request 中，用作下一步的爬取。
func ParseURL(contents []byte, req *collect.Request) collect.ParseResult {
	re := regexp.MustCompile(urlListRe)

	matches := re.FindAllSubmatch(contents, -1)
	result := collect.ParseResult{}

	for _, m := range matches {
		u := string(m[1])
		result.Requests = append(
			result.Requests, &collect.Request{
				Task:  req.Task,
				Url:   u,
				Depth: req.Depth + 1,
				ParseFunc: func(c []byte, request *collect.Request) collect.ParseResult {
					return GetContent(c, u)
				},
			})
	}

	return result
}

const ContentRe = `<div class="topic-content">[\s\S]*?阳台[\s\S]*?<div class="aside">`

// GetContent 发现正文中有对应的文字，就将当前帖子的 URL 写入到 Items 当中
func GetContent(contents []byte, url string) collect.ParseResult {
	re := regexp.MustCompile(ContentRe)

	ok := re.Match(contents)
	if !ok {
		return collect.ParseResult{
			Items: []interface{}{},
		}
	}

	result := collect.ParseResult{
		Items: []interface{}{url},
	}

	return result
}