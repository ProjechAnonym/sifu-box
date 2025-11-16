# sifu-box

sing-box 转换程序

## **注意**

sifu-box@2.0版本运行逻辑与v1.0版本区别较大, v1.0的部署方案不适用于v2.0版本。v2.0版本采用信号沟通sing-box内核, 不再需要配置sing-box的systemctl文件, 只需要配置sifu-box的就可以, sing-box上传到工作目录的core文件夹下 

## 安装

将压缩包下载之后解压即可

### 命令

```bash
apt-get update
apt-get install -y tar sudo vim acl
mkdir /opt/sifu-box
# 确保存在opt/sifubox文件夹,压缩包上传到root文件夹下,如果不是root用户可以改成绝对路径
tar -zxvf sifu-box-*.tar.gz -C /opt/sifubox --strip-components=1

# 删除压缩包
rm -rf sifu-box-*.tar.gz

```

- 以 root 用户运行的 systemctl 文件

```bash
cat > /etc/systemd/system/sifu-box.service <<EOF
[Unit]
Description=A config file transform Service
After=network.target

[Service]
Type=simple
ExecStart=/opt/sifu-box/bin/sifu-box run -c /opt/sifu-box/config.yaml -d /opt/sifu-box/lib -l 0.0.0.0:8080
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

- 以一般用户启动服务

```bash
useradd -r -s /usr/sbin/nologin sifubox
chown -R sifubox /opt/sifu-box
chgrp -R sifubox /opt/sifu-box
chmod -R 755 /opt/sifu-box
chmod u+x /opt/sifu-box/bin/sifu-box

cat > /usr/lib/systemd/system/sifu-box.service <<EOF
[Unit]
Description=A config file transform Service
After=network.target nss-lookup.target network-online.target

[Service]

CapabilityBoundingSet=CAP_NET_ADMIN CAP_NET_BIND_SERVICE CAP_SYS_PTRACE CAP_DAC_READ_SEARCH
AmbientCapabilities=CAP_NET_ADMIN CAP_NET_BIND_SERVICE CAP_SYS_PTRACE CAP_DAC_READ_SEARCH

User=sifubox
Group=sifubox
ExecStart=/opt/sifu-box/bin/sifu-box run -c /opt/sifu-box/config.yaml -d /opt/sifu-box/lib -l 0.0.0.0:8080
Restart=on-failure
RestartSec=10s
[Install]
WantedBy=multi-user.target
EOF

```

关于 sing-box 和 mosdns 的配置有时效问题,请移步博客[sing-box 和 mosdns 配置](https://vercel-blog.sifulin.top/zh-cn/2024/07/11/two-sexy-bitches-singbox-and-mosdns/)

### 运行命令讲解

sifu-box 目前仅接受 run 命令,参数如下

1. -c `该参数指定配置文件的路径`
2. -d `该参数指定工作目录, 数据库以及日志还有sing-box的lib目录会在该目录下`
3. -l `监听的地址`

### 必备文件

- 工作文件目录, `需要保证该目录具有读写权限`

| bin            | core           | lib                | logs               | sifu-box.db      | static                 | temp                   |
| -------------- | -------------- | ------------------ | ------------------ | ---------------- | ---------------------- | ---------------------- |
| sifu-box的目录 | sing-box的目录 | sing-box的工作目录 | sifu-box的日志目录 | sifu-box的数据库 | sifu-box的前端页面文件 | sifu-box的临时存储目录 |

```
.
|-- bin
|   `-- sifu-box
|-- core
|   `-- sing-box
|-- lib
|   `-- ui
|-- logs
|-- sifu-box.db
|-- static
|   `-- dist
`-- temp
    |-- config
    |-- providers
    `-- rulesets
```

1. `setting.config.yaml` 配置文件

| SMTP                                       | USER                                                         | TEMPLATE                                                     |
| ------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| 邮件服务器配置, 未来或许会用于添加hook通知 | 用户信息, username和password是管理员账户; code是访客代码, key则是JWT的加密密钥 | 默认的模板, 因为该程序的模板大部分字段需要自己填写, 有一个这个默认模板比较方便 |

```yaml
smtp:
    host: "smtp.qq.com"
    port: 465
    email: "1982396@qq.com"
    password: "123456"
user:
    username: "sifulin"
    password: "wsr19990902"
    code: "123456"
    key: "syius"
template:
    ntp:
      enabled: false
    inbounds:
      - address:
          - 172.18.0.1/30
          - fdfe:dcba:9876::1/126
        auto_route: true
        interface_name: tun0
        mtu: 9000
        stack: mixed
        strict_route: true
        tag: tun-in
        type: tun
    outbounds_group:
      - type: direct
        tag: direct
      - type: selector
        tag: select
        providers:
          - M78
        tag_groups:
          - auto
      - type: urltest
        tag: auto
        providers:
          - M78
    dns:
      final: external
      strategy: prefer_ipv4
      reverse_mapping: true
      rules:
        - action: route
          domain_keyword:
            - siub.top
          server: internal
      servers:
        - detour: select
          server: 8.8.8.8
          server_port: 853
          tag: external
          type: tls
        - server: 192.168.10.3
          server_port: 5335
          tag: internal
          type: udp
    experiment:
      clash_api:
        external_controller: 192.168.10.6:8080
        external_ui: ui
        secret: "123456"
    log:
      disabled: true
      timestamp: false
    route:
      rules:
        - action: route
          outbound: direct
          user:
            - bind
        - action: hijack-dns
          port:
            - 53
        - action: hijack-dns
          protocol:
            - dns
        - action: route
          ip_is_private: true
          outbound: direct
        - action: reject
          protocol:
            - quic
      rule_set:
        - type: remote
          tag: geoip-cn
          format: binary
          url: https://raw.githubusercontent.com/SagerNet/sing-geoip/refs/heads/rule-set/geoip-cn.srs
          download_detour: select
          update_interval: 1d
      final: select
      auto_detect_interface: true
      default_domain_resolver:
        server: internal
    providers:
      - M78
```

## 特性

1. 支持自动定时更新配置文件,默认每周一的 4:30 更新
2. 支持根据机场链接,已有 yaml 文件生成配置文件
3. 支持添加规则集到生成的配置文件
4. 支持自定义不同的模板生成配置文件并托管这些文件,其他设备轻松获取配置文件
5. 集成 Yacd 面板功能

## 结尾

[API 文档](https://github.com/ProjechAnonym/sifu-box/blob/main/API%E6%96%87%E6%A1%A3.md)我没细写, 就是直接导出来的, 开发过程中可能会有更改, 如果想编写新的 ui 建议自行测试返回结果

再来点规则集的下载链接

1. [MetaCubeX](https://github.com/MetaCubeX/meta-rules-dat/tree/sing)
2. [官方规则集](https://github.com/SagerNet/sing-geoip/tree/rule-set)
3. [保证量大管饱规则集](https://github.com/senshinya/singbox_ruleset/tree/main/rule)
