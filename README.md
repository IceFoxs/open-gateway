# open-gateway

#### 运行
```shell
make run
```

#### 构建
```shell
make build
```

### 压缩
```shell
#安装upx
brew install upx
#压缩
upx -o output_file input_file
#macos
upx --force-macos -o output_file input_file
解压
upx -d -o example_uncompressed.exe example_compressed.exe
```

### windows编译
```shell
go build  trimpath -o ./dist/opengateway.exe  ./cmd/opengateway/
```