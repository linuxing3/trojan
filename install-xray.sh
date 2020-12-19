#!/bin/bash
# Author: Jrohy
# github: https://github.com/Jrohy/trojan

#定义操作变量, 0为否, 1为是
HELP=0

REMOVE=0

UPDATE=0

DOWNLAOD_URL="https://github.com/Jrohy/trojan/releases/download/"

VERSION_CHECK="https://api.github.com/repos/Jrohy/trojan/releases/latest"

SERVICE_URL="https://raw.githubusercontent.com/Jrohy/trojan/master/asset/xray-web.service"

[[ -e /var/lib/xray-manager ]] && UPDATE=1

#Centos 临时取消别名
[[ -f /etc/redhat-release && -z $(echo $SHELL|grep zsh) ]] && unalias -a

[[ -z $(echo $SHELL|grep zsh) ]] && SHELL_WAY="bash" || SHELL_WAY="zsh"

#######color code########
RED="31m"
GREEN="32m"
YELLOW="33m"
BLUE="36m"
FUCHSIA="35m"

colorEcho(){
    COLOR=$1
    echo -e "\033[${COLOR}${@:2}\033[0m"
}

#######get params#########
while [[ $# > 0 ]];do
    KEY="$1"
    case $KEY in
        --remove)
        REMOVE=1
        ;;
        -h|--help)
        HELP=1
        ;;
        *)
                # unknown option
        ;;
    esac
    shift # past argument or value
done
#############################

help(){
    echo "bash $0 [-h|--help] [--remove]"
    echo "  -h, --help           Show help"
    echo "      --remove         remove trojan"
    return 0
}

removeXray() {
    #移除trojan
    rm -rf /usr/bin/xray >/dev/null 2>&1
    rm -rf /usr/local/etc/xray >/dev/null 2>&1
    rm -f /etc/systemd/system/xray.service >/dev/null 2>&1

    #移除xray管理程序
    rm -f /usr/local/bin/xray >/dev/null 2>&1
    rm -rf /var/lib/xray-manager >/dev/null 2>&1
    rm -f /etc/systemd/system/xray-web.service >/dev/null 2>&1

    systemctl daemon-reload

    #移除xray的专用db
    docker rm -f xray-mysql xray-mariadb >/dev/null 2>&1
    rm -rf /home/mysql /home/mariadb >/dev/null 2>&1
    
    #移除环境变量
    sed -i '/xray/d' ~/.${SHELL_WAY}rc
    source ~/.${SHELL_WAY}rc

    colorEcho ${GREEN} "uninstall success!"
}

checkSys() {
    #检查是否为Root
    [ $(id -u) != "0" ] && { colorEcho ${RED} "Error: You must be root to run this script"; exit 1; }
    if [[ $(uname -m 2> /dev/null) != x86_64 ]]; then
        colorEcho $YELLOW "Please run this script on x86_64 machine."
        exit 1
    fi

    if [[ `command -v apt-get` ]];then
        PACKAGE_MANAGER='apt-get'
    elif [[ `command -v dnf` ]];then
        PACKAGE_MANAGER='dnf'
    elif [[ `command -v yum` ]];then
        PACKAGE_MANAGER='yum'
    else
        colorEcho $RED "Not support OS!"
        exit 1
    fi

    # 缺失/usr/local/bin路径时自动添加
    [[ -z `echo $PATH|grep /usr/local/bin` ]] && { echo 'export PATH=$PATH:/usr/local/bin' >> /etc/bashrc; source /etc/bashrc; }
}

#安装依赖
installDependent(){
    if [[ ${PACKAGE_MANAGER} == 'dnf' || ${PACKAGE_MANAGER} == 'yum' ]];then
        ${PACKAGE_MANAGER} install socat bash-completion -y
    else
        ${PACKAGE_MANAGER} update
        ${PACKAGE_MANAGER} install socat bash-completion xz-utils -y
    fi
}

