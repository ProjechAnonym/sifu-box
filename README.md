# sifu-box

singbox 转换程序

## 安装

将压缩包下载之后解压即可

### 命令

还没写再说

### 路径配置

sifu-box 本身的路径要求不严格,但是生成的 singbox 配置文件默认设置在`/opt/singbox/config.json`路径下,此外 singbox 的启动关闭是通过`systemctl`命令控制的,务必确保配置了系统服务

### 必备文件

```
.
|-- dist
|-- config
| |-- Proxy.config.yaml
| `-- Server.config.yaml
|-- static
|   `-- Default
`-- template
    `-- Default.template.yaml
```

1. `config` 目录下存放配置文件
   Proxy.config.yaml 为代理配置文件，Server.config.yaml 为服务器配置文件
   其中 Proxy.config.yaml 文件内容如下:

```yaml
url: [
    # 订阅链接列表,每个订阅链接应包含如下几项
    {
      path: "https://sub2.smallstrawberry.com/api/v1/client/subscribe?toke", # 订阅链接
      proxy: true, # 是否使用代理下载配置文件,仅服务模式有效
      label: 一速云, # 机场的名称
      remote: true, # 是否是远程订阅,如果否则path应为配置文件的绝对路径
    },
  ]
rule_set: [
    # 规则集列表,每个规则集应包含如下几项
    {
      label: chatgpt, # 规则集的名称
      value: {
          type: local, # 远程规则集还是本地规则集,可选值: local, remote
          path: /opt/singbox/chatgpt.txt, # 规则集的路径,如果type为local,则应为配置文件的绝对路径,如果type为remote,则应为订阅链接
          format: binary, # 规则集文件的格式,可选值: binary, source,
          china: false, # 是否中国地区规则集,如果为true,则在singbox配置中会直连出站
        },
    },
    { label: baidu, value: {
          type: remote,
          path: https://raw.githubusercontent.com/SagerNet/sing-geosite/rule-set/geosite-cn.srs,
          format: binary,
          china: true,
          # 以下两项仅type为remote时生效
          update_interval: 1d, # 更新时间间隔,默认为1天
          download_detour: select, ## 下载时使用的出站策略,默认为select
        } },
  ]
```

Server.config.yaml 文件内容如下:

```yaml
cors: { origins: ["*"] } # 允许所有跨域请求,为了安全性可以更改为你访问的域名,如果只是内网使用则无所谓
key: sifu # 前端登录密码,默认为sifu
token: $199wsr*dianhua1532# # 服务器模式可以根据不同模板文件生成不同的配置文件,比如ios的配置文件,为保证安全性参考机场的认证模式会将这段token进行MD5加密放入url参数中
server_mode: true # 是否为服务器模式,端口默认为8080
```

2.  `static`目录下存放生成的配置文件 3.`template` 目录下存放模板文件
    默认应该具备 Default.template.yaml 模板文件,大部分内容参考 [singbox 的官方 wiki](https://sing-box.sagernet.org/zh/configuration/),部分有区别的模块说明如下:

```yaml
outbounds:
  custom_outbound: [
      # 默认的出站,不配置上不了网的很基本的东西,此部分和wiki一致
      [
        { type: dns, tag: dns-out },
        { type: direct, tag: direct },
        { type: block, tag: block },
      ],
      # 自定义节点的列表,这部分是用于添加自建节点的信息的,比如自建shadowsocks节点
      [
        {
          type: shadowsocks,
          tag: "自建香港",
          server: "sifu.top",
          server_port: 0,
          method: "cipher",
          password: "wsr19990902",
        },
      ],
    ]
route:
  rule_set: [
      # 兜底用的geoip规则集,删除无法生成配置文件
      {
        tag: geoip-cn,
        type: remote,
        format: binary,
        url: https://raw.githubusercontent.com/SagerNet/sing-geoip/rule-set/geoip-cn.srs,
        download_detour: select,
        update_interval: 1d,
      },
      {
        tag: geosite-cn,
        type: remote,
        format: binary,
        url: https://raw.githubusercontent.com/SagerNet/sing-geosite/rule-set/geosite-cn.srs,
        download_detour: select,
        update_interval: 1d,
      },
    ]
  rules:
    # 默认的出站规则,就是一些很基本不调就上不了网的
    default:
      [
        { protocol: dns, outbound: dns-out },
        { ip_is_private: true, outbound: direct },
        { protocol: [quic], outbound: block },
      ]

    # 分流规则,这个配置会比较灵活,在Proxy配置中的规则集会在分流规则的第一条开始添加,然后才会是shunt规则,所以还可以在shunt添加一些自定义兜底分流规则
    shunt:
      [
        {
          type: logical,
          mode: and,
          rules:
            [
              { rule_set: geosite-cn, invert: true },
              { rule_set: geoip-cn, invert: true },
            ],
          outbound: select,
        },
        {
          type: logical,
          mode: and,
          rules: [{ rule_set: geosite-cn }, { rule_set: geoip-cn }],
          outbound: direct,
        },
      ]
```

4. `dist` 目录下则是前端静态文件,不需要修改

## 特性

1. 支持自动定时更新配置文件,默认每周一的 4:30 更新
2. 支持根据机场链接,已有 yaml 文件生成配置文件
3. 支持添加规则集到生成的配置文件
4. 支持统一管理多台主机的 singbox 配置,只要打开 ssh
