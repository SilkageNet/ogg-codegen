# OGG-Codegen

> Support OpenAPI 3.0 to generate golang gin middleware server code.

OpenAPI 可以很方便的我们定义HTTP服务的协议，以及文档展示。但在实际开发的过程中，常存在以下问题：

- 客户端和服务端在对其做实现的时候，经常会与协议有些偏差，不便于问题的定位。
- 请求和响应的结构体定义这种无脑的code，有点浪费开发实际。
- OpenAPI原本就支持Schema的定义，而在业务实现时却还需要手写Schema校验的代码。

因此，促使了 `ogg-codegen` 的出现😂。命名：`o(OpenAPI)g(Golang)g(Gin)-codegen`

## Getting Started

### Install

可直接通过 `go install` 安装到本地：

```shell
go install github.com/SilkageNet/ogg-codegen@latest
```

### Quick Started

工具开箱即用，若当前目录存在 `swagger.yaml` 文件，直接不带任何参数执行 `ogg-codegen` 即可生成一份默认的代码：

```shell
ogg-codegen
```

接着就是 **四行代码编写一个Mock服务 😂😂😂（垃圾开源框架Slogan）**

```go
package test
func main() {
    engine := gin.New()
    swagger.Register(engine, &swagger.MockServer{})
    srv := &http.Server{Addr: ":8989", Handler: engine}
    _ = srv.ListenAndServe()
}
```

### Usage

具体的命令行参数可以通过 `ogg-codegen --help` 查看。

```shell
$ ogg-codegen --help
NAME:
   ogg-codegen - A new cli application

USAGE:
   ogg-codegen [global options] command [command options] [arguments...]

VERSION:
   1.1.1

DESCRIPTION:
   OpenAPI gin golang codegen toolkit.

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --file value, -f value        Specify a path to a OpenAPI 3.0 spec file (default: "swagger.yaml")
   --cookie value                When the file path is URL and authentication is required, you can specify cookie through this parameter
   --out value, -o value         Specify a file output path (default: ".")
   --package value, -p value     The package name for generated code (default: "swagger")
   --tag value, -t value         Specify the tag list to generate
   --exSchema value, --es value  Specifies a list of schemas that do not need to be generated
   --merge, -m                   Merge files only (default: false)
   --schema                      Specifies whether the schema code needs to be generated (default: true)
   --server                      Specifies whether the server code needs to be generated (default: true)
   --import value                Specify the import list to be added additionally
   --help, -h                    show help (default: false)
   --version, -v                 print the version (default: false)
```

#### --file / -f

指定 `OpenAPI` 文档路径，支持URL和本地Path，默认值为当前目录的 `swagger.yaml`。

```shell
$ ogg-codegen -f openapi.yaml
```

#### --cookie

当 `--file` 指定路径为URL时，且访问需要鉴权时，可以通过本参数设置Cookie。

```shell
$ ogg-codegen -f URL \
              --cookie _gitlab_session=0ec19fa0212e4f28d962130fc3a758ca
```

#### --out / -o

指定生成文件的输出目录，默认为当前目录。

#### --package / -p

指定生成代码的包名 `package` 。

```shell
$ ogg-codegen -p api
```

#### --tag / -t

生成指定 `tag` 的接口实现，默认生成所有 `tag`。当后面不同类型的接口将由不同服务实现时，可以根据该参数指定自己这个服务需要实现的 `tag` ，这样就可以避免生成冗余的代码。

```shell
$ ogg-codegen --tag auth
              -t user
```

#### --exSchema / -es

指定不需要生成的 `Schema` ，默认会生成所有 `Schema`。

```shell
$ ogg-codegen --exSchema Response
              -es User
```

#### --merge / -m

指定是否为合并swagger文件，默认为 `fase`。当指定为 `true` 时，仅合并文件不会执行代码生成。

#### --schema

指定是否需要生成 `Schema` 代码，默认为 `true`。

```shell
$ ogg-codegen --schema false
```

#### --server

指定是否需要生成 `Server` 代码，默认为 `true`。

```shell
$ ogg-codegen --server false
```

