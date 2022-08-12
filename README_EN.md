# X-UI

[CN](./README.md)| EN  
X-UI is a webUI panel based on Xray-core which supports multi protocols and multi users  
This project is a fork of [vaxilu&#39;s project](https://github.com/vaxilu/x-ui),and it is a experiental project which used by myself for learning golang  
For some basic usages,please visit my [blog post](https://coderfan.net/how-to-use-x-ui-pannel-to-set-up-proxies-for-bypassing-gfw.html)  
If you need more language options ,please open a issue and let me know that

# changes 
- 2022.08.11：Support multi users on the same port;add CPU limit exceed  alert  
- 2022.07.28：Add acme standalone mode for cert issue；add  mechanism to keep X-UI alive even there exist crashes
- 2022.07.24：Add base path auto generate feature for security;add traffice reset automatically;add device alert
- 2022.07.21：Add more translations;add restart/stop xray service in Web panel
- 2022.07.11：Add time expiration notify for each inbound;add traffic limit notify for each inbound;add get url link command/inbound copy command in telegram bot  
- 2022.07.03：Add transport options in Trojan protocol;restruct Telegram bot for convenience  
- 2022.06.19：Add shadowsocks 2022 Ciphers,add inbounds search,traffic clear function in WebUI
- 2022.05.14：Add Telegram bot commands,support enable/disable/delete/status check
- 2022.04.25：Add SSH login notify
- 2022.04.23：Add WebUi login notify
- 2022.04.16：Add Telegram bot set up in WebUi pannel
- 2022.04.12：Optimize Telegram bot notify,more human friendly
- 2022.04.06：Add cert issue function，optimize installation/update and add telegram bot notify

# basics

- support system status info check
- support multi protocols and multi users
- support protocols：vmess、vless、trojan、shadowsocks、dokodemo-door、socks、http
- support many transport method including tcp、udp、ws、kcp etc
- traffic counting,traffic restrict and time restrcit
- support custom configuration template
- support https access fot WebUI
- support SSL cert issue by Acme
- support telegram bot notify and control
- more functions in control menu  

for more detailed usages,plz see [WIKI](https://github.com/FranzKafkaYu/x-ui/wiki)

# installation
Make sure your system `bash` and `curl` and `network` are ready,here we go

```
bash <(curl -Ls https://raw.githubusercontent.com/FranzKafkaYu/x-ui/master/install.sh)
```

## shortcut  
After Installation，you can input `x-ui`to enter control menu，current menu details：
```
  x-ui 面板管理脚本
  0. 退出脚本
————————————————
  1. 安装 x-ui
  2. 更新 x-ui
  3. 卸载 x-ui
————————————————
  4. 重置用户名密码
  5. 重置面板设置
  6. 设置面板端口
  7. 查看当前面板设置
————————————————
  8. 启动 x-ui
  9. 停止 x-ui
  10. 重启 x-ui
  11. 查看 x-ui 状态
  12. 查看 x-ui 日志
————————————————
  13. 设置 x-ui 开机自启
  14. 取消 x-ui 开机自启
————————————————
  15. 一键安装 bbr (最新内核)
  16. 一键申请SSL证书(acme申请)
 
面板状态: 已运行
是否开机自启: 是
xray 状态: 运行

请输入选择 [0-16]: 
```

## Suggested system as follows:
- CentOS 7+
- Ubuntu 16+
- Debian 8+

# telegram

[CoderfanBaby](https://t.me/CoderfanBaby)  
[FranzKafka‘sPrivateGroup](https://t.me/franzkafayu)

# credits
- [vaxilu/x-ui](https://github.com/vaxilu/x-ui)
- [XTLS/Xray-core](https://github.com/XTLS/Xray-core)
- [telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)

## Stargazers over time

[![Stargazers over time](https://starchart.cc/FranzKafkaYu/x-ui.svg)](https://starchart.cc/FranzKafkaYu/x-ui)
