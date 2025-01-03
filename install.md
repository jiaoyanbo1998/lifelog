下载protobuf
    https://github.com/protocolbuffers/protobuf/releases

安装Go和gRPC插件
    go get google.golang.org/protobuf/cmd/protoc-gen-go@latest （将ProtoBuf编译为.go）
    go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest （将ProtoBuf编译为.grpc）

buf lint：检查proto文件是否合法

buf generate：生成.go文件，生成.grpc文件

