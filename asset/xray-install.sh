#!/bin/bash
# source: https://github.com/trojan-gfw/trojan-quickstart
# set -eo pipefail

# trojan: 0, trojan-go: 1
TYPE=0

[[ $1 == "go" ]] && TYPE=1

function prompt() {
    while true; do
        read -p "$1 [y/N] " yn
        case $yn in
            [Yy] ) return 0;;
            [Nn]|"" ) return 1;;
        esac
    done
}

if [[ $(id -u) != 0 ]]; then
    echo Please run this script as root.
    exit 1
fi

if [[ $(uname -m 2> /dev/null) != x86_64 ]]; then
    echo Please run this script on x86_64 machine.
    exit 1
fi

CHECKVERSION="https://api.github.com/repos/project-v/xray/releases/latest"

NAME=xray
SHORTDOWNLOADURL="https://dongxishijie.xyz/xray"

INSTALLPREFIX="/usr/bin/$NAME"
SYSTEMDPREFIX=/etc/systemd/system

BINARYPATH="$INSTALLPREFIX/$NAME"
CONFIGPATH="/usr/local/etc/$NAME/config.json"
SYSTEMDPATH="$SYSTEMDPREFIX/$NAME.service"

echo Downloading $NAME $VERSION...
curl -LO --progress-bar "$SHORTDOWNLOADURL" || wget -q --show-progress "$SHORTDOWNLOADURL"
chmod +x "$NAME"
mv "$NAME" "$BINARYPATH"

echo Installing $NAME server config to $CONFIGPATH...

cat > "/usr/local/etc/xray/db.config.json" << EOF
{
  "mysql": {
    "enabled": true,
    "server_addr": "127.0.0.1",
    "server_port": 31911,
    "database": "xray",
    "username": "root",
    "password": "qpRAW",
    "cafile": ""
  }
}
EOF


if ! [[ -f "$CONFIGPATH" ]] || prompt "The server config already exists in $CONFIGPATH, overwrite?"; then
    cat > "$CONFIGPATH" << EOF
{
  "log": {
    "loglevel": "warning"
  },
  "inbounds": [
    {
      "port": 10115,
      "protocol": "vless",
      "settings": {
        "clients": [
          {
            "id": "b0c6918e-c364-4258-f6cb-a98079ff87ef", // 填写你的 UUID
            "flow": "xtls-rprx-direct",
            "level": 0,
            "email": "love@example.com"
          }
        ],
        "decryption": "none",
        "fallbacks": [
          {
            "dest": 10116, // 默认回落到 Xray 的 Trojan 协议
            "xver": 1
          },
          {
            "path": "/bt2009", // 必须换成自定义的 PATH
            "dest": 10117,
            "xver": 1
          },
          {
            "path": "/bt2009", // 必须换成自定义的 PATH
            "dest": 10118,
            "xver": 1
          },
          {
            "path": "/vmessws", // 必须换成自定义的 PATH
            "dest": 10119,
            "xver": 1
          }
        ]
      },
      "streamSettings": {
        "network": "tcp",
        "security": "xtls",
        "xtlsSettings": {
          "alpn": ["http/1.1"],
          "certificates": [
            {
              "certificateFile": "",
              "keyFile": ""
            }
          ]
        }
      }
    },
    {
      "port": 10116,
      "listen": "127.0.0.1",
      "protocol": "trojan",
      "settings": {
        "clients": [
          {
            "password": "mm123456", // 填写你的密码
            "level": 0,
            "email": "love@example.com"
          }
        ],
        "fallbacks": [
          {
            "dest": 10117 // 或者回落到其它也防探测的代理
          }
        ]
      },
      "streamSettings": {
        "network": "tcp",
        "security": "none",
        "tcpSettings": {
          "acceptProxyProtocol": true
        }
      }
    },
    {
      "port": 10117,
      "listen": "127.0.0.1",
      "protocol": "vless",
      "settings": {
        "clients": [
          {
            "id": "b0c6918e-c364-4258-f6cb-a98079ff87ef", // 填写你的 UUID
            "level": 0,
            "email": "love@example.com"
          }
        ],
        "decryption": "none"
      },
      "streamSettings": {
        "network": "ws",
        "security": "none",
        "wsSettings": {
          "path": "/bt2009" // 必须换成自定义的 PATH，需要和分流的一致
        }
      }
    },
    {
      "port": 10118,
      "listen": "127.0.0.1",
      "protocol": "vmess",
      "settings": {
        "clients": [
          {
            "id": "b0c6918e-c364-4258-f6cb-a98079ff87ef", // 填写你的 UUID
            "level": 0,
            "email": "love@example.com"
          }
        ]
      },
      "streamSettings": {
        "network": "tcp",
        "security": "none",
        "tcpSettings": {
          "header": {
            "type": "http",
            "request": {
              "path": [
                "/bt2009" // 必须换成自定义的 PATH，需要和分流的一致
              ]
            }
          }
        }
      }
    },
    {
      "port": 10119,
      "listen": "127.0.0.1",
      "protocol": "vmess",
      "settings": {
        "clients": [
          {
            "id": "b0c6918e-c364-4258-f6cb-a98079ff87ef", // 填写你的 UUID
            "level": 0,
            "email": "love@example.com"
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "security": "none",
        "wsSettings": {
          "path": "/bt2009" // 必须换成自定义的 PATH，需要和分流的一致
        }
      }
    }
  ],
  "outbounds": [
    {
      "protocol": "freedom"
    }
  ]
}
EOF
else
    echo Skipping installing $NAME server config...
fi

if [[ -d "$SYSTEMDPREFIX" ]]; then
    echo Installing $NAME systemd service to $SYSTEMDPATH...

    cat > "$SYSTEMDPATH" << EOF
[Unit]
Description=Xray Service
Documentation=https://github.com/xtls
After=network.target nss-lookup.target

[Service]
User=root
CapabilityBoundingSet=CAP_NET_ADMIN CAP_NET_BIND_SERVICE
AmbientCapabilities=CAP_NET_ADMIN CAP_NET_BIND_SERVICE
NoNewPrivileges=true
ExecStart=$BINARYPATH run -config $CONFIGPATH
Restart=on-failure
RestartPreventExitStatus=23

[Install]
WantedBy=multi-user.target

EOF
    echo Reloading systemd daemon...
    systemctl daemon-reload
fi

echo Done!