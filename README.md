# golang-logger

###### golang结合logrus、elasticsearch、gin网络框架封装的日志库

#### 安装
```
go get github.com/sirupsen/logrus
go get github.com/lestrrat-go/file-rotatelogs
go get github.com/rifflock/lfshook
go get github.com/olivere/elastic
go get github.com/yanjiaqi03/golang-logger
```

#### 简单方式
主要参考logrus用法即可
> https://github.com/sirupsen/logrus
```
var logger *logrus.Logger = logger.Instance()
logger.Debug(...)
logger.Info(...)
logger.Warn(...)
logger.Error(...)
logger.Fatal(...)
```

#### 高级配置
```
var logger *logrus.Logger = logger.NewBuilder()
                                  .SetLevel(logrus.InfoLevel) // 日志输出等级
                                  .SetLogPath("log/") // 本地输出路径
                                  .SetName("IndexName") // ElasticSearch索引
                                  .SetElasticHost("http://xxx:9200") // ElasticSearchHost
                                  .SetLocalHost("xxxx") // ElasticSearch中上报输出日志服务器地址，默认为本机IP，可以不用填写
                                  .Build()
```

#### 与Gin框架结合
> https://github.com/gin-gonic/gin
```
engine := gin.Default()
router := engine.Use(logger.GinLogMiddleWare()) // 已经为您封装好日志中间件
...
engine.Run(":8081")
```

由于每个请求的日志都是分散的，为了方便查询某一个请求对应的所有日志，我们规定了一个`gin_trace_id`字段来标识每一个请求。
```
// 使用logger.Context(c)在请求中进行日志上报，每一条日志就都会带有gin_trace_id字段
logger.Context(c).WithFields(logrus.Fields{
		"animal": "walrus",
	}).Info("first step")
```
