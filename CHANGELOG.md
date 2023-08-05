##  feat:add parse js
1. 动态规则引擎
2. 实现动态规则的引擎，可 Javascript 虚拟机，因为使用 JS 操作网页有天然的优势
3. otto 是用 Go 编写的 JavaScript 虚拟机 https://github.com/robertkrimen/otto

## feat:design task and rules
1. 一个爬虫任务会针对不同的网站有不同的处理规则，但现在的处理方式导致多个规则之间是割裂的，不便于统一管理。
2. 我们在添加初始爬虫网站 URL 时，这些种子任务是在 main 函数中注入进去的，与任务的规则之间是割裂的。但是我们需要将初始爬虫 URL 与处理规则进行统一的管理。
3. 当前的爬虫任务还是需要手动初始化才能运行，可配置化程度比较低。我们希望这些写好的静态任务在程序初始化时能够自动加载。而通过外部接口，或者只要在配置文件中指定一个任务名就能将任务调度起来。
4. 更进一步，当前的任务和规则都是静态的，静态指的是代码需要提前写好，重新编译运行才能够在运行中被调用。能否动态地增加任务和任务的规则，让程序能够动态地解析我们的规则呢？

## feat:add solving of failure
1. 进行失败处理
2. 爬取网站时，网络超时等诸多潜在风险都可能导致爬取失败
3. 可以对失败的任务进行重试。但是如果网站多次失败，那就没有必要反复重试了，可以将它们放入单独的队列中
4. 为了防止失败请求日积月久导致的内存泄露，同时也为了在程序崩溃后能够再次加载这些失败网站，最后还需要将这些失败网站持久化到数据库或文件中。

## feat:add random user-agent
1. 设置随机 User-Agent
2. 为了避免服务器检测到我们使用了同一个 User-Agent，继而判断出是同一个客户端在发出请求
3. 操作的本质就是将浏览器的不同型号与不同版本拼接起来，组成一个新的 User-Agent

## feat:add priority setting in request
1. 爬虫任务的优先级有时并不是相同的，一些任务需要优先处理

## feat:add hash to keep unique request
1. 构建一个新的结构 Crawler 作为全局的爬取实例，将之前 Schedule 中的 options 迁移到 Crawler 中，Schedule 只处理与调度有关的工作，并抽象为了 Scheduler 接口 

## refactor:add task by refactoring request
1. 之前的 Request 结构体会在每一次请求时发生变化，现在需要一个字段能够表示一整个网站的爬取任务
2. 抽离出一个新的结构 Task 作为一个爬虫任务，而 Request 则作为单独的请求存在。

## feat:add max depth limit
1. 设置爬虫的最大深度

## feat:add max depth limit
1. 设置爬虫的最大深度

## feat:add option design model
1. 用函数式选项模式改造调度引擎的初始化配置

## feat:add engine
调度引擎主要目标是完成下面几个功能：
1. 创建调度程序，接收任务并将任务存储起来；
2. 执行调度任务，通过一定的调度算法将任务调度到合适的 worker 中执行；
3. 创建指定数量的 worker，完成实际任务的处理；
4. 创建数据处理协程，对爬取到的数据进行进一步处理。

## feat:add bfs and cookie
1. 用广度优先搜索实战爬虫
2. 用 Cookie 突破反爬封锁

## feat:add zap log
1. 添加日志

## feat:add proxy
1. 添加代理

## feat:add take up an engine
之前将爬取网站信息的代码封装为了 fetch 函数，完成了第一轮的功能抽象。
但爬取的网站会越来越复杂，加上服务器本身的反爬机制等原因，我们需要用到不同的爬取技术。
例如后面会讲到的模拟浏览器访问、代理访问等。需要切换爬取方法，用模块化的方式对功能进行组合、测试，
需要对爬取网站数据的代码模块进行接口抽象。
1. 添加采取引擎

## feat:reptile a site and solve encoding
1. 简单爬取澎湃新闻页面
2. 进行字符编码处理




