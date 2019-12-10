# Fabric Client

### About modules and vendor of Golang

#### 使用 golang 的 modules

* go.mod 文件中的 module 指定的是 go.mod 所有目录的包路径（相对于GOPATH, 可自定义）
* go.mod 文件中会包含项目中所有用到的包、包的版本、Golang的版本
* 如果使用`go.mod`时默认使用的包是在 $GOPATH/mod 目录，前提是指定了`GO111MODULE=on`，否则会用 vender，最后会用 $GOPATH/src

```bash
GO111MODULE=on go mod init
```

`go mod vender`会将 go.mod 中指定版本的包全部放入 vendor 目录下

```bash
GO111MODULE=on go mod vender
```