# x-ui
支持多协议多用户的 xray 面板

# 功能介绍
- 系统状态监控
- 支持多用户多协议，网页可视化操作
- 支持的协议：vmess、vless、trojan、shadowsocks、dokodemo-door、socks、http
- 支持配置更多传输配置
- 账号流量统计
- 可自定义 xray 配置模板
- 支持 https 访问面板（自备域名 + ssl 证书）
- 更多高级配置项，详见面板

# 安装&升级
```
bash <(curl -Ls https://raw.githubusercontent.com/sprov065/x-ui/master/install.sh) 0.2.0
```

## 建议系统
- CentOS 7+
- Ubuntu 16+
- Debian 8+

# 常见问题
## 与 v2-ui 关系
x-ui 相当于 v2-ui 的加强版，未来会加入更多功能，待 x-ui 功能稳定后，v2-ui 将不再提供更新

x-ui 可与 v2-ui 并存，数据不互通，不影响对方的运行

## 从 v2-ui 迁移
首先在安装了 v2-ui 的服务器上安装最新版 x-ui，然后使用以下命令进行迁移，将迁移本机 v2-ui 的`所有 inbound 账号数据`至 x-ui，`面板设置和用户名密码不会迁移`
> 迁移成功后请`关闭 v2-ui` 并且`重启 x-ui`，否则 v2-ui 的 inbound 会与 x-ui 的 inbound 会产生`端口冲突`
```
x-ui v2-ui
```

# Telegram
群组：https://t.me/sprov_blog

频道：https://t.me/sprov_channel

## Stargazers over time

[![Stargazers over time](https://starchart.cc/sprov065/x-ui.svg)](https://starchart.cc/sprov065/x-ui)
