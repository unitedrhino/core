# 产品概述
[![Go](https://github.com/zeromicro/go-zero/workflows/Go/badge.svg?branch=master)](https://github.com/unitedrhino/things/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/zeromicro/go-zero)](https://goreportcard.com/report/github.com/unitedrhino/things)
[![Go Reference](https://pkg.go.dev/badge/github.com/unitedrhino/things.svg)](https://pkg.go.dev/github.com/unitedrhino/things)

**联犀** 是一款基于 Go 语言开发的商业级 SaaS 云原生微服务物联网平台，致力于帮助企业快速构建自己的物联网应用，实现快速业务落地。  
这个仓库是联犀中中台模块,是联犀中最为核心的模块,所有和具体业务无关的功能都会在这里去实现,
如租户,项目,用户,角色,应用等等,用户借助该模块也可以直接开发自己的业务,不一定是物联网的.

[文档](https://doc.unitedrhino.com/)
## 技术优势
- **快速商用**：提供web,小程序,app,设备sdk及模组,少量开发即可上线。
- **超强的拓展能力**：一套代码同时支持k8s,docker,可以通过http,grpc,nats及ws快速集成,完善的租户管理,应用管理,少量代码即可快速实现自己的业务
- **高性能**：使用golang编写,选用高性能组件(emqx,nats,tdengine),极致的性能优化保证极端情况的稳定

<img  src="./doc/assets/SaaS平台.png">  


## 开源社区
- **GitHub**: [联犀 GitHub](https://github.com/unitedrhino/core)
- **Gitee**: [联犀 Gitee](https://gitee.com/unitedrhino/core)










## 服务介绍
### syssvr(系统管理服务)
syssvr是saas系统的核心模块,是最基础的模块,不依赖其他任何模块

| 功能   | 说明                                                                             |
|------|:-------------------------------------------------------------------------------|
| 用户管理 | 提供用户登录(微信,钉钉,手机,邮箱等)登出,会话校验保持等                                                 |
| 租户管理 | 提供租户的管理及OEM,租户功能授权,租户应用授权,租户模块授权,租户菜单授权                                        |
| 项目管理 | 提供项目的管理及项目配置的管理                                                                |
| 区域管理 | 区域是项目下的以树结构存在的,如 华东项目 xx街道 xx房 xx室,街道,房间和室都是一个区域                               |
| 数据权限 | 提供区域和项目的数据权限管理,如项目组(家庭组)及房间的授权,可以细到管理权限,读写权限和只读权限                              |
| 通知管理 | 通知的通道,模版及配置管理,支持短信,邮箱,钉钉(机器人及webhook),微信推送,企业微信站内信及电话通知...                     |
| 授权管理 | 功能权限管理,通过继承goctl工具可以直接生成配置文件导入到系统中,将多个接口集成为一个授权,并提供高效的管理                       |
| 角色管理 | 支持多角色,高性能,可以授权项目和区域,也可以授权功能权限,应用模块及菜单权限                                        |
| 应用管理 | 分为web应用,app应用,小程序应用,而一个web应用由多个模块组成,如系统管理模块,物联网模块,营销模块等                        |
| 模块管理 | 模块目前只用于web, 每个模块都是独立的系统,如物联网,平台管理,系统管理,营销管理等                                   |
| 字典管理 | 提供加强版的字典,不仅支持列表形式,还支持树结构                                                       |
| 插槽管理 | 在系统需要拓展的地方可以使用插槽系统进行实时通知,实时校验,如区域创建的时候可以通过插槽系统来让物联网系统判断该区域是否能创建子区域,以便让系统的拓展性更强 |
| 日志管理 | 提供操作日志,登录日志                                                                    |
| 运营管理 | 提供工单管理及反馈的功能                                                                   |

### datasvr(数据管理服务)
在不同的系统中,数据分析是非常重要的一个部分,而这个部分对前后端来说都是耗时耗力的行为,前端有低代码解决大屏的问题,后端也有数据管理服务可以配置后即可动态获取前端所需要的数据,方便快捷
<img src="./doc/assets/\数据分析示例.png"/></img>

### timed(timedjobsvr,timedschedulersvr) 定时任务生产者和消费者模块
支持定时任务及延时任务,底层使用 [asynq](https://github.com/hibiken/asynq) 进行实现
支持以下触发方式:
* 定时触发
* 延时触发
* 消息队列触发

支持以下任务执行方式:
* 消息队列发送
* sql执行
* 脚本执行

同时可以很方便的在管理平台上看到任务的执行记录和执行结果,同时也支持立马执行及修改任务信息  
<img src="./doc/assets/\任务管理.png"/></img>








## 技术栈

### 后端
1. 微服务框架：[go-zero](https://go-zero.dev/)
2. 高性能缓存：[redis](https://redis.io/)
3. 高性能消息队列：[nats](https://docs.nats.io/)
4. 关系型数据库：[mysql (推荐使用 MariaDB 或 MySQL 5.7)](https://mariadb.com/) 或 pgsql，未来将支持更多数据库
5. 微服务注册中心（单体可不使用）：etcd
6. 云原生轻量级对象存储：[minio](https://min.io/)

### 前端
1. 渐进式 JavaScript 框架：[vue](https://cn.vuejs.org/)
2. 企业级设计组件：[ant design](https://antdv.com/docs/vue/introduce-cn/)

### 小程序
1. [uniapp vue3](https://uniapp.dcloud.net.cn/)

### app(安卓, iOS, 鸿蒙)
1. [uniapp x](https://doc.dcloud.net.cn/uni-app-x/)
## 贡献者

感谢所有已经做出贡献的人!

### 后端

<a href="https://github.com/unitedrhino/things/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=unitedrhino/things" />
</a>


## 社区

- 微信(加我拉微信群): `godLei6` (需备注“来自github”)
- [官网](https://doc.unitedrhino.com/)
- 微信二维码
- <img style="width: 300px;" src="./doc/assets/微信二维码.jpg">

## 收藏

<img src="https://starchart.cc/unitedrhino/things.svg">
