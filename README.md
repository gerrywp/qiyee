# qiyee

# [![gin](https://img.shields.io/badge/Gin-v1.9.1-blue)](https://github.com/gin-gonic/gin) [![gorm](https://img.shields.io/badge/Gorm-v1.25.2-blue)](https://gorm.io/zh_CN/docs/) [![sqlite](https://img.shields.io/badge/Sqlite-v1.5.2-blue)](https://www.sqlite.org/index.html)

## 介绍

互联网上的开源企业建站都过于复杂了，基本都是使用的Mysql等需要独立安装的数据库。企业门户网站本身就只用做信息展示用，用sqlite完全可以支持。体积最小的企业门户网站，简单；易用；单文件部署。

技术框架：Gin+Gorm+Sqlite  
页面展示：AdminLTE（3.2.0）【发现了更干净好用的框架tabler(https://tabler.io/admin-template),考虑后续替换。同时摒弃掉iframe,采用alpinejs做tab页展示】

## GHAT(gin+htmx+alpinejs+tabler)

1. Gin：Go语言的Web框架，性能高，易用性强。
2. HTMX：一个前端库，允许你在HTML中使用属性来发起HTTP请求并更新页面内容，而无需编写大量的JavaScript代码。
3. Alpine.js：一个轻量级的JavaScript框架，提供了类似于Vue.js的响应式数据绑定和组件化开发，但体积更小，适合用于简单的交互。
4. Tabler：一个现代化的开源管理面板模板，提供了丰富的UI组件和布局，适合用于构建管理后台和仪表盘。

登录页面和管理后台使用Alpine.js+Tabler实现，参考docs下的模板文件，HTMX只用在主页面加载一次即可，tab子页面可以直接使用。

## 前置

~~因项目使用了SQlite数据库，需要有GCC开发环境才能启动。
可以手动下下载wingw也可以是用choco安装：~~

```bash
#自行百度，省事，但是需要翻墙
choco install wingw -y
```

~~或者通过sourceforge手动下载
https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-win32/seh/~~

**移除前置条件，已替换原CGO实现的sqlite驱动，为纯GO实现的sqlite驱动。**

## 使用说明

1. clone本项目
2. 本地运行

```cmd
cd qiyee
go run cmd\main.go
```

3. 打开项目根目录下的pai.db3的sqlite数据库文件,在users数据库表新增你的用户名密码
4. 管理后端地址：<http://localhost:8080/pai/login>
   默认用户名密码为：pai/7654321
5. 企业主页：<http://localhost:8080>

## 截图

![主页](./docs/images/Snipaste_2026-05-08_20-29-18.png '主页')
![关于](./docs/images/Snipaste_2026-05-08_20-29-53.png '关于')
![产品](./docs/images/Snipaste_2026-05-08_20-30-08.png '产品')
![后端](./docs/images/Snipaste_2026-05-08_20-30-19.png '后端')
![站点信息](./docs/images/Snipaste_2026-05-08_20-30-42.png '站点信息')

## 谢谢

<a target="_blank" rel="noopener noreferrer" href="https://github.com/gerrywp/qiyee/blob/main/wx_20240314175148.jpg">
<img src="./wx_20240314175148.jpg" alt="微信收款码" title="微信支付" height="160" width="160"/>
</a>
<a target="_blank" rel="noopener noreferrer" href="https://github.com/gerrywp/qiyee/blob/main/zfb_20240314175158.jpg" >
<img src="./zfb_20240314175158.jpg" alt="支付宝收款码" title="支付宝支付" height="160" width="160"/>
</a>

东风吹尽西风起
