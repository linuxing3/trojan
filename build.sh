#!/bin/bash
echo "开始打包"
packr2
go build -ldflags "-s -w -X 'trojan/xray.MVersion=`git describe --tags $(git rev-list --tags --max-count=1)`' -X 'trojan/xray.BuildDate=`TZ=Asia/Shanghai date "+%Y%m%d-%H%M"`' -X 'trojan/xray.GoVersion=`go version|awk '{print $3,$4}'`' -X 'trojan/xray.GitVersion=`git rev-parse HEAD`'" -o "result/xxxray" .
chmod +x result/xxxray
scp result/xxxray root@dongxishijie.xyz:xxxray
echo "打开Mysql服务器方法如下："
echo docker exec -it xray-mariadb mysql -u root -p
echo use xray;
echo show tables;
echo select * from users;
echo "清理临时文件"
packr2 clean
rm -rf result
