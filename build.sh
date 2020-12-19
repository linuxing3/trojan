#!/bin/bash

PROJECT="linuxing3/trojan"

packr2

go build -ldflags "-s -w -X 'trojan/xray.MVersion=`git describe --tags $(git rev-list --tags --max-count=1)`' -X 'trojan/xray.BuildDate=`TZ=Asia/Shanghai date "+%Y%m%d-%H%M"`' -X 'trojan/xray.GoVersion=`go version|awk '{print $3,$4}'`' -X 'trojan/xray.GitVersion=`git rev-parse HEAD`'" -o "result/xxxray" .

scp result/xxxray root@dongxishijie.xyz:xxxray

# docker exec -it xray-mariadb mysql -u root -p

# use xray;

# show tables;

# select * from users;

packr2 clean

rm -rf result
