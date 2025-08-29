# UnitedRhino Core - SaaS Middle Platform Core Module

[![Go](https://github.com/zeromicro/go-zero/workflows/Go/badge.svg?branch=master)](https://github.com/unitedrhino/core/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/unitedrhino/core)](https://goreportcard.com/report/github.com/unitedrhino/core)
[![Go Reference](https://pkg.go.dev/badge/github.com/unitedrhino/core.svg)](https://pkg.go.dev/github.com/unitedrhino/core)
![GitHub Repo stars](https://img.shields.io/github/stars/unitedrhino/core)

> 📖 [English](README.en.md) | [中文](README.md)

## 🚀 Product Introduction

**UnitedRhino Core** is the SaaS middle platform core module of the UnitedRhino IoT platform, providing comprehensive multi-tenant, multi-project, and multi-application management capabilities. This is the most core module of the UnitedRhino platform, where all basic functions unrelated to specific business are implemented, such as tenant management, project management, user management, role permissions, application management, etc.

With this module, users can directly develop their own business systems, not limited to IoT applications, and quickly build various SaaS platforms.

> 📖 [Complete Documentation](https://doc.unitedrhino.com/) | 🌐 [Online Demo](https://doc.unitedrhino.com/use/ezkveztg/)

---

## ✨ Core Features

### 🏢 Multi-tenant & Multi-project Architecture
- Supports multi-tenant and multi-project capabilities, enabling low-cost custom project development
- Complete tenant function authorization, application authorization, module authorization, and menu authorization
- Flexible regional management with tree-structured organization support

### 🔧 Flexible Deployment Modes
- One codebase simultaneously supports k8s and docker deployment
- Supports monolithic, microservice, and cluster modes
- Quick integration through HTTP, gRPC, NATS, and WebSocket

### ⚡ High-Performance Design
- Written in Golang with high-performance components
- Integrates high-performance components like EMQX, NATS, TDengine
- Extreme performance optimization ensures stability under extreme conditions

### 🛠️ Rapid Development Capabilities
- Provides Web, mini-program, App, device SDK, and module support
- Quick deployment with minimal development, rapidly implementing business requirements
- Complete permission management and data permission control

## 🏗️ System Architecture

UnitedRhino Core serves as the core module of the SaaS middle platform, providing complete basic service support for upper-layer business.

### SaaS Middle Platform Architecture Design

![SaaS Platform Architecture](./doc/assets/SaaS平台.png)

---

## 🔧 Core Services

### syssvr - System Management Service

syssvr is the core module of the SaaS system, the most basic module that doesn't depend on any other modules.

| Function Module | Function Description |
|---------|---------|
| **User Management** | Provides user login (WeChat, DingTalk, phone, email, etc.), logout, session validation, and maintenance |
| **Tenant Management** | Provides tenant management, OEM, tenant function authorization, tenant application authorization, tenant module authorization, and tenant menu authorization |
| **Project Management** | Provides project management and project configuration management |
| **Regional Management** | Regions are tree-structured under projects, e.g., East China Project → XX Street → XX Building → XX Room |
| **Data Permissions** | Provides regional and project data permission management, supporting project groups (family groups) and room authorization with fine-grained permission control |
| **Notification Management** | Notification channels, templates, and configuration management, supporting SMS, email, DingTalk, WeChat push, WeCom in-site messages, and phone notifications |
| **Authorization Management** | Function permission management, integrating multiple interfaces into one authorization through goctl tool-generated configuration files |
| **Role Management** | Supports multiple roles, high performance, can authorize projects and regions, as well as function permissions, application modules, and menu permissions |
| **Application Management** | Divided into Web applications, App applications, and mini-program applications, with each Web application composed of multiple modules |
| **Module Management** | Modules are currently used for Web, each module being an independent system, such as IoT, platform management, system management, marketing management, etc. |
| **Dictionary Management** | Provides enhanced dictionary support, not only supporting list format but also tree structure |
| **Slot Management** | Uses slot system for real-time notifications and validations in system expansion areas, enhancing system extensibility |
| **Log Management** | Provides operation logs and login logs |
| **Operations Management** | Provides work order management and feedback functionality |

### datasvr - Data Management Service

In different systems, data analysis is a very important part, and this part is time-consuming and labor-intensive for both frontend and backend. Frontend has low-code solutions for dashboard issues, while backend also has data management services that can be configured to dynamically obtain data required by frontend, convenient and fast.

![Data Analysis Example](./doc/assets/数据分析示例.png)

### timed - Scheduled Task Service

Supports scheduled tasks and delayed tasks, implemented using [asynq](https://github.com/hibiken/asynq) at the bottom layer.

#### Trigger Methods
- **Scheduled Trigger**: Execute tasks according to time plans
- **Delayed Trigger**: Execute tasks after specified delay time
- **Message Queue Trigger**: Trigger tasks through message queues

#### Execution Methods
- **Message Queue Sending**: Send messages to specified queues
- **SQL Execution**: Execute database operations
- **Script Execution**: Execute custom scripts

You can also conveniently view task execution records and results on the management platform, supporting immediate execution and task information modification.

![Task Management](./doc/assets/任务管理.png)

---

## 🛠️ Technology Stack

### Backend Technology
- **Microservice Framework**: [go-zero](https://go-zero.dev/) - High-performance microservice framework
- **High-Performance Cache**: [Redis](https://redis.io/) - In-memory data structure store
- **Message Queue**: [NATS](https://docs.nats.io/) - High-performance messaging system
- **Relational Database**: [MySQL/MariaDB](https://mariadb.com/) or PostgreSQL
- **Service Registry**: [etcd](https://etcd.io/) (microservice mode)
- **Object Storage**: [MinIO](https://min.io/) - Cloud-native lightweight object storage

### Frontend Technology
- **Framework**: [Vue.js](https://vuejs.org/) - Progressive JavaScript framework
- **UI Components**: [Ant Design Vue](https://antdv.com/) - Enterprise-grade design components

### Mobile
- **Mini-Program**: [uni-app Vue3](https://uniapp.dcloud.net.cn/) - Cross-platform development framework
- **App**: [uni-app X](https://doc.dcloud.net.cn/uni-app-x/) - Supports Android, iOS, HarmonyOS

---

## 🚀 Quick Start

### 📋 Requirements
- **Go**: 1.19+
- **Database**: MySQL 5.7+ or PostgreSQL
- **Cache**: Redis 6.0+
- **Container**: Docker (optional, recommended)

### 🛠️ Quick Deployment

#### 📚 Detailed Deployment Guide
From environment preparation to service startup, step-by-step deployment guide

[📖 View Deployment Documentation](https://doc.unitedrhino.com/use/046431/)

### 💡 Having Issues?

- **📖 View Documentation**: [Complete Documentation](https://doc.unitedrhino.com/)
- **🐛 Submit Issue**: [GitHub Issues](https://github.com/unitedrhino/core/issues)
- **💬 Join Community**: Scan QR code to join WeChat group for technical support

---

## 💬 Contact Us

### 📱 WeChat Community

> 💬 **Group already has 500+ IoT developers, looking forward to your joining!**

![WeCom QR Code](./doc/assets/企业微信二维码.png)

**Scan to join and start your IoT journey!**

### 📢 Official Account

Follow our official account for more exciting content:

![Official Account](./doc/assets/公众号.jpg)

### 📞 Other Contact Methods

- **WeChat**: godLei6
- **Website**: [https://doc.unitedrhino.com/](https://doc.unitedrhino.com/)
- **GitHub Issues**: [Submit Feedback](https://github.com/unitedrhino/things/issues)

## 🤝 Open Source Community

- **GitHub**: [UnitedRhino GitHub](https://github.com/unitedrhino/things)
- **Gitee**: [UnitedRhino Gitee](https://gitee.com/unitedrhino/things)
- **Website**: [UnitedRhino Website](https://doc.unitedrhino.com/)

---

## 👥 Contributors

Thanks to everyone who has contributed!

[![Contributors](https://contributors-img.web.app/image?repo=unitedrhino/core)](https://github.com/unitedrhino/core/graphs/contributors)

---

## ⭐ Star History

![Star History](https://starchart.cc/unitedrhino/core.svg)

> 💡 **Note**: For latest version updates, please visit: [Gitee](https://gitee.com/unitedrhino/core)

---

## 📄 License

This project is licensed under the [Apache License 2.0](LICENSE).

---

## 🚀 Start Your SaaS Middle Platform Journey

If this project helps you, please give us a ⭐ Star

💬 Join our community and learn together with 500+ developers

[⭐ Star on GitHub](https://github.com/unitedrhino/core) | [⭐ Star on Gitee](https://gitee.com/unitedrhino/core)

*Made with ❤️ by UnitedRhino Team*
