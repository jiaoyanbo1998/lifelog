技术栈：Gin GORM MySQL Redis Kafka JWT Viper MinIO Zap ELK Wire等。

项目介绍：LifeLog是一个使用Go语言和Gin框架开发的web应用，用来收集用户的日常行为信息，收集到的文字、图片和视频用于我们实验室课题组的LifeLog研究，目前已经收集了大概4万条数据。

登录功能：使用JWT和Gin的Middleware进行登录校验，使用Lua脚本限制发送短信的次数和验证短信的次数，使用布隆过滤器实现手机号黑名单，使用Zset实现滑动窗口限流，使用Zset实现历史密码库，使用bcrypt对密码进行加密等，实现了一个安全性较高的登录功能。

退出功能：使用长短token和redis实现了一个安全性较高的退出功能。

短信服务：通过增强安全性，提高可用性和可观测性，构建了一个更加稳健的短信服务。

日常管理：我的日常列表，增加我的日常，删除我的日常，修改我的日常，查询我的日常，好友日常列表实现了一个丰富的缓存方案，如缓存预加载，缓存特定的日常，添加缓存监控，设计特定的缓存淘汰策略等。

互动管理：点赞，评论，收藏，阅读，互关列表，关注列表，粉丝列表。使用Sarama实现Kafka异步批量消费评论和点赞，使用邻接表 维护评论的树形结构等。

文件上传：使用MinIO作为文件系统，存储图片和视频。小文件直接上传，大文件采用分片上传。

日志采集：使用ELK实现日志的统一采集，存储和检索。使用Zap框架采集日志，使用Filebeat读取.log文件，然后使用Logstash过滤插件处理后，将日志发送到Elasticsearch统一存储，最后使用Kibana可视化检索日志。


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

