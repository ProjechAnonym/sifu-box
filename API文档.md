---
title: 默认模块
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: true
code_clipboard: true
highlight_theme: darkula
headingLevel: 2
generator: "@tarslib/widdershins v4.0.30"

---

# 默认模块

Base URLs:

* <a href="http://prod-cn.your-api-server.com">正式环境: http://prod-cn.your-api-server.com</a>

* <a href="http://dev-cn.your-api-server.com">开发环境: http://dev-cn.your-api-server.com</a>

# Authentication

# 登录

## POST 游客登录接口

POST /api/login/visitor

> Body 请求参数

```yaml
code: "123456"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» code|body|string| 是 |none|

> 返回示例

```json
{
  "message": {
    "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsInV1aWQiOiI3NDBkZDZiZS04NmE4LTQwMDYtYmE3MC03ZWY5OTA0MTFhYzYiLCJleHAiOjE3NTg0NTk4NjF9.5CTn9HPkUbM9iYFbKcDvtnQQc2WpV4UsjJ_bvkO0uWA",
    "admin": false
  }
}
```

```json
{
  "message": "密钥错误"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST 管理员登录接口

POST /api/login/admin

> Body 请求参数

```yaml
username: sifulin
password: wsr199

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» username|body|string| 是 |none|
|» password|body|string| 是 |none|

> 返回示例

```json
{
  "message": {
    "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwidXVpZCI6ImMxODc0YWQ2LWU0YzEtNGRlMi05YzY1LTA1OGVhMGY5NTJkMSIsImV4cCI6MTc1ODQ1OTkxNn0.L7bpsr9cD0YHQQe3ULnCB1fUmUfXdnfsRsXOmead5FA",
    "admin": true
  }
}
```

```json
{
  "message": "用户名或密码错误"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## GET 自动登录接口

GET /api/verify

> Body 请求参数

```yaml
{}

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |none|
|body|body|object| 否 |none|

> 返回示例

```json
{
  "message": {
    "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwidXVpZCI6ImMxODc0YWQ2LWU0YzEtNGRlMi05YzY1LTA1OGVhMGY5NTJkMSIsImV4cCI6MTc1ODQ1OTkxNn0.L7bpsr9cD0YHQQe3ULnCB1fUmUfXdnfsRsXOmead5FA",
    "admin": true
  }
}
```

```json
{
  "message": {
    "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsInV1aWQiOiI1OTc2NTUwZS02YWZhLTQwMzUtYTc4Yy05YzE5MmQwODFlY2QiLCJleHAiOjE3NTg0NjAwMjN9.bZHsosiInkVPyEwB-WjUDRskmTPm3s7GdMP3Gpc3OYk",
    "admin": false
  }
}
```

```json
"{\"message\":\"解析\\\"authorization\"字段失败\"}"
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

# 配置

## POST 添加机场-云端

POST /api/configuration/add/provider/remote

> Body 请求参数

```json
[
  {
    "name": "M782",
    "path": "https://sub.m78sc.cn/api/v1/client/subscribe?token=083387dce0f02a10e8115379f9871c6d1",
    "remote": true
  }
]
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |none|
|body|body|array[object]| 否 |none|

> 返回示例

> 200 Response

```json
{
  "message": [
    {
      "message": "添加机场\"M78\"失败: [ent: constraint failed: UNIQUE constraint failed: providers.name, providers.path]",
      "status": false
    },
    {
      "message": "添加机场\"自定义1\"成功",
      "status": true
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## PATCH 修改机场

PATCH /api/configuration/edit/provider

> Body 请求参数

```yaml
name: M78
path: https://sub.m78sc.cn/api/v1/client/subscribe?token=083387dce0f02a10e8115379f9871c6d
remote: "true"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |none|
|body|body|object| 否 |none|
|» name|body|string| 否 |none|
|» path|body|string| 否 |none|
|» remote|body|boolean| 否 |none|

> 返回示例

```json
{
  "message": [
    {
      "message": "添加机场\"M78\"失败: [ent: constraint failed: UNIQUE constraint failed: providers.name, providers.path]",
      "status": false
    },
    {
      "message": "添加机场\"自定义1\"成功",
      "status": true
    }
  ]
}
```

```json
{
  "message": "修改机场\"M78\"成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST 添加规则集-云端

POST /api/configuration/add/ruleset/remote

> Body 请求参数

```json
"[\n    {\n        \"name\": \"china-ip\",\n        \"path\": \"https://github.com/SagerNet/sing-geoip/raw/refs/heads/rule-set/geoip-cn.srs\",\n        \"remote\": true,\n        \"binary\": true,\n        \"download_detour\": \"direct\",\n        \"update_interval\": \"1d\"\n    },\n    // {\n    //     \"name\": \"bing\",\n    //     \"path\": \"https://github.com/MetaCubeX/meta-rules-dat/raw/refs/heads/sing/geo/geosite/bing.srs\",\n    //     \"remote\": true,\n    //     \"binary\": true,\n    //     \"download_detour\": \"direct\",\n    //     \"update_interval\": \"1d\"\n    // }\n]"
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |none|
|body|body|array[object]| 否 |none|

> 返回示例

```json
{
  "message": [
    {
      "message": "添加机场\"M78\"失败: [ent: constraint failed: UNIQUE constraint failed: providers.name, providers.path]",
      "status": false
    },
    {
      "message": "添加机场\"自定义1\"成功",
      "status": true
    }
  ]
}
```

```json
{
  "message": [
    {
      "message": "添加机场\"china-ip\"成功",
      "status": true
    },
    {
      "message": "添加机场\"bing\"成功",
      "status": true
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## PATCH 修改规则集

PATCH /api/configuration/edit/ruleset

> Body 请求参数

```yaml
name: china-ip
path: https://github.com/SagerNet/sing-geoip/raw/refs/heads/rule-set/geoip-cn.srs
remote: "true"
binary: "true"
download_detour: direct
update_interval: 12d

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |none|
|body|body|object| 否 |none|
|» name|body|string| 否 |none|
|» path|body|string| 否 |none|
|» remote|body|boolean| 否 |none|
|» binary|body|boolean| 否 |none|
|» download_detour|body|string| 否 |none|
|» update_interval|body|string| 否 |none|

> 返回示例

```json
{
  "message": [
    {
      "message": "添加机场\"M78\"失败: [ent: constraint failed: UNIQUE constraint failed: providers.name, providers.path]",
      "status": false
    },
    {
      "message": "添加机场\"自定义1\"成功",
      "status": true
    }
  ]
}
```

```json
{
  "message": "修改机场\"M78\"成功"
}
```

```json
{
  "message": "修改机场\"M78\"成功"
}
```

```json
{
  "message": "修改机场\"china-ip\"成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## GET 获取配置

GET /api/configuration/fetch

> Body 请求参数

```yaml
{}

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |none|
|body|body|object| 否 |none|

> 返回示例

> 207 Response

```json
{
  "message": [
    {
      "message": [
        {
          "id": 1,
          "name": "M78",
          "path": "https://sub.m78sc.cn/api/v1/client/subscribe?token=083387dce0f02a10e8115379f9871c6d",
          "remote": true
        },
        {
          "id": 3,
          "name": "M718",
          "path": "/opt/sifubox/uploads/providers/e166c2105f66562d61a2afce952763fb.yaml",
          "remote": false
        },
        {
          "id": 4,
          "name": "M728",
          "path": "/opt/sifubox/uploads/providers/a09a24f2166195e220194252ac1386a7.yaml",
          "remote": false
        }
      ],
      "status": true,
      "type": "provider"
    },
    {
      "message": [
        {
          "id": 6,
          "name": "Bing",
          "path": "/opt/sifubox/uploads/rulesets/112a7e096595d1c32c4ecdfd9e56b66c.srs",
          "remote": false,
          "update_interval": "",
          "binary": true,
          "download_detour": ""
        },
        {
          "id": 7,
          "name": "geoip-cn",
          "path": "/opt/sifubox/uploads/rulesets/0b1e892fe64acae39ac4e5106f55474f.srs",
          "remote": false,
          "update_interval": "",
          "binary": true,
          "download_detour": ""
        },
        {
          "id": 9,
          "name": "china-ip",
          "path": "https://github.com/SagerNet/sing-geoip/raw/refs/heads/rule-set/geoip-cn.srs",
          "remote": true,
          "update_interval": "1d",
          "binary": true,
          "download_detour": "direct"
        }
      ],
      "status": true,
      "type": "ruleset"
    },
    {
      "message": [
        {
          "id": 1,
          "name": "default",
          "ntp": {
            "enabled": false,
            "server": ""
          },
          "inbounds": [
            {
              "address": [
                "172.18.0.1/30",
                "fdfe:dcba:9876::1/126"
              ],
              "auto_route": true,
              "interface_name": "tun0",
              "mtu": 1500,
              "stack": "mixed",
              "strict_route": true,
              "tag": "tun-in",
              "type": "tun"
            }
          ],
          "outbounds_group": [
            {
              "type": "direct",
              "tag": "direct",
              "providers": null,
              "tag_groups": null
            },
            {
              "type": "selector",
              "tag": "selector",
              "providers": [
                "M78"
              ],
              "tag_groups": null
            },
            {
              "type": "urltest",
              "tag": "auto",
              "providers": [
                "M78"
              ],
              "tag_groups": null
            }
          ],
          "dns": {
            "final": "google",
            "strategy": "prefer_ipv4",
            "servers": [
              {
                "server": "8.8.8.8",
                "server_port": 853,
                "tag": "google",
                "type": "tls"
              },
              {
                "server": "1.1.1.1",
                "server_port": 853,
                "tag": "cloudflare",
                "type": "tls"
              }
            ]
          },
          "experiment": {
            "clash_api": {
              "external_controller": "127.0.0.1:9090",
              "external_ui": "ui",
              "secret": "123456"
            }
          },
          "log": {
            "disabled": true,
            "output": "",
            "timestamp": false,
            "level": ""
          },
          "route": {
            "rules": [
              {
                "action": "route",
                "outbound": "direct",
                "user": [
                  "bind"
                ]
              },
              {
                "action": "hijack-dns",
                "port": [
                  53
                ]
              },
              {
                "action": "hijack-dns",
                "protocol": [
                  "dns"
                ]
              },
              {
                "action": "route",
                "ip_is_private": true,
                "outbound": "direct"
              },
              {
                "action": "reject",
                "protocol": [
                  "quic"
                ]
              }
            ],
            "final": "direct",
            "default_domain_resolver": {
              "server": "google"
            }
          },
          "providers": [
            "M78"
          ]
        },
        {
          "id": 2,
          "name": "default1",
          "ntp": {
            "enabled": false,
            "server": ""
          },
          "inbounds": [
            {
              "address": [
                "172.18.0.1/30",
                "fdfe:dcba:9876::1/126"
              ],
              "auto_route": true,
              "interface_name": "tun0",
              "mtu": 1500,
              "stack": "mixed",
              "strict_route": true,
              "tag": "tun-in",
              "type": "tun"
            }
          ],
          "outbounds_group": [
            {
              "type": "direct",
              "tag": "direct",
              "providers": null,
              "tag_groups": null
            },
            {
              "type": "selector",
              "tag": "selector",
              "providers": [
                "M78"
              ],
              "tag_groups": null
            },
            {
              "type": "urltest",
              "tag": "auto",
              "providers": [
                "M78"
              ],
              "tag_groups": null
            }
          ],
          "dns": {
            "final": "google",
            "strategy": "prefer_ipv4",
            "servers": [
              {
                "server": "8.8.8.8",
                "server_port": 853,
                "tag": "google",
                "type": "tls"
              },
              {
                "server": "1.1.1.1",
                "server_port": 853,
                "tag": "cloudflare",
                "type": "tls"
              }
            ]
          },
          "experiment": {
            "clash_api": {
              "external_controller": "127.0.0.1:9090",
              "external_ui": "ui",
              "secret": "123456"
            }
          },
          "log": {
            "disabled": true,
            "output": "",
            "timestamp": false,
            "level": ""
          },
          "route": {
            "rules": [
              {
                "action": "route",
                "outbound": "direct",
                "user": [
                  "bind"
                ]
              },
              {
                "action": "hijack-dns",
                "port": [
                  53
                ]
              },
              {
                "action": "hijack-dns",
                "protocol": [
                  "dns"
                ]
              },
              {
                "action": "route",
                "ip_is_private": true,
                "outbound": "direct"
              },
              {
                "action": "reject",
                "protocol": [
                  "quic"
                ]
              }
            ],
            "final": "direct",
            "default_domain_resolver": {
              "server": "google"
            }
          },
          "providers": [
            "M78"
          ]
        },
        {
          "id": 3,
          "name": "default2",
          "ntp": {
            "enabled": false,
            "server": ""
          },
          "inbounds": [
            {
              "address": [
                "172.18.0.1/30",
                "fdfe:dcba:9876::1/126"
              ],
              "auto_route": true,
              "interface_name": "tun0",
              "mtu": 1500,
              "stack": "mixed",
              "strict_route": true,
              "tag": "tun-in",
              "type": "tun"
            }
          ],
          "outbounds_group": [
            {
              "type": "direct",
              "tag": "direct",
              "providers": null,
              "tag_groups": null
            },
            {
              "type": "selector",
              "tag": "selector",
              "providers": [
                "M78"
              ],
              "tag_groups": [
                "auto"
              ]
            },
            {
              "type": "urltest",
              "tag": "auto",
              "providers": [
                "M78"
              ],
              "tag_groups": null
            }
          ],
          "dns": {
            "final": "google",
            "strategy": "prefer_ipv4",
            "servers": [
              {
                "server": "8.8.8.8",
                "server_port": 853,
                "tag": "google",
                "type": "tls"
              },
              {
                "server": "1.1.1.1",
                "server_port": 853,
                "tag": "cloudflare",
                "type": "tls"
              }
            ]
          },
          "experiment": {
            "clash_api": {
              "external_controller": "127.0.0.1:9090",
              "external_ui": "ui",
              "secret": "123456"
            }
          },
          "log": {
            "disabled": true,
            "output": "",
            "timestamp": false,
            "level": ""
          },
          "route": {
            "rules": [
              {
                "action": "route",
                "outbound": "direct",
                "user": [
                  "bind"
                ]
              },
              {
                "action": "hijack-dns",
                "port": [
                  53
                ]
              },
              {
                "action": "hijack-dns",
                "protocol": [
                  "dns"
                ]
              },
              {
                "action": "route",
                "ip_is_private": true,
                "outbound": "direct"
              },
              {
                "action": "reject",
                "protocol": [
                  "quic"
                ]
              }
            ],
            "final": "direct",
            "default_domain_resolver": {
              "server": "google"
            }
          },
          "providers": [
            "M78"
          ]
        }
      ],
      "status": true,
      "type": "template"
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|207|[Multi-Status](https://tools.ietf.org/html/rfc2518#section-10.2)|none|Inline|

### 返回数据结构

状态码 **207**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|[object]|true|none||none|
|»» message|[object]|true|none||none|
|»»» id|integer|true|none||none|
|»»» name|string|true|none||none|
|»»» path|string|true|none||none|
|»»» remote|boolean|true|none||none|
|»»» update_interval|string|true|none||none|
|»»» binary|boolean|true|none||none|
|»»» download_detour|string|true|none||none|
|»»» ntp|object|true|none||none|
|»»»» enabled|boolean|true|none||none|
|»»»» server|string|true|none||none|
|»»» inbounds|[object]|true|none||none|
|»»»» address|[string]|true|none||none|
|»»»» auto_route|boolean|true|none||none|
|»»»» interface_name|string|true|none||none|
|»»»» mtu|integer|true|none||none|
|»»»» stack|string|true|none||none|
|»»»» strict_route|boolean|true|none||none|
|»»»» tag|string|true|none||none|
|»»»» type|string|true|none||none|
|»»» outbounds_group|[object]|true|none||none|
|»»»» type|string|true|none||none|
|»»»» tag|string|true|none||none|
|»»»» providers|[string]|true|none||none|
|»»»» tag_groups|[string]¦null|true|none||none|
|»»» dns|object|true|none||none|
|»»»» final|string|true|none||none|
|»»»» strategy|string|true|none||none|
|»»»» servers|[object]|true|none||none|
|»»»»» server|string|true|none||none|
|»»»»» server_port|integer|true|none||none|
|»»»»» tag|string|true|none||none|
|»»»»» type|string|true|none||none|
|»»» experiment|object|true|none||none|
|»»»» clash_api|object|true|none||none|
|»»»»» external_controller|string|true|none||none|
|»»»»» external_ui|string|true|none||none|
|»»»»» secret|string|true|none||none|
|»»» log|object|true|none||none|
|»»»» disabled|boolean|true|none||none|
|»»»» output|string|true|none||none|
|»»»» timestamp|boolean|true|none||none|
|»»»» level|string|true|none||none|
|»»» route|object|true|none||none|
|»»»» rules|[object]|true|none||none|
|»»»»» action|string|true|none||none|
|»»»»» outbound|string|true|none||none|
|»»»»» user|[string]|true|none||none|
|»»»»» port|[integer]|true|none||none|
|»»»»» protocol|[string]|true|none||none|
|»»»»» ip_is_private|boolean|true|none||none|
|»»»» final|string|true|none||none|
|»»»» default_domain_resolver|object|true|none||none|
|»»»»» server|string|true|none||none|
|»»» providers|[string]|true|none||none|
|»» status|boolean|true|none||none|
|»» type|string|true|none||none|

## DELETE 删除机场

DELETE /api/configuration/delete/provider

> Body 请求参数

```yaml
name:
  - M718
  - M728

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |none|
|body|body|object| 否 |none|
|» name|body|[string]| 是 |名称|

> 返回示例

```json
{
  "message": [
    {
      "message": "添加机场\"M78\"失败: [ent: constraint failed: UNIQUE constraint failed: providers.name, providers.path]",
      "status": false
    },
    {
      "message": "添加机场\"自定义1\"成功",
      "status": true
    }
  ]
}
```

```json
[
  {
    "message": "删除机场\"自定义\"成功",
    "status": true
  },
  {
    "message": "删除机场\"自定义1\"成功",
    "status": true
  }
]
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|207|[Multi-Status](https://tools.ietf.org/html/rfc2518#section-10.2)|none|Inline|

### 返回数据结构

状态码 **207**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» status|boolean|true|none||none|

## DELETE 删除规则集

DELETE /api/configuration/delete/ruleset

> Body 请求参数

```yaml
name:
  - china-ip

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |none|
|body|body|object| 否 |none|
|» name|body|[string]| 是 |名称|

> 返回示例

> 207 Response

```json
[
  {
    "message": "删除规则集\"china-ip\"成功",
    "status": true
  }
]
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|207|[Multi-Status](https://tools.ietf.org/html/rfc2518#section-10.2)|none|Inline|

### 返回数据结构

状态码 **207**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» status|boolean|true|none||none|

## POST 添加机场-文件

POST /api/configuration/add/provider/local

> Body 请求参数

```yaml
file:
  - file://C:\Users\19822\Desktop\M718.yaml
  - file://C:\Users\19822\Desktop\M728.yaml

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |none|
|body|body|object| 否 |none|
|» file|body|string(binary)| 否 |none|

> 返回示例

> 200 Response

```json
{
  "message": [
    {
      "message": "添加机场\"M78\"失败: [ent: constraint failed: UNIQUE constraint failed: providers.name, providers.path]",
      "status": false
    },
    {
      "message": "添加机场\"自定义1\"成功",
      "status": true
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST 添加规则集-本地

POST /api/configuration/add/ruleset/local

> Body 请求参数

```yaml
file:
  - file://C:\Users\19822\Downloads\Bing.srs
  - file://C:\Users\19822\Downloads\geoip-cn.srs
  - file://C:\Users\19822\Downloads\Bing@cn.json

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |none|
|body|body|object| 否 |none|
|» file|body|string(binary)| 否 |none|

> 返回示例

> 200 Response

```json
{
  "message": [
    {
      "message": "添加规则集\"Bing\"失败: [ent: constraint failed: UNIQUE constraint failed: rulesets.name, rulesets.path]",
      "status": false
    },
    {
      "message": "添加规则集\"geoip-cn\"失败: [ent: constraint failed: UNIQUE constraint failed: rulesets.name, rulesets.path]",
      "status": false
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST 添加模板

POST /api/configuration/add/template

> Body 请求参数

```json
{
  "name": "default1",
  "inbounds": [
    {
      "tag": "tun-in",
      "type": "tun",
      "interface_name": "tun0",
      "mtu": 1500,
      "stack": "mixed",
      "auto_route": true,
      "strict_route": true,
      "address": [
        "172.18.0.1/30",
        "fdfe:dcba:9876::1/126"
      ]
    }
  ],
  "outbounds_group": [
    {
      "type": "direct",
      "tag": "direct"
    },
    {
      "type": "selector",
      "tag": "selector",
      "providers": [
        "M78"
      ],
      "tag_groups": [
        "auto"
      ]
    },
    {
      "type": "urltest",
      "tag": "auto",
      "providers": [
        "M78"
      ]
    }
  ],
  "dns": {
    "servers": [
      {
        "tag": "google",
        "type": "tls",
        "server": "8.8.8.8",
        "server_port": 853,
        "detour": "selector"
      },
      {
        "tag": "cloudflare",
        "type": "tls",
        "server": "1.1.1.1",
        "server_port": 853
      }
    ],
    "rules": [
      {
        "rule_set": [
          "Bing",
          "china-ip"
        ],
        "action": "route",
        "server": "cloudflare"
      }
    ],
    "final": "google",
    "strategy": "prefer_ipv4"
  },
  "route": {
    "auto_detect_interface": true,
    "default_domain_resolver": {
      "server": "cloudflare"
    },
    "final": "selector",
    "rule_set": [
      {
        "type": "local",
        "tag": "Bing",
        "path": "/opt/sifubox/uploads/rulesets/112a7e096595d1c32c4ecdfd9e56b66c.srs",
        "format": "binary"
      }
    ],
    "rules": [
      {
        "user": [
          "bind"
        ],
        "action": "route",
        "outbound": "direct"
      },
      {
        "port": [
          53
        ],
        "action": "hijack-dns"
      },
      {
        "protocol": [
          "dns"
        ],
        "action": "hijack-dns"
      },
      {
        "ip_is_private": true,
        "action": "route",
        "outbound": "direct"
      },
      {
        "protocol": [
          "quic"
        ],
        "action": "reject"
      },
      {
        "rule_set": [
          "Bing",
          "china-ip"
        ],
        "action": "route",
        "outbound": "direct"
      },
      {
        "rule_set": [
          "geoip-cn",
          "china-ip"
        ],
        "action": "route",
        "outbound": "direct"
      }
    ]
  },
  "log": {
    "disabled": true
  },
  "experiment": {
    "clash_api": {
      "external_controller": "127.0.0.1:8081",
      "external_ui": "ui",
      "secret": "123456"
    }
  },
  "providers": [
    "M78"
  ]
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |none|
|body|body|object| 否 |none|
|» *anonymous*|body|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "message": "添加模板\"default3\"成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## DELETE 删除模板

DELETE /api/configuration/delete/template

> Body 请求参数

```yaml
name:
  - default
  - default1
  - default2
  - default3
  - default4

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |none|
|body|body|object| 否 |none|
|» name|body|[string]| 是 |名称|

> 返回示例

> 207 Response

```json
[
  {
    "message": "删除模板\"default3\"成功",
    "status": true
  },
  {
    "message": "删除模板\"default4\"成功",
    "status": true
  },
  {
    "message": "查找模板\"default5\"失败: [ent: template not found]",
    "status": false
  }
]
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|207|[Multi-Status](https://tools.ietf.org/html/rfc2518#section-10.2)|none|Inline|

### 返回数据结构

状态码 **207**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» status|boolean|true|none||none|

## PATCH 修改模板

PATCH /api/configuration/edit/template

> Body 请求参数

```json
{
  "name": "default",
  "inbounds": [
    {
      "tag": "tun-in",
      "type": "tun",
      "interface_name": "tun0",
      "mtu": 1500,
      "stack": "mixed",
      "auto_route": true,
      "strict_route": true,
      "address": [
        "172.18.0.1/30",
        "fdfe:dcba:9876::1/126"
      ]
    }
  ],
  "outbounds_group": [
    {
      "type": "direct",
      "tag": "direct"
    },
    {
      "type": "selector",
      "tag": "selector",
      "providers": [
        "M78"
      ],
      "tag_groups": [
        "auto"
      ]
    },
    {
      "type": "urltest",
      "tag": "auto",
      "providers": [
        "M78"
      ]
    }
  ],
  "dns": {
    "servers": [
      {
        "tag": "google",
        "type": "tls",
        "server": "8.8.8.8",
        "server_port": 853,
        "detour": "selector"
      },
      {
        "tag": "cloudflare",
        "type": "tls",
        "server": "1.1.1.1",
        "server_port": 853
      }
    ],
    "rules": [
      {
        "rule_set": [
          "Bing",
          "china-ip"
        ],
        "action": "route",
        "server": "cloudflare"
      }
    ],
    "final": "google",
    "strategy": "prefer_ipv4"
  },
  "route": {
    "auto_detect_interface": true,
    "default_domain_resolver": {
      "server": "cloudflare"
    },
    "final": "selector",
    "rule_set": [
      {
        "type": "local",
        "tag": "Bing",
        "path": "/opt/sifubox/uploads/rulesets/112a7e096595d1c32c4ecdfd9e56b66c.srs",
        "format": "binary"
      }
    ],
    "rules": [
      {
        "user": [
          "bind"
        ],
        "action": "route",
        "outbound": "direct"
      },
      {
        "port": [
          53
        ],
        "action": "hijack-dns"
      },
      {
        "protocol": [
          "dns"
        ],
        "action": "hijack-dns"
      },
      {
        "ip_is_private": true,
        "action": "route",
        "outbound": "direct"
      },
      {
        "protocol": [
          "quic"
        ],
        "action": "reject"
      },
      {
        "rule_set": [
          "Bing",
          "china-ip"
        ],
        "action": "route",
        "outbound": "direct"
      },
      {
        "rule_set": [
          "geoip-cn",
          "china-ip"
        ],
        "action": "route",
        "outbound": "direct"
      }
    ]
  },
  "log": {
    "disabled": true
  },
  "experiment": {
    "clash_api": {
      "external_controller": "0.0.0.0:8080",
      "external_ui": "ui",
      "secret": "123456"
    }
  },
  "providers": [
    "M78"
  ]
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |none|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "message": "修改模板\"default1\"成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## GET 获取yacd面板信息

GET /api/configuration/yacd

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "message": {
    "url": "http://192.168.50.5",
    "secret": "123456"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## GET 获取默认模板

GET /api/configuration/default/template

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

# 迁移

## GET 获取配置yaml文件

GET /api/migrate/export

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 否 |none|

> 返回示例

> 200 Response

```
"providers:\n    - name: M78\n      path: https: //sub.m78sc.cn/api/v1/client/subscribe?token=083387dce0f02a10e8115379f9871c6d\n      remote: true\n      uuid: e77efd6eea4bac28a59fb71ad3816b57\n      nodes:\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v1a02.78787878.top\n          server_port: 31001\n          tag: \"剩余流量：140.95 GB\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v1a02.78787878.top\n          server_port: 31001\n          tag: \"距离下次重置剩余：27 天\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v1a02.78787878.top\n          server_port: 31001\n          tag: \"套餐到期：2025-10-28\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v1a02.78787878.top\n          server_port: 31001\n          tag: \"官网：m78.at\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v1a01.78787878.top\n          server_port: 31001\n          tag: \"电信用户建议使用01或Premium节点\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v1a001.78787878.top\n          server_port: 31001\n          tag: \"\\U0001F1ED\\U0001F1F0香港01 家宽（推荐）\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v1a002.78787878.top\n          server_port: 31002\n          tag: \"\\U0001F1ED\\U0001F1F0香港02 家宽\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v2a003.78787878.top\n          server_port: 31003\n          tag: \"\\U0001F1ED\\U0001F1F0香港03 家宽\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v3s01.78787878.top\n          server_port: 31013\n          tag: \"\\U0001F1ED\\U0001F1F0香港Premium01｜x2\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v1a001.78787878.top\n          server_port: 32001\n          tag: \"\\U0001F1F8\\U0001F1EC新加坡01\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v1a002.78787878.top\n          server_port: 32012\n          tag: \"\\U0001F1F8\\U0001F1EC新加坡02\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v2a003.78787878.top\n          server_port: 32003\n          tag: \"\\U0001F1F8\\U0001F1EC新加坡03\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v3s01.78787878.top\n          server_port: 32013\n          tag: \"\\U0001F1F8\\U0001F1EC新加坡Premium01｜x2\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v1a001.78787878.top\n          server_port: 34001\n          tag: \"\\U0001F1E8\\U0001F1F3台湾01\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v1a002.78787878.top\n          server_port: 34002\n          tag: \"\\U0001F1E8\\U0001F1F3台湾02\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v2a003.78787878.top\n          server_port: 34003\n          tag: \"\\U0001F1E8\\U0001F1F3台湾03\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v1a001.78787878.top\n          server_port: 33001\n          tag: \"\\U0001F1EF\\U0001F1F5日本01\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v1a002.78787878.top\n          server_port: 33002\n          tag: \"\\U0001F1EF\\U0001F1F5日本02\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v2a003.78787878.top\n          server_port: 33003\n          tag: \"\\U0001F1EF\\U0001F1F5日本03\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v3s01.78787878.top\n          server_port: 33003\n          tag: \"\\U0001F1EF\\U0001F1F5日本Premium01｜x2\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v1a001.78787878.top\n          server_port: 35001\n          tag: \"\\U0001F1FA\\U0001F1F8美国01\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v1a002.78787878.top\n          server_port: 35002\n          tag: \"\\U0001F1FA\\U0001F1F8美国02\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: v2a003.78787878.top\n          server_port: 35013\n          tag: \"\\U0001F1FA\\U0001F1F8美国03\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: us-direct01.78787878.top\n          server_port: 45001\n          tag: \"美国直连01(需要IPv6)|x0.3|测试\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: aa02rare.78787878.top\n          server_port: 38463\n          tag: \"\\U0001F1F0\\U0001F1F7韩国\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: aa01rare.78787878.top\n          server_port: 30549\n          tag: \"\\U0001F1E6\\U0001F1F6南极 可改B站IP\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: aa01rare.78787878.top\n          server_port: 30551\n          tag: \"\\U0001F1E6\\U0001F1F7阿根廷\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: aa01rare.78787878.top\n          server_port: 18087\n          tag: \"\\U0001F1EA\\U0001F1EC埃及\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: ss02rare.78787878.top\n          server_port: 31475\n          tag: \"\\U0001F1F9\\U0001F1F7土耳其\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: ss03rare.78787878.top\n          server_port: 25080\n          tag: \"\\U0001F1FA\\U0001F1E6乌克兰\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: ss03rare.78787878.top\n          server_port: 27734\n          tag: \"\\U0001F1EB\\U0001F1F7法国\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: ss03rare.78787878.top\n          server_port: 35919\n          tag: \"\\U0001F1E9\\U0001F1EA德国\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: ss03rare.78787878.top\n          server_port: 38599\n          tag: \"\\U0001F1F2\\U0001F1FE马来西亚\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: ss01rare.78787878.top\n          server_port: 12267\n          tag: \"\\U0001F1F3\\U0001F1F1荷兰\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: ss02rare.78787878.top\n          server_port: 31923\n          tag: \"\\U0001F1F3\\U0001F1EC尼日利亚\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: aa01rare.78787878.top\n          server_port: 29790\n          tag: \"\\U0001F1F9\\U0001F1ED泰国\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: aa02rare.78787878.top\n          server_port: 39044\n          tag: \"\\U0001F1F7\\U0001F1FA俄罗斯\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: aa02rare.78787878.top\n          server_port: 20748\n          tag: \"\\U0001F1EE\\U0001F1E9印度尼西亚\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: aa02rare.78787878.top\n          server_port: 41861\n          tag: \"\\U0001F1FB\\U0001F1F3越南\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: aa02rare.78787878.top\n          server_port: 49188\n          tag: \"\\U0001F1EE\\U0001F1F3印度\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: aa02rare.78787878.top\n          server_port: 34136\n          tag: \"\\U0001F1F5\\U0001F1ED菲律宾｜x2\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: aa02rare.78787878.top\n          server_port: 41708\n          tag: \"\\U0001F1F2\\U0001F1F4澳门\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: aa02rare.78787878.top\n          server_port: 18353\n          tag: \"\\U0001F1E6\\U0001F1FA澳大利亚\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: aa02rare.78787878.top\n          server_port: 24817\n          tag: \"\\U0001F1EC\\U0001F1E7英国\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: aa02rare.78787878.top\n          server_port: 38188\n          tag: \"\\U0001F1E7\\U0001F1F7巴西\\r\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: aa02rare.78787878.top\n          server_port: 35555\n          tag: \"\\U0001F1E6\\U0001F1EA迪拜\\r\"\n          type: shadowsocks\n      templates:\n        - default\n        - default1\n    - name: M718\n      path: /opt/sifubox/uploads/providers/e166c2105f66562d61a2afce952763fb.yaml\n      remote: false\n      uuid: 725f689fc67f1cb521788196d1c8db25\n      nodes:\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: aa02rare.78787878.top\n          server_port: 38188\n          tag: \"\\U0001F1E7\\U0001F1F7巴西111\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: aa02rare.78787878.top\n          server_port: 35555\n          tag: \"\\U0001F1E6\\U0001F1EA迪拜23\"\n          type: shadowsocks\n    - name: M728\n      path: /opt/sifubox/uploads/providers/a09a24f2166195e220194252ac1386a7.yaml\n      remote: false\n      uuid: d7ea95894a555d2bbf1e3e8b0b08da0d\n      nodes:\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: aa02rare.78787878.top\n          server_port: 38188\n          tag: \"\\U0001F1E7\\U0001F1F7巴西1121\"\n          type: shadowsocks\n        - method: chacha20-ietf-poly1305\n          password: 4494bdef-c36f-4778-890a-fc7e09ee6087\n          server: aa02rare.78787878.top\n          server_port: 35555\n          tag: \"\\U0001F1E6\\U0001F1EA迪拜432\"\n          type: shadowsocks\nrulesets:\n    - name: Bing\n      path: /opt/sifubox/uploads/rulesets/112a7e096595d1c32c4ecdfd9e56b66c.srs\n      remote: false\n      binary: true\n      templates:\n        - default1\n        - default\n    - name: geoip-cn\n      path: /opt/sifubox/uploads/rulesets/0b1e892fe64acae39ac4e5106f55474f.srs\n      remote: false\n      binary: true\n      templates:\n        - default1\n        - default\n    - name: china-ip\n      path: https: //github.com/SagerNet/sing-geoip/raw/refs/heads/rule-set/geoip-cn.srs\n      remote: true\n      update_interval: 1d\n      binary: true\n      download_detour: direct\n      templates:\n        - default\n        - default1\ntemplates:\n    - name: default\n      ntp:\n        enabled: false\n        server: \"\"\n      inbounds:\n        - address:\n            - 172.18.0.1/30\n            - fdfe:dcba: 9876: : 1/126\n          auto_route: true\n          interface_name: tun0\n          mtu: 1500\n          stack: mixed\n          strict_route: true\n          tag: tun-in\n          type: tun\n      outbounds_group:\n        - type: direct\n          tag: direct\n          providers: []\n          tag_groups: []\n        - type: selector\n          tag: selector\n          providers:\n            - M78\n          tag_groups:\n            - auto\n        - type: urltest\n          tag: auto\n          providers:\n            - M78\n            - M718\n          tag_groups: []\n      dns:\n        final: google\n        strategy: prefer_ipv4\n        servers:\n            - server: 8.8.8.8\n              server_port: 853\n              tag: google\n              type: tls\n            - server: 1.1.1.1\n              server_port: 853\n              tag: cloudflare\n              type: tls\n      experiment:\n        clash_api:\n            external_controller: 127.0.0.1: 9090\n            external_ui: ui\n            secret: \"123456\"\n      log:\n        disabled: true\n        output: \"\"\n        timestamp: false\n        level: \"\"\n      route:\n        rule_set:\n            - type: remote\n              tag: china-ip\n              format: binary\n              url: https: //github.com/MetaCubeX/meta-rules-dat/raw/bd4354ba7f11a22883b919ac9fb9f7034fb51b31/geo/geoip/cn.srs\n              download_detour: direct\n              update_interval: 1d\n        final: direct\n        default_domain_resolver:\n            server: google\n      providers:\n        - M78\n        - M718\n    - name: default1\n      ntp:\n        enabled: false\n        server: \"\"\n      inbounds:\n        - address:\n            - 172.18.0.1/30\n            - fdfe:dcba: 9876: : 1/126\n          auto_route: true\n          interface_name: tun0\n          mtu: 1500\n          stack: mixed\n          strict_route: true\n          tag: tun-in\n          type: tun\n      outbounds_group:\n        - type: direct\n          tag: direct\n          providers: []\n          tag_groups: []\n        - type: selector\n          tag: selector\n          providers:\n            - M78\n            - M728\n          tag_groups:\n            - auto\n        - type: urltest\n          tag: auto\n          providers:\n            - M78\n            - M728\n          tag_groups: []\n      dns:\n        final: google\n        strategy: prefer_ipv4\n        rules:\n            - action: route\n              outbound: direct\n              rule_set:\n                - Bing\n                - china-ip\n        servers:\n            - server: 8.8.8.8\n              server_port: 853\n              tag: google\n              type: tls\n            - server: 1.1.1.1\n              server_port: 853\n              tag: cloudflare\n              type: tls\n      experiment:\n        clash_api:\n            external_controller: 127.0.0.1: 9090\n            external_ui: ui\n            secret: \"123456\"\n      log:\n        disabled: true\n        output: \"\"\n        timestamp: false\n        level: \"\"\n      route:\n        rules:\n            - action: route\n              outbound: direct\n              user:\n                - bind\n            - action: hijack-dns\n              port:\n                - 53\n            - action: hijack-dns\n              protocol:\n                - dns\n            - action: route\n              ip_is_private: true\n              outbound: direct\n            - action: reject\n              protocol:\n                - quic\n            - action: route\n              outbound: direct\n              rule_set:\n                - Bing\n                - china-ip\n        rule_set:\n            - type: remote\n              tag: china-ip\n              format: binary\n              url: https: //github.com/MetaCubeX/meta-rules-dat/raw/bd4354ba7f11a22883b919ac9fb9f7034fb51b31/geo/geoip/cn.srs\n              download_detour: direct\n              update_interval: 1d\n            - type: local\n              tag: Bing\n              format: binary\n              path: /opt/sifubox/uploads/rulesets/112a7e096595d1c32c4ecdfd9e56b66c.srs\n        final: direct\n        default_domain_resolver:\n            server: google\n      providers:\n        - M78\n        - M718\n        - M728\n"
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST 上传配置文件

POST /api/migrate/import

> Body 请求参数

```yaml
file: file://C:\Users\19822\Downloads\migrate.yaml

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |none|
|body|body|object| 否 |none|
|» file|body|string(binary)| 否 |none|

> 返回示例

> 207 Response

```json
[
  {
    "message": "添加机场\"M78\"失败: [ent: constraint failed: UNIQUE constraint failed: providers.name, providers.path]",
    "status": false
  },
  {
    "message": "添加机场\"M718\"失败: [ent: constraint failed: UNIQUE constraint failed: providers.name, providers.path]",
    "status": false
  },
  {
    "message": "添加机场\"M728\"失败: [ent: constraint failed: UNIQUE constraint failed: providers.name, providers.path]",
    "status": false
  },
  {
    "message": "添加规则集\"Bing\"失败: [ent: constraint failed: UNIQUE constraint failed: rulesets.name, rulesets.path]",
    "status": false
  },
  {
    "message": "添加规则集\"geoip-cn\"失败: [ent: constraint failed: UNIQUE constraint failed: rulesets.name, rulesets.path]",
    "status": false
  },
  {
    "message": "添加规则集\"china-ip\"失败: [ent: constraint failed: UNIQUE constraint failed: rulesets.name, rulesets.path]",
    "status": false
  },
  {
    "message": "添加模板\"default\"失败: [ent: constraint failed: UNIQUE constraint failed: templates.name]",
    "status": false
  },
  {
    "message": "添加模板\"default1\"失败: [ent: constraint failed: UNIQUE constraint failed: templates.name]",
    "status": false
  }
]
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|207|[Multi-Status](https://tools.ietf.org/html/rfc2518#section-10.2)|none|Inline|

### 返回数据结构

状态码 **207**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» status|boolean|true|none||none|

# 托管

## GET 配置文件列表

GET /api/files/list

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "message": [
    {
      "expire_time": "1763389797",
      "name": "default",
      "path": "c21f969b5f03d33d43e04f8f136e7682.json",
      "signature": "ec035154ff9640613de532ab59156b8c02a09056fc4746c364ff6c0a53a6182a"
    },
    {
      "expire_time": "1763389797",
      "name": "default1",
      "path": "c89feeb47bb697c7cb7f46b1e8186fef.json",
      "signature": "f1e3b27a9f2378ddca4a333cc0936575978bccc9226db081dae5683828191fd0"
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|[object]|true|none||none|
|»» expire_time|string|true|none||none|
|»» name|string|true|none||none|
|»» path|string|true|none||none|
|»» signature|string|true|none||none|

## GET 获取配置文件

GET /api/files/download/default1/1760669215/aaf5d0287f9dd2c24da17065a7ebb93e703bd75a7a927cdf3709eaa61d5371a5/c89feeb47bb697c7cb7f46b1e8186fef.json

> 返回示例

> 500 Response

```json
{
  "message": "文件已过期"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|none|Inline|

### 返回数据结构

状态码 **500**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|

# 应用

## POST 设置模板

POST /api/application/template

> Body 请求参数

```yaml
name: default

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 否 |none|
|body|body|object| 否 |none|
|» name|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST 重设更新间隔

POST /api/application/interval

> Body 请求参数

```yaml
interval: ""

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |none|
|body|body|object| 否 |none|
|» interval|body|string| 是 |如果为空则取消定时任务|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

# 执行

## GET 发送指令

GET /api/execute/check

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## GET 刷新文件

GET /api/execute/refresh

> Body 请求参数

```yaml
{}

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |none|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

# 数据模型