#### --import

当 `yaml` 文件有其他拓展的类型时，可能会额外添加一些import，我们可以通过该参数来指定。指定格式为 `import:ignore_code`。肯定会好奇为什么会有 `ignore_code` 吧，因为并不是所有生成的代码文件都会用到该这个import，如果不忽略的话，就会报错啦：

```shell
$ ogg-codegen --import "**/ds:_ tex.JsInt64"
              --import "**/tsf:_ = tex.ToInt"
```

## Extensions

为了生成更加易用的代码，我们使用了 `OpenAPI` 原生的拓展API，以此来支持一些不错的 feature🎉。

### x-go-type

`x-go-type` 是 `Schema` 拓展属性，它的值为 `string` 类型。也就是指定一个拓展类型。因为生成文件默认导出了 `github.com/pinealctx/neptune/tex` 这个 package，那么可以直接使用这个 package 下面的的类型，比如：

- `tex.JsInt64`：对int64的拓展，客户端接收到的时候为字符串，而服务端按照int64处理。
- `tex.UnixStamp`：对时间戳的拓展，客户端接收到的时候为字符串，而服务端int64处理，对应mysql时，按照 `datetime` 类型映射。

当你要导入其他 pakcage 下的类型的时候，记得通过 `--import` 导入你要使用 package。

### x-go-enum

`x-go-enum` 是 `Schema` 拓展属性，它的值为 `map` 类型。指定该 `Schema` 为一个枚举类型，以及它的各个枚举值：

```yaml
LoginType:
  type: number
  description: Login type
  enum: 
    - 0
    - 1
  x-go-tag:
    PasswdLogin: 0
    SMSLogin: 1
```

生成的代码如下：

```go
package test

type LoginType int

const (
    PasswdLogin LoginType = 0
    SMSLogin    LoginType = 1
)
```

### x-go-tags

`x-go-tags` 是 `Schema` 拓展属性，它的值为 `map` 类型。指定该值需额外拓展的tag，比如可以为json字段添加 `gorm` 的tag。

### x-omitempty

`x-omitempty` 是 `Schema` 拓展属性，它的值为 `boolean` 类型。指定JSON的tag是否添加 `omitempty`。

### x-go-file-ext

`x-go-file-ext` 是 `Schema` 拓展属性，它的值为 `[]string` 类型。仅当为上传文件时支持，用于校验文件的拓展名。

### x-serve-file

`x-serve-file` 是 `Path` 拓展属性，它的值为 `boolean` 类型。当服务端返回文件类型的时候处理方式不同于别的请求，因此需要该拓展属性来指定。

```yaml
  /file/download:
    get:
      tags: [file]
      description: 下载文件
      operationId: DownloadFile
      x-serve-file: true
      security:
        - TokenAuth: []
      parameters:
        - $ref: "#/components/parameters/Language"
        - $ref: "#/components/parameters/DeviceType"
        - name: id
          in: query
          description: 文件ID
          required: true
          schema:
            type: string
      responses:
        200:
          description: OK
          content:
            /*:
              schema:
                type: string
                format: binary
```

## Notices

考虑到 `OpenAPI 3.0` 的完整实现复杂度过高，而我们日常业务场景也并未使用所有特性，因此在做代码生成的时候，我们需要对某些特性做相应的限制，以及增加一些我们独有的特性。

1. 需尽量将结构体独立的定义在 `components` 里。这样既方便结构体的复用，也方便做结构体的解析和校验。

2. 避免定义嵌套的结构体。一般的请求Body或者Schema的各个属性尽量平铺展示，避免多层嵌套，若遇到某些通用的结构体，应抽一个公共的Schema。这样会避免这类代码的生成：

    ```go
    package test

    type A struct {
      B struct {
        C string
      }
    }
    ```

3. 代码生成工具会通过 `tag` 来将 `path` 请求进行分类，因此我们仅允许每个 `path` 仅有一个 `tag`。

## Write at the end

工具目前还有很多特性在完善过程中，👏欢迎大家提issue和feature。


**约定大于配置，工具大于约定。**

