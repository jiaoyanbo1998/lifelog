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

