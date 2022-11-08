# x-ui

 xray  panel with multi-protocol multi-user support

# Function introduction

- System status monitoring
- Support multi-user multi-protocol, web page visualization operation
- Supported protocols: vmess, vless, trojan, shadowsocks, dokodemo-door, socks, http
- Support to configure more transport configurations
-  Traffic statistics, limit traffic, limit expiration time
- Customizable xray Configuration templates
- Support https Access panel (self-provided domain name + ssl certificate)
- Support one-click SSL certificate application and automatic renewal
- More advanced configuration items, see panel for details

# Install & Upgrade

```
bash <(curl -Ls https://raw.githubusercontent.com/msameim181/x-ui/main/install.sh)
```

## Manual install & upgrade

1. First download the latest compressed package from https://github.com/msameim181/x-ui/releases , generally choose `amd64` architecture
2. Then upload the compressed package to the  `/root/` directory of the server, and use the  `root` user to log in to the server

> If your server cpu  architecture is not  `amd64`, replace  `amd64` in the command with another architecture

```
cd /root/
rm x-ui/ /usr/local/x-ui/ /usr/bin/x-ui -rf
tar zxvf x-ui-linux-amd64.tar.gz
chmod +x x-ui/x-ui x-ui/bin/xray-linux-* x-ui/x-ui.sh
cp x-ui/x-ui.sh /usr/bin/x-ui
cp -f x-ui/x-ui.service /etc/systemd/system/
mv x-ui/ /usr/local/
systemctl daemon-reload
systemctl enable x-ui
systemctl restart x-ui
```

## Install using docker

> This docker tutorial and  docker image are provided by [Chasing66](https://github.com/Chasing66)

1. Install docker

```shell
curl -fsSL https://get.docker.com | sh
```

2. Install x-ui

```shell
mkdir x-ui && cd x-ui
docker run -itd --network=host \
    -v $PWD/db/:/etc/x-ui/ \
    -v $PWD/cert/:/root/cert/ \
    --name x-ui --restart=unless-stopped \
    enwaiax/x-ui:latest
```

> Build  own image

```shell
docker build -t x-ui .
```

## SSL certificate application

> This function and tutorial are provided by [FranzKafkaYu](https://github.com/FranzKafkaYu)

The script has a built-in SSL certificate application function. To use this script to apply for a certificate, the following conditions must be met:

- Know Cloudflare registered email
- Aware of Cloudflare Global API Key
- The domain name has been resolved to the current server through cloudflare

How to get Cloudflare Global API Key:
    ![](media/bda84fbc2ede834deaba1c173a932223.png)
    ![](media/d13ffd6a73f938d1037d0708e31433bf.png)

When using, just enter  `domain name`, `mailbox`, `API KEY`, the diagram is as follows:
        ![](media/2022-04-04_141259.png)

Precautions:

- The script uses DNS API for certificate request
- Use Let'sEncrypt as the CA by default
- The certificate installation directory is the /root/cert directory
- The certificates applied for by this script are all generic domain name certificates

## Tg robot use (under development, temporarily unavailable)

> This function and tutorial are provided by [FranzKafkaYu](https://github.com/FranzKafkaYu)

X-UI supports daily traffic notification, panel login reminder and other functions through Tg robot. To use Tg robot, you need to apply by yourself
For specific application tutorials, please refer to [blog link](https://coderfan.net/how-to-use-telegram-bot-to-alarm-you-when-someone-login-into-your-vps.html)
Instructions for use: Set robot-related parameters in the background of the panel, including

- Tg Robot Token
- Tgbot ChatId
- Tg Robot cycle running time, using crontab syntax  

Reference syntax:
- 30 * * * * * //Notify at the 30s of every point
- @hourly      //Hourly notification
- @daily       //Daily notification (0:00 AM)
- @every 8h    //Notify every 8 hours  

TG notification content:
- Node traffic usage
- Panel login reminder
- Node expiration reminder
- Traffic warning reminder  

More features are planned...
## Suggestion system

- CentOS 7+
- Ubuntu 16+
- Debian 8+

# FAQ

## Migrated from  v2-ui 

First install the latest version of  x-ui on the server where  v2-ui  is installed, and then use the following command to migrate, which will migrate  `all  inbound  account data` of the local  v2-ui  to  x-ui, `Panel settings and username passwords are not migrated`

> After the migration is successful, please `close  v2-ui` and  `restart  x-ui`, otherwise the  inbound  of  v2-ui  will be the same as the  inbound  of  x-ui    will generate   `port conflict`

```
x-ui v2-ui
```

## issue Close

All kinds of white problems see high blood pressure

## Stargazers over time

[![Stargazers over time](https://starchart.cc/vaxilu/x-ui.svg)](https://starchart.cc/vaxilu/x-ui)
