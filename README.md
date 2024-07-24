## 目录结构
> ├── api                         暴露接口目录
> ├── bootstrap                   初始化启动目录
> │  ├── app.go
> │  └── internal
> │     ├── config.go
> │     ├── ...
> ├── cmd                         启动目录
> │  └── main.go                  启动文件
> ├── config
> │  └── config.yaml              配置文件
> ├── docs                        接口文档
> ├── domain                      领域模型
> ├── go.mod
> ├── go.sum
> ├── internal                    内部业务代码
> ├── middleware                  中间件
> ├── repository
> │  └── dao
> └── service