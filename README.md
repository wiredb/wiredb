
<div align="center">
    <img src="cmd/wiredb-org.png" style="width: 86px; height: auto; display: inline-block;">
</div>

<p align="center">WireDB is a NoSQL that supports multiple data types based on Log-structured file system.</p>


---


[![Go Report Card](https://goreportcard.com/badge/github.com/auula/wiredb)](https://goreportcard.com/report/github.com/auula/wiredb)
[![Go Reference](https://pkg.go.dev/badge/github.com/auula/wiredb.svg)](https://pkg.go.dev/github.com/auula/wiredb)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/55bc449808ca4d0c80c0122f170d7313)](https://app.codacy.com/gh/auula/wiredb/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![codecov](https://codecov.io/gh/wiredb/wiredb/graph/badge.svg?token=ekQ3KzyXtm)](https://codecov.io/gh/wiredb/wiredb)
[![DeepSource](https://app.deepsource.com/gh/wiredb/wiredb.svg/?label=active+issues&show_trend=true&token=sJBjq88ZxurlEgiOu_ukQ3O_)](https://app.deepsource.com/gh/wiredb/wiredb/?ref=repository-badge)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![release](https://img.shields.io/github/release/wiredb/wiredb.svg)](https://github.com/wiredb/wiredb/releases)



---

[简体中文](#) | [English](#)

---

## 🎉 Feature

- 支持多种内置的数据结构
- 高吞吐量、低延迟、高效批量数据写入
- 支持磁盘数据存储和磁盘垃圾数据回收
- 支持磁盘数据静态加密和静态数据压缩
- 支持 IP 白名单功能保障数据的安全访问
- 支持通过基于 RESTful API 协议操作数据

---

## 🚀 Quick Start

使用 Docker 可以快速部署 [`wiredb:latest`](https://hub.docker.com/r/auula/wiredb) 的镜像来测试 WireDB 提供的服务，运行以下命令即可拉取 WireDB 镜像：

```bash
docker pull auula/wiredb:v1.0.0
```

运行 WireDB 镜像启动容器服务，并且映射端口到外部主机网络，执行下面的命令：

```bash
docker run -p 2668:2668 auula/wiredb:v1.0.0
```

WireDB 提供使用 RESTful API 的方式进行数据交互，理论上任意具备 HTTP 协议的客户端都支持访问和操作 WireDB 服务实例。在调用 RESTful API 时需要在请求头中添加 `Auth-Token` 进行鉴权，该密钥由 WireDB 进程自动生成，可通过容器运行时日志获取，使用以下命令查看启动日志：

```bash
root@2c2m:~# docker logs 46ae91bc73a6
                         _            ____
                 _    __(_)______ ___/ / /
                | |/|/ / / __/ -_) _  / _ \
                |__,__/_/_/  \__/\_,_/_.__/  v1.0.0

  WireDB is a NoSQL database based on Log-structured file system.
  Software License: Apache 2.0  Website: https://wiredb.github.io

[WIREDB:C] 2025/02/27 10:07:01 [WARN]	The default password is: T9EHAvi5dcIpPK9G#ADlVj4NB 👈
[WIREDB:C] 2025/02/27 10:07:01 [INFO]	Logging output initialized successfully
[WIREDB:C] 2025/02/27 10:07:01 [INFO]	Loading and parsing region data files...
[WIREDB:C] 2025/02/27 10:07:01 [INFO]	Region compression activated successfully
[WIREDB:C] 2025/02/27 10:07:01 [INFO]	File system setup completed successfully
[WIREDB:C] 2025/02/27 10:07:01 [INFO]	HTTP server started at http://172.0.0.1:2668 🚀
```

---

## 🕹️ RESTful API 

目前 WireDB 服务对外提供数据交互接口是基于 HTTP 协议的 RESTful API ，只需要通过支持  `HTTP` 协议客户端软件就可以进行数据操作。这里推荐使用 `curl` 软件进行数据交互操作，WireDB 内部提供了多种数据结构抽象，例如 Table 、List 、ZSet 、Set 、Number 、Text 类型，这些数据类型对应着常见的业务代码所需使用的数据结构，这里以 Table 类型结构为例进行 RESTful API 数据交互的演示。


Table 结构类似于 JSON 及任何有映射关系的半结构化数据，例如编程语言中的 struct 和 class 字段都可以使用 Table 进行存储，下面是一个 Table 结构 JSON 抽象：

```json
{
    "table": {
        "is_valid": false,
        "items": [
            {
                "id": 1,
                "name": "Item 1"
            },
            {
                "id": 2,
                "name": "Item 2"
            }
        ],
        "meta": {
            "version": "2.0",
            "author": "Leon Ding"
        }
    },
    "ttl": 120
}
```

下面是 curl 进行数据存储操作的例子，由于是 RESTful API 设计风格，需要在 HTTP 的请求路径 URL 加上数据类型信息。注意存储使用 HTTP 协议的 PUT 方法进行操作，使用 PUT 方法会直接创建新数据版本覆盖掉旧的数据版本，命令如下：

```bash
curl -X PUT http://localhost:2668/table/key-01 -v \
     -H "Content-Type: application/json" \
     -H "Auth-Token: T9EHAvi5dcIpPK9G#ADlVj4NB" \
     --data @tests/table.json
```

获取数据的方式只需要将 HTTP 的请求改为 GET 方式就获取 Key 响应的存储记录，命令如下：

```bash
curl -X GET http://localhost:2668/table/key-01 -v \
-H "Auth-Token: T9EHAvi5dcIpPK9G#ADlVj4NB" 
```

删除对于数据记录，只需要将 HTTP 的请求改为 DELETE 的方式即可，命令如下：

```bash
curl -X DELETE http://localhost:2668/table/key-01 -v \
-H "Auth-Token: T9EHAvi5dcIpPK9G#ADlVj4NB" 
```

更为复杂的查询和复杂更新操作，将在后续的版本更新中添加支持。


---

## 🧪 Benchmark Test

由于底层存储引擎是以 Append-Only Log 的方式将所有的操作写入到数据文件中，所以这里给出的测试用例报告，是针对的其核心文件系统 [`vfs`](./vfs/) 包的写入性能测试的结果。运行测试代码的硬件设备配置信息为（Intel i5-7360U, 8GB LPDDR3 RAM），写入基准测试结果如下：

```bash
$: go test -benchmem -run=^$ -bench ^BenchmarkVFSWrite$ github.com/auula/wiredkv/vfs
goos: darwin
goarch: amd64
pkg: github.com/auula/wiredkv/vfs
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
BenchmarkVFSWrite-4   	  130216	      9682 ns/op	    1757 B/op	      44 allocs/op
PASS
ok  	github.com/auula/wiredkv/vfs	2.544s
```

在项目根目录下有一个 [`tools.sh`](./tools.sh) 的工具脚本文件，可以快速帮助完成各项辅助工作。

---

推荐使用 Linux 发型版本来运行 WireDB 服务，WireDB 服务进程依赖配置文件中的参数，在运行 WireDB 服务之前将下面的配置内容写到 `config.yaml` 中：

```yaml
port: 2668                              # 服务 HTTP 协议端口
mode: "std"                             # 默认为 std 标准库
path: "/tmp/wiredb"                     # 数据库文件存储目录
auth: "Are we wide open to the world?"  # 访问 HTTP 协议的秘密
logpath: "/tmp/wiredb/out.log"          # WireDB 在运行时程序产生的日志存储文件
debug: false        # 是否开启 debug 模式
region:             # 数据区
    enable: true    # 是否开启数据压缩功能
    second: 1800    # 默认垃圾回收器执行周期单位为秒
    threshold: 3    # 默认个数据文件大小，单位 GB
encryptor:          # 是否开启静态数据加密功能
    enable: false
    secret: "your-static-data-secret!"
compressor:         # 是否开启静态数据压缩功能
    enable: false
allowip:            # 白名单 IP 列表，可以去掉这个字段，去掉之后白名单就不会开启
    - 192.168.31.221
    - 192.168.101.225
    - 127.0.0.1
```

---

## 🌟 Stargazers over time

[![Stargazers over time](https://starchart.cc/wiredb/wiredb.svg?background=%23ffffff&axis=%23333333&line=%23f84206)](https://starchart.cc/wiredb/wiredb)


---

## 👬 Contributors

🤝 Thanks to all the contributors below! 

![Contributors](https://contributors-img.web.app/image?repo=wiredb/wiredb)




