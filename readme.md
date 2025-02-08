# sifu-box

singbox 转换程序

## **注意**

该程序依赖 sing-box,请先安装 sing-box,并且由于该程序依赖 systemctl,所以需要以高权限用户运行,同时自动切换 sing-box 的配置文件也需要对其有相应权限, 务必确保权限足够,若 sing-box 与 sifu-box 都使用 root 用户部署应该不会存在权限问题

## 安装

将压缩包下载之后解压即可

### 命令

**建议以高权限用户运行,因为该程序通过 systemctl 控制 singbox 的状态,权限不够会无法使用 systemctl**

```bash
apt-get update
apt-get install -y tar sudo vim acl
# 确保存在opt/sifubox文件夹,压缩包上传到root文件夹下,如果不是root用户可以改成绝对路径
tar -xvf sifu-box-*.tar --strip-components 1 -C /opt/sifubox/


```

- 以 root 用户运行的 systemctl 文件

```bash
cat > /etc/systemd/system/sifu-box.service <<EOF
[Unit]
Description=A config file transform Service
After=network.target

[Service]
Type=simple
ExecStart=/opt/sifubox/bin/sifu-box run -c /opt/sifubox/config -w /opt/test -l :8080 -s
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
chown -R sifubox /opt/sifubox
chgrp -R sifubox /opt/sifubox
# 因为sifu-box需要对sing-box的配置文件进行自动修改, 所以需要给sing-box的配置文件添加读写权限
setfacl -m u:sifubox:rw /opt/singbox/config.json
cat > /etc/systemd/system/sifu-box.service <<EOF
[Unit]
Description=A config file transform Service
After=network.target

[Service]
Type=simple
User=sifubox
Group=sifubox
ExecStart=/opt/sifubox/bin/sifu-box run -c /opt/sifubox/config -w /opt/test -l :8080 -s
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF
# 如果使用systemctl控制sing-box的启停, 则需要设置sudo权限, 一般的用户是没有权限执行systemctl的
cat >> /etc/sudoers.d/sifubox-nopasswd << EOF
sifubox ALL=(ALL) NOPASSWD: /bin/systemctl * sifu-box
EOF
```

