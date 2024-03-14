# qiyee

<a href="https://github.com/gin-gonic/gin" target="_blank">
    <img src="https://img.shields.io/badge/Gin-v1.9.1-blue" alt="Gin">
</a>
<a href="https://gorm.io/zh_CN/docs/" target="_blank">
    <img src="https://img.shields.io/badge/Gorm-v1.25.2-blue" alt="Gorm">
</a>
<a href="https://3vshej.cn/AdminLTE/" target="_blank">
    <img src="https://img.shields.io/badge/AdminLTE-3.2.0-blue" alt="AdminLTE">
</a>
<a href="https://www.sqlite.org/index.html" target="_blank">
    <img src="https://img.shields.io/badge/Sqlite-v1.5.2-blue" alt="Sqlite">
</a>

## 介绍
互联网上的开源企业建站都过于复杂了，基本都是使用的Mysql等需要独立安装的数据库。企业门户网站本身就只用做信息展示用，用sqlite完全可以支持。因此独立开发一个最小单元的企业门户网站，达到简洁；易用；一个文件独立部署的目的。

技术框架：Gin+Gorm+Sqlite  
页面展示：AdminLTE

## 使用说明
1. clone本项目
2. 本地运行
```cmd
cd qiyee\cmd
go run main.go
```
3. 打开项目根目录下的pai.db3的sqlite数据库文件,在users数据库表新增你的用户名密码
4. 浏览器访问<http://localhost:8080/pai/login>  
默认用户名密码为：pai/7654321

## 打赏

<a href="https://github.com/gerrywp/qiyee/blob/main/wx_20240314175148.jpg" target="_blank">
<img src="./wx_20240314175148.jpg" alt="微信收款码" title="微信支付" height="160" width="160"/>
</a>
<a href="https://github.com/gerrywp/qiyee/blob/main/zfb_20240314175158.jpg" target="_blank">
<img src="./zfb_20240314175158.jpg" alt="支付宝收款码" title="支付宝支付" height="160" width="160"/>
</a>

您的支持是我前进的动力
