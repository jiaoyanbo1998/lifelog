技术栈：Gin、gRPC、GORM、MySQL、Redis、JWT、Kafka、Viper、MinIO、Zap、ELK、Wire等

项目功能：
1.登录服务
将JWT的校验逻辑封装为Gin的Middleware，实现了统一的登录校验。
使用Lua脚本限制发送短信的次数和验证短信的次数，防止非法用户恶意消耗短信资源。
使用布隆过滤器实现手机号黑名单，阻止非法用户发送短信。
使用Lua脚本实现滑动窗口限流，防止服务器被大量请求击垮。
使用Redis的Zset数据类型实现历史密码库，不允许用户使用历史密码当作新密码。
使用bcrypt不可逆哈希算法对密码进行加密，避免数据库存储明文。
使用长短token机制，短token用于登录校验，长token用于给短token续约。
2.退出服务
使用Redis实现了token的实时失效管理，有效防止token滥用和非法访问。
3.短信服务
通过增强安全性，提高可用性和可观测性，构建了一个更加稳健的短信服务。
4.生活记录服务
主要功能包括生活记录的分页查询和增删改。
设计了线上库与制作库的双库架构，支持在不影响用户观看的前提下修改生活记录。
设计了一个缓存方案，如缓存前几页生活记录，业务缓存预加载，应用启动预加载，缓存一致性解	 	决方案，redis崩溃应对方案，添加缓存监控，设计特定的缓存淘汰策略等，将系统的QPS从239提升
到1594，P99延迟从1.22秒降低到193毫秒。
5.收藏夹服务
主要功能包括收藏夹的分页查询和增删改，收藏夹详情展示。
设计了一个缓存方案，将系统的QPS从659提升到1766，P99延迟从929毫秒降低到206毫秒。
6.评论服务
使用邻接表维护评论的树形结构，确保评论的高效存储与检索。
通过Sarama实现Kafka异步批量消费评论，提升系统的并发处理能力。
7.文件服务
使用MinIO作为文件存储系统，支持小文件直接上传，大文件分片上传。
8.互动服务
实现了点赞、收藏、阅读、关注等功能，提升了用户的参与感与平台的社交属性。
9.热榜服务：？？？？？？？？？？？？？？？？？？？？？？
10.Feed流服务
实现了评论、点赞、收藏、关注、阅读的feed流服务，确保用户能够查看到最新消息。
11.日志采集
使用ELK实现了日志的统一采集，存储与可视化检索。
使用Zap框架采集日志，使用Filebeat读取.log文件，使用Logstash过滤插件处理日志，最终将日志	
发送到ElasticSearch中，并使用Kibana可视化检索日志。
12.压测
会使用k6对接口进行压测，并将压测结果在Grafana平台可视化展示。
13.Docker
会简单使用docker-compose来部署开发环境


go get -u github.com/gin-gonic/gin
go get github.com/google/wire/cmd/wire@latest
go get github.com/dlclark/regexp2
go get github.com/golang-jwt/jwt/v5
go get github.com/google/uuid
go get github.com/prometheus/client_golang/prometheus
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get github.com/redis/go-redis/v9
go get github.com/spaolacci/murmur3
go get github.com/robfig/cron/v3
go get github.com/google/wire
go get -u github.com/natefinch/lumberjack
go get github.com/spf13/viper
go get google.golang.org/grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
go get go.etcd.io/etcd/client/v3
go get -u github.com/zeromicro/go-zero
go get -u github.com/jinzhu/copier
go install github.com/link1st/go-stress-testing@latest
go get github.com/segmentio/kafka-go
go install github.com/minio/minio@latest

下载protobuf
https://github.com/protocolbuffers/protobuf/releases

安装Go和gRPC插件
go get google.golang.org/protobuf/cmd/protoc-gen-go@latest （将ProtoBuf编译为.go）
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest （将ProtoBuf编译为.grpc）

buf lint：检查proto文件是否合法

buf generate：生成.go文件，生成.grpc文件

