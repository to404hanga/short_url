# short_url

## 项目技术栈

1. Gin
2. Gorm
3. wire
4. mysql
5. redis
6. zap
7. viper
8. ...

## 项目目录结构

```shell
├─pkg 公共包
│  ├─generator 短链接生成与校验相关
│  ├─logfile 创建日志文件
│  ├─sharding 分表算法
│  └─sign 签名算法
│     └─epay EPAY签名算法
├─proto proto定义文件与生成文件
│  └─short_url short_url.proto生成内容
├─rpc 对外提供微服务
│  ├─config 配置文件
│  ├─grpc 对外提供rpc接口
│  ├─ioc 组装层
│  ├─job 定时任务
│  ├─log 日志
│  ├─repository 存储层
│  │  ├─cache 缓存层
│  │  └─dao 数据访问对象层
│  └─service 服务层
├─scripts 脚本
│  ├─etcd_data etcd-data挂载目录
│  ├─jmeter Jmeter压测计划文件
│  └─mysql Mysql初始化文件
├─test 测试文件
└─web 对外提供http服务
   ├─config 配置文件
   ├─ioc 组装层
   ├─log 日志
   ├─middlewares 中间件
   └─routes 路由层
```

## 使用前

1. 先参考 **short_url/config/config_dev.yaml** 的格式，新建一个 short_url/config/config.yaml 的配置文件。
2. 无需手动创建 short_url/log/log.txt 和 short_url/log/error_output.txt，启动时会自动创建。

## **你必须知道**

1. **不要**提交任何敏感信息，例如`api_key`、`address`或`password`。
2. 您可以使用配置文件`config.yaml`来存储某些敏感信息，但不要试图提交它。每次修改`config.yaml`的结构后，您必须同步更新`config.yaml.template`。
3. 任何时候不要用 `git push --force` 除非你知道你在干什么。
