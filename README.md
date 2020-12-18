# xray
![](https://img.shields.io/github/v/release/linuxing3/trojan.svg) 
![](https://img.shields.io/docker/pulls/linuxing3/trojan.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/linuxing3/trojan)](https://goreportcard.com/report/github.com/linuxing3/trojan)
[![Downloads](https://img.shields.io/github/downloads/linuxing3/trojan/total.svg)](https://img.shields.io/github/downloads/linuxing3/trojan/total.svg)
[![License](https://img.shields.io/badge/license-GPL%20V3-blue.svg?longCache=true)](https://www.gnu.org/licenses/gpl-3.0.en.html)


trojan多用户管理部署程序

## 功能
- 在线web页面和命令行两种方式管理xray多用户
- 启动 / 停止 / 重启 xray 服务端
- 支持流量统计和流量限制
- 命令行模式管理, 支持命令补全
- 集成acme.sh证书申请
- 生成客户端配置文件
- 在线实时查看xray日志
- 支持xray://分享链接和二维码分享(二维码仅限web页面)
- 限制用户使用期限

## 安装方式
*xray使用请提前准备好服务器可用的域名*  

###  a. 一键脚本安装
```
#安装/更新
source <(curl -sL https://raw.githubusercontent.com/linuxing3/trojan/xray/install-xray.sh)

#卸载
source <(curl -sL https://raw.githubusercontent.com/linuxing3/trojan/xray/install-xray.sh) --remove

```
安装完后输入'xray'可进入管理程序   
浏览器访问 https://域名 可在线web页面管理xray用户  
前端页面源码地址: [xray-web](https://github.com/linuxing3/xray-web)

### b. docker运行
1. 安装mysql  

因为mariadb内存使用比mysql至少减少一半, 所以推荐使用mariadb数据库
```
docker run --name xray-mariadb --restart=always -p 3306:3306 -v /home/mariadb:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=xray -e MYSQL_ROOT_HOST=% -e MYSQL_DATABASE=xray -d mariadb:10.2
```
端口和root密码以及持久化目录都可以改成其他的

2. 安装xray
```
docker run -it -d --name xray --net=host --restart=always --privileged linuxing3/xray init
```
运行完后进入容器 `docker exec -it xray bash`, 然后输入'xray'即可进行初始化安装   

启动web服务: `systemctl start xray-web`   

设置自启动: `systemctl enable xray-web`

更新管理程序: `source <(curl -sL https://raw.githubusercontent.com/linuxing3/trojan/xray/install-xray.sh)`

## 运行截图
![avatar](asset/1.png)
![avatar](asset/2.png)

## 命令行
```
Usage:
  xray [flags]
  xray [command]

Available Commands:
  add         添加用户
  clean       清空指定用户流量
  completion  自动命令补全(支持bash和zsh)
  del         删除用户
  help        Help about any command
  info        用户信息列表
  log         查看xray日志
  restart     重启xray
  start       启动xray
  status      查看xray状态
  stop        停止xray
  tls         证书安装
  update      更新xray
  updateWeb   更新xray管理程序
  version     显示版本号
  web         以web方式启动

Flags:
  -h, --help   help for xray
```