关于 sing-box 和 mosdns 的配置有时效问题,请移步博客[sing-box 和 mosdns 配置](https://vercel-blog.sifulin.top/zh-cn/2024/07/11/two-sexy-bitches-singbox-and-mosdns/)

### 运行命令讲解

sifu-box 目前仅接受 run 命令,参数如下

1. -c `该参数指定配置文件的目录`
2. -w `该参数指定工作目录, 数据库以及日志等文件会在该目录下`
3. -l `监听的地址`
4. -s `是否启用服务模式, 出现该参数必须设置-l参数, 没有该参数则sifu-box直接读取配置文件中的机场规则集以及模板生成相应的配置文件在工作目录下`

### 必备文件

- 配置文件目录

```
.
|
`-- config
  `-- setting.config.yaml

```

- 工作文件目录, `需要保证该目录具有读写权限`

```
.
`-- static
  |-- dist
  `-- template
      `-- default.template.yaml
```

1. `config` 目录下存放配置文件
   setting.config.yaml 为配置文件

```yaml
application:
  singbox:
    listen: http://192.168.1.2:9090 # sing-box的clashAPI 监听地址
    secret: 123456 # sing-box的clashAPI 的密钥
    work_dir: /opt/singbox/config # sing-box的工作目录
    config_path: /opt/singbox/config.json # sing-box的配置文件路径
    binary_path: /opt/singbox/singbox # sing-box的文件路径
    commands: # 控制sing-box的命令, 一般用户要使用sudo
      boot_command:
        name: systemctl
        args:
          - start
          - sing-box
      stop_command:
        name: systemctl
        args:
          - stop
          - sing-box
      check_command:
        name: systemctl
        args:
          - status
          - sing-box
      reload_command:
        name: systemctl
        args:
          - reload
          - sing-box
      restart_command:
        name: systemctl
        args:
          - restart
          - sing-box
  server:
    user:
      username: sifulin # 管理员账户
      password: adminadmin # 管理员密码
      email: 19826@qq.com # 管理员邮箱
      code: 123456 # 非管理员用户登录时的密钥
      private_key: sifulin # 管理员私钥, 用于生成JWT密钥, 设个难点的就行
    ssl: # ssl证书,没有则可以删掉这个项目
      public: /opt/singbox/config/cert.pem
      private: /opt/singbox/config/key.pem

configuration: # 配置项, 如果是服务模式可不填写此项, 当然该处的配置会在程序运行时读取并写入sqlite数据库中
  providers: # 配置源
    - name: "夜煞云"
      path: "https://45.137.180.216/api/v1/client/nyth?token=aa751df806caf70136e33c9310531fe6"
      detour: select # 在生成的配置文件中, 路由规则会将该机场的域名出站设为这里设置的出站口
      remote: true # 是否位于本地, true为网络链接下载
  rulesets:
    - type: remote # 是否位于本地, remote为网络链接下载
      path: https://github.com/MetaCubeX/meta-rules-dat/raw/sing/geo/geosite/cn.srs
      format: binary # binary和source两种格式, 二进制格式更快
      tag: china-site # 标签, 用于匹配规则
      download_detour: select # 下载时使用的路由
      update_interval: 1d # 更新时间间隔
      label: china # 标签, 不同于tag, 相同label不同tag的规则集会在路由规则中集合起来共同使用一个出站口
      china: true # 是否位于中国大陆
      name_server: internal # 使用的dns, 可以不填, 如果是域名的话可以指定这组规则集的要使用的DNS, 留空则完全要看是否会匹配到DNS规则集指定的出站口, 没匹配到则会使用默认的DNS
    - type: remote
      path: https://github.com/MetaCubeX/meta-rules-dat/raw/sing/geo/geoip/cn.srs
      format: binary
      tag: china-ip
      download_detour: select
      update_interval: 1d
      label: china
      china: true
  templates:
    default:
      log:
        disabled: false
        timestamp: true
        output: singbox.log
      experimental:
        cache_file:
          enabled: true
        clash_api:
          external_controller: 0.0.0.0:9090
          external_ui: ui
          external_ui_download_detour: select
          external_ui_download_url: https://github.com/MetaCubeX/Yacd-meta/archive/gh-pages.zip
          secret: "123456"
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
      dns:
        final: external
        strategy: prefer_ipv4
        reverse_mapping: true
        servers:
          - tag: external
            type: tls
            server_port: 853
            server: "8.8.8.8"
            detour: select
          - tag: internal
            type: udp
            server: "192.168.1.2"
            server_port: 5353
            detour: direct
        rules: []
      route:
        # 默认的出站规则,就是一些很基本不调就上不了网的
        rules:
          - user:
              - bind
            action: route
            outbound: direct
          - action: hijack-dns
            port:
              - 53
          - action: hijack-dns
            protocol:
              - dns
          - action: route
            outbound: direct
            ip_is_private: true
          - protocol:
              - quic
            action: reject
        # 默认出站节点不能改,很多配置下载需要代理,我设置的都是select
        final: select
        auto_detect_interface: true
      outbounds:
        - tag: direct
          type: direct
      custom_outbounds:
        - tag: 香港
          type: shadowsocks
          method: chacha20-ietf-poly1305
          server: 127.0.0.1
          server_port: 443
          password: 123456dsaf
    ios:
      log:
        disabled: false
        timestamp: true
        output: singbox.log
      experimental:
        cache_file:
          enabled: true
        clash_api:
          external_controller: 0.0.0.0:9090
          external_ui: ui
          external_ui_download_detour: select
          external_ui_download_url: https://github.com/MetaCubeX/Yacd-meta/archive/gh-pages.zip
          secret: "123456"
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
      dns:
        final: external
        strategy: prefer_ipv4
        reverse_mapping: true
        servers:
          - tag: external
            type: tls
            server_port: 853
            server: "8.8.8.8"
            detour: select
          - tag: internal
            type: udp
            server: "192.168.1.2"
            server_port: 5353
            detour: direct
        rules: []
      route:
        # 默认的出站规则,就是一些很基本不调就上不了网的
        rules:
          - user:
              - bind
            action: route
            outbound: direct
          - action: hijack-dns
            port:
              - 53
          - action: hijack-dns
            protocol:
              - dns
          - action: route
            outbound: direct
            ip_is_private: true
          - protocol:
              - quic
            action: reject
        # 默认出站节点不能改,很多配置下载需要代理,我设置的都是select
        final: select
        auto_detect_interface: true
      outbounds:
        - tag: direct
          type: direct
      custom_outbounds: # 自定义出站节点
        - tag: 香港
          type: shadowsocks
          method: chacha20-ietf-poly1305
          server: 127.0.0.1
          server_port: 443
          password: 123456dsaf
```

4. `dist` 目录下则是前端静态文件,不需要修改

## 特性

1. 支持自动定时更新配置文件,默认每周一的 4:30 更新
2. 支持根据机场链接,已有 yaml 文件生成配置文件
3. 支持添加规则集到生成的配置文件
4. 支持自定义不同的模板生成配置文件并托管这些文件,其他设备轻松获取配置文件
5. 集成 Yacd 面板功能

## 结尾

API 文档我没细写, 就是直接导出来的, 开发过程中可能会有更改, 如果想编写新的 ui 建议自行测试返回结果
