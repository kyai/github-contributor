# github-contributors

快速生成Github代码贡献者一览图，可用于README等文件。

## Usage

```go
go get github.com/kyai/github-contributors
```

#### 以`golang/go`项目为例:
```
$ github-contributors golang/go

The image is saved in github-avatar-cache/golang-go.jpeg
```

#### 生成图片如下:

![preview.jpeg](https://raw.githubusercontent.com/kyai/github-contributors/master/preview.jpeg)

## 自定义配置

```go
go get github.com/kyai/github-contributors/creator
```

```go
import "github.com/kyai/github-contributors/creator"
```

```go
creator.New(repo string).Set(config *Config).Create()
```
