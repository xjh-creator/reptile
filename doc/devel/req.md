##正则表达式、xpath、xml 的选择
1. 正则表达式通常比较复杂而且性能低下，在实际运用过程中，我们一般采用 XPath 与 CSS 选择器进行结构化查询。
2. 比较这两种查询方法，会发现 XPath 是为 XML 文档设计的，而 CSS 选择器是为 HTML 文档专门设计的，更加简单，也更主流。

##User-Agent
###大多数浏览器使用以下格式发送 User-Agent：
```
Mozilla/5.0 (操作系统信息) 运行平台(运行平台细节) <扩展信息>
```
###使用不同的浏览器、设备，User-Agent 都会略有不同。不同应用程序的 User-Agent 参考如下：
```
Lynx: Lynx/2.8.8pre.4 libwww-FM/2.14 SSL-MM/1.4.1 GNUTLS/2.12.23

Wget: Wget/1.15 (linux-gnu)

Curl: curl/7.35.0

Samsung Galaxy Note 4: Mozilla/5.0 (Linux; Android 6.0.1; SAMSUNG SM-N910F Build/MMB29M) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/4.0 Chrome/44.0.2403.133 Mobile Safari/537.36

Apple iPhone: Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_1 like Mac OS X) AppleWebKit/603.1.30 (KHTML, like Gecko) Version/10.0 Mobile/14E304 Safari/602.1

Apple iPad: Mozilla/5.0 (iPad; CPU OS 8_4_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12H321 Safari/600.1.4

Microsoft Internet Explorer 11 / IE 11: Mozilla/5.0 (compatible, MSIE 11, Windows NT 6.3; Trident/7.0; rv:11.0) like Gecko
```

##远程访问浏览器
仅仅在请求头中传递 User-Agent 是不够的。正如我们之前提到过的，浏览器引擎会对 HTML 与 CSS 文件进行渲染，
并且执行 JavaScript 脚本，还可能会完成一些实时推送、异步调用工作。这导致内容会被延迟展示，
无法直接通过简单的 http.Get 方法获取到数据。更进一步的，
有些数据需要进行用户交互，例如我们需要点击某些按钮才能获得这些信息。这就迫切地需要我们具有模拟浏览器的能力，
或者更简单一点：直接操作浏览器，让浏览器来帮助我们爬取数据。
###要借助浏览器的能力实现自动化爬取，目前依靠的技术有以下三种：
1. 借助浏览器驱动协议（WebDriver protocol）远程与浏览器交互；
2. 借助谷歌开发者工具协议（CDP，Chrome DevTools Protocol）远程与浏览器交互；
3. 在浏览器应用程序中注入要执行的 JavaScript，典型的工具有 Cypress， TestCafe。（主要用于测试）