setupCron() {
    if [[ `crontab -l 2>/dev/null|grep acme` ]]; then
        if [[ -z `crontab -l 2>/dev/null|grep xray-web` || `crontab -l 2>/dev/null|grep xray-web|grep "&"` ]]; then
            #计算北京时间早上3点时VPS的实际时间
            ORIGIN_TIME_ZONE=$(date -R|awk '{printf"%d",$6}')
            LOCAL_TIME_ZONE=${ORIGIN_TIME_ZONE%00}
            BEIJING_ZONE=8
            BEIJING_UPDATE_TIME=3
            DIFF_ZONE=$[$BEIJING_ZONE-$LOCAL_TIME_ZONE]
            LOCAL_TIME=$[$BEIJING_UPDATE_TIME-$DIFF_ZONE]
            if [ $LOCAL_TIME -lt 0 ];then
                LOCAL_TIME=$[24+$LOCAL_TIME]
            elif [ $LOCAL_TIME -ge 24 ];then
                LOCAL_TIME=$[$LOCAL_TIME-24]
            fi
            crontab -l 2>/dev/null|sed '/acme.sh/d' > crontab.txt
            echo "0 ${LOCAL_TIME}"' * * * systemctl stop xray-web; "/root/.acme.sh"/acme.sh --cron --home "/root/.acme.sh" > /dev/null; systemctl start xray-web' >> crontab.txt
            crontab crontab.txt
            rm -f crontab.txt
        fi
    fi
}

installXray(){
    local SHOW_TIP=0
    if [[ $UPDATE == 1 ]];then
        systemctl stop xray-web >/dev/null 2>&1
        rm -f /usr/local/bin/xray
    fi
    LASTEST_VERSION=$(curl -H 'Cache-Control: no-cache' -s "$VERSION_CHECK" | grep 'tag_name' | cut -d\" -f4)

    echo "正在下载管理程序`colorEcho $BLUE $LASTEST_VERSION`版本..."
    # 这里的xray不是xray的原始版本，而是本软件
    curl -L "$DOWNLAOD_URL/$LASTEST_VERSION/xray" -o /usr/local/bin/xray
    chmod +x /usr/local/bin/xray

    if [[ ! -e /etc/systemd/system/xray-web.service ]];then
        SHOW_TIP=1
        curl -L $SERVICE_URL -o /etc/systemd/system/xray-web.service
        systemctl daemon-reload
        systemctl enable xray-web
    fi
    #命令补全环境变量
    [[ -z $(grep xray ~/.${SHELL_WAY}rc) ]] && echo "source <(xray completion ${SHELL_WAY})" >> ~/.${SHELL_WAY}rc
    source ~/.${SHELL_WAY}rc
    if [[ $UPDATE == 0 ]];then
        colorEcho $GREEN "安装xray管理程序成功!\n"
        echo -e "运行命令`colorEcho $BLUE xray`可进行xray管理\n"
        /usr/local/bin/xray
    else
        if [[ `cat /usr/local/etc/xray/config.json|grep -w "\"db\""` ]];then
            sed -i "s/\"db\"/\"database\"/g" /usr/local/etc/xray/config.json
            systemctl restart xray
        fi
        /usr/local/bin/xray upgrade db
        if [[ -z `cat /usr/local/etc/xray/config.json|grep sni` ]];then
            /usr/local/bin/xray upgrade config
        fi
        systemctl restart xray-web
        colorEcho $GREEN "更新xray管理程序成功!\n"
    fi
    setupCron
    [[ $SHOW_TIP == 1 ]] && echo "浏览器访问'`colorEcho $BLUE https://域名`'可在线xray多用户管理"
}

main(){
    [[ ${HELP} == 1 ]] && help && return
    [[ ${REMOVE} == 1 ]] && removeXray && return
    [[ $UPDATE == 0 ]] && echo "正在安装xray管理程序.." || echo "正在更新xray管理程序.."
    checkSys
    [[ $UPDATE == 0 ]] && installDependent
    echo "安装xray的加强管理程序"
    installXray
}

main