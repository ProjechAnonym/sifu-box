# sifu-box

singbox 转换程序

## 安装

将压缩包下载之后解压即可

### 命令

**建议以高权限用户运行,因为该程序通过 systemctl 控制 singbox 的状态,权限不够会无法使用 systemctl**

```bash
apt-get update
apt-get install -y tar sudo vim
# 确保存在opt/sifubox文件夹,压缩包上传到root文件夹下,如果不是root用户可以改成绝对路径
tar -xvf sifu-box-*.tar --strip-components 1 -C /opt/sifubox/
cat > /etc/systemd/system/sifu-box.service <<EOF
[Unit]
Description=A config file transform Service
After=network.target

[Service]
Type=simple
ExecStart=/opt/sifubox/sifu-box
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload
# 开机自启
systemctl enable sifu-box
# 启动服务
systemctl start sifu-box
```

关于 sing-box 和 mosdns 的配置有时效问题,请移步博客[sing-box 和 mosdns 配置](https://vercel-blog.sifulin.top/zh-cn/2024/07/11/two-sexy-bitches-singbox-and-mosdns/)

### 路径配置

sifu-box 路径要求为在`/opt/singbox/bin/sing-box`,不在这个路径的话 sing-box 自动升级功能是无法使用的。生成的 singbox 配置文件默认设置在`/opt/singbox/config.json`路径下,此外 singbox 的启动关闭是通过`systemctl`命令控制的,务必确保配置了系统服务

### 必备文件

```

.
|-- dist
|-- config
| |-- proxy.config.yaml
| |-- recover.template.yaml
| `-- mode.config.yaml
`-- template
    `-- default.template.yaml

```

1. `config` 目录下存放配置文件
   proxy.config.yaml 为代理配置文件，mode.config.yaml 为服务器配置文件
   其中 proxy.config.yaml 文件内容如下:

```yaml
providers: [
    # 订阅链接列表,每个订阅链接应包含如下几项
    {
      path: "https://sub2.smallstrawberry.com/api/v1/client/subscribe?toke", # 订阅链接
      proxy: true, # 是否使用代理下载配置文件,仅服务模式有效
      name: 一速云, # 机场的名称
      remote: true, # 是否是远程订阅,如果是本地配置文件则path应为配置文件的绝对路径
    },
  ]
rulesets: [
    # 规则集列表,每个规则集应包含如下几项
    {
      tag: chatgpt, # 规则集的名称
      type: local, # 远程规则集还是本地规则集,可选值: local, remote
      path: /opt/singbox/chatgpt.json, # 规则集的路径,仅type为local时有效
      url: https://raw.githubusercontent.com/SagerNet/sing-geoip/rule-set/geoip-cn.srs, # 规则集的链接,仅type为remote时有效
      format: binary, # 规则集文件的格式,可选值: binary, source,
      china: false, # 是否中国地区规则集,如果为true,则在singbox配置中会直连出站
      dnsRule: "external", # 与singbox的dns搭配使用,命中该规则集的dns请求会从指定的dns服务器出站
      label: china, # 规则集的组标签,与tag不同,这个是搭配使用的,比如china-ip和china-site会在route中共同组成一个规则出站
      download_detour: select, # 下载时使用的出站策略,默认为select
      update_interval: 1d, # 更新时间间隔,默认为1天
    },
  ]
```

mode.config.yaml 文件内容如下:

```yaml
cors:
  - "*" # 允许所有跨域请求,为了安全性可以更改为你访问的域名,如果只是内网使用则无所谓
token: sifu # 前端登录密码,默认为sifu
key: $199wsr*dianhua1532# # 服务器模式可以根据不同模板文件生成不同的配置文件,比如ios的配置文件,为保证安全性参考机场的认证模式会将这段token进行MD5加密放入url参数中
mode: true # 是否为服务器模式
listen: "[::]:8080" # 服务器模式监听的端口,可以指定监听的ip和端口,比如[192.168.1.1]:9090,默认监听本机所有ip
```

2. `template` 目录下存放模板文件
   默认应该具备 default.template.yaml 模板文件,大部分内容参考 [singbox 的官方 wiki](https://sing-box.sagernet.org/zh/configuration/),部分有区别的模块说明如下:

```yaml
customOutbounds: [
    # 自定义节点的列表,这部分是用于添加自建节点的信息的,比如自建shadowsocks节点
    {
      type: shadowsocks,
      tag: "自建香港",
      server: "sifu.top",
      server_port: 0,
      method: "cipher",
      password: "wsr19990902",
    },
  ]
```

4. `dist` 目录下则是前端静态文件,不需要修改

## 特性

1. 支持自动定时更新配置文件,默认每周一的 4:30 更新
2. 支持根据机场链接,已有 yaml 文件生成配置文件
3. 支持添加规则集到生成的配置文件
4. 支持统一管理多台主机的 singbox 配置,只要打开 ssh 配置好 22 端口
5. 支持自定义不同的模板生成配置文件并托管这些文件,实现移动设备轻松获取配置文件
