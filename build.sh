#!/bin/bash

PROJECT="linuxing3/trojan"

#获取当前的这个脚本所在绝对路径
SHELL_PATH=$(cd `dirname $0`; pwd)

cd $SHELL_PATH

packr2

go build -ldflags "-s -w -X 'trojan/xray.MVersion=`git describe --tags $(git rev-list --tags --max-count=1)`' -X 'trojan/xray.BuildDate=`TZ=Asia/Shanghai date "+%Y%m%d-%H%M"`' -X 'trojan/xray.GoVersion=`go version|awk '{print $3,$4}'`' -X 'trojan/xray.GitVersion=`git rev-parse HEAD`'" -o "result/xray" .

cd result

cd $SHELL_PATH

packr2 clean

rm -rf result
