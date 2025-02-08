# 登录

## POST 登录接口

POST /api/login/admin

根据路径不同, 区分管理员登录还是访客登录
* /api/login/admin
管理员登录, 需要用户名和密码
* /api/login/visitor
访客登录, 需要管理员提供的密钥, 在配置文件中设置

> Body 请求参数

```yaml
username: sifulin
password: wsr19990902
code: "123456"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» username|body|string| 否 |当管理员账户登录时需要|
|» password|body|string| 否 |当管理员账户登录时需要|
|» code|body|string| 否 |当访客登录时需要, 由配置文件设置|

> 返回示例

```json
{
    "message": {
        "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwidXVpZCI6Ijk1ZjkzMDJlLTI4M2UtNGU5Yi1iZjBhLTZkZGFkZDMxM2JhZCIsImV4cCI6MTczOTExMzg5Nn0.AwSyoG_DEXcPf-ZxgKdJq962H-cPXg4YDJ3Qis7uvQI",
        "admin": true
    }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

| 名称      | 类型    | 必选 | 约束 | 中文名 | 说明                            |
| --------- | ------- | ---- | ---- | ------ | ------------------------------- |
| » message | object  | true | none | 消息   | 验证成功将返回更新时间的JWT密钥 |
| »» jwt    | string  | true | none |        | none                            |
| »» admin  | boolean | true | none |        | none                            |

## GET 自动登录接口

GET /api/verify

验证保存在浏览器的JWT密钥, 有效则保留原本信息并更新过时时间, 无效则登录失败

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用于验证身份的JWT令牌|

> 返回示例

```json
{
  "message": {
    "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwidXVpZCI6Ijk1ZjkzMDJlLTI4M2UtNGU5Yi1iZjBhLTZkZGFkZDMxM2JhZCIsImV4cCI6MTczOTExMzg5Nn0.AwSyoG_DEXcPf-ZxgKdJq962H-cPXg4YDJ3Qis7uvQI",
    "admin": true
  }
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
|» message|object|true|none|消息|验证成功将返回更新时间的JWT密钥|
|»» jwt|string|true|none||none|
|»» admin|boolean|true|none||none|

# 配置

## GET 获取配置

GET /api/configuration/fetch

该接口识别必须设置header的authorization字段,并通过解析该字段判断是否为管理员,只有管理员才有权限获取配置信息

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |浏览器存储的JWT密钥|

> 返回示例

> 该接口会返回所有的配置信息, 包括机场、规则集以及模板

```json
{
  "message": {
    "providers": [
      {
        "name": "夜煞云",
        "path": "https://45.137.180.216/api/v1/client/nyth?token=aa751df806caf70136e33c9310531fe6",
        "remote": true,
        "detour": "select"
      },
      {
        "name": "github1",
        "path": "C:/Users/19822/Downloads/20250119.yaml",
        "remote": false,
        "detour": "select"
      }
    ],
    "rulesets": [
      {
        "tag": "china-site",
        "type": "remote",
        "path": "https://github.com/MetaCubeX/meta-rules-dat/raw/sing/geo/geosite/cn.srs",
        "format": "binary",
        "china": true,
        "name_server": "internal",
        "label": "china",
        "download_detour": "select",
        "update_interval": "1d"
      },
      {
        "tag": "china-ip",
        "type": "remote",
        "path": "https://github.com/MetaCubeX/meta-rules-dat/raw/sing/geo/geoip/cn.srs",
        "format": "binary",
        "china": true,
        "label": "china",
        "download_detour": "select",
        "update_interval": "1d"
      }
    ],
    "templates": {
      "default": {
        "log": {
          "output": "singbox.log",
          "disabled": false,
          "timestamp": true
        },
        "experimental": {
          "cache_file": {
            "enabled": true
          },
          "clash_api": {
            "external_controller": "0.0.0.0:9090",
            "external_ui": "ui",
            "external_ui_download_url": "https://github.com/MetaCubeX/Yacd-meta/archive/gh-pages.zip",
            "external_ui_download_detour": "select",
            "Secret": "123456"
          }
        },
        "inbounds": [
          {
            "type": "tun",
            "tag": "tun-in",
            "interface_name": "tun0",
            "address": [
              "172.18.0.1/30",
              "fdfe:dcba:9876::1/126"
            ],
            "mtu": 9000,
            "auto_route": true,
            "strict_route": true,
            "stack": "mixed"
          }
        ],
        "dns": {
          "final": "external",
          "strategy": "prefer_ipv4",
          "reverse_mapping": true,
          "servers": [
            {
              "tag": "external",
              "address": "tls://8.8.8.8",
              "detour": "select"
            },
            {
              "tag": "internal",
              "address": "172.23.30.94:5353",
              "detour": "direct"
            },
            {
              "tag": "dns_block",
              "address": "rcode://refused"
            }
          ],
          "rules": [
            {
              "outbound": [
                "any"
              ],
              "action": "route",
              "server": "internal"
            }
          ]
        },
        "route": {
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
            }
          ],
          "final": "select",
          "auto_detect_interface": true
        },
        "outbounds": [
          {
            "tag": "direct",
            "type": "direct"
          }
        ],
        "custom_outbounds": [
          {
            "method": "chacha20-ietf-poly1305",
            "password": "123456dsaf",
            "server": "127.0.0.1",
            "server_port": 443,
            "tag": "香港",
            "type": "shadowsocks"
          }
        ]
      },
      "ios": {
        "log": {
          "disabled": true
        },
        "experimental": {
          "cache_file": {
            "enabled": true
          },
          "clash_api": {
            "external_controller": "0.0.0.0:9090",
            "external_ui": "ui",
            "external_ui_download_url": "https://github.com/MetaCubeX/Yacd-meta/archive/gh-pages.zip",
            "external_ui_download_detour": "select",
            "Secret": "123456"
          }
        },
        "inbounds": [
          {
            "type": "tun",
            "tag": "tun-in",
            "interface_name": "tun0",
            "address": [
              "172.18.0.1/30",
              "fdfe:dcba:9876::1/126"
            ],
            "mtu": 9000,
            "auto_route": true,
            "strict_route": true,
            "stack": "mixed"
          }
        ],
        "dns": {
          "final": "external",
          "strategy": "prefer_ipv4",
          "reverse_mapping": true,
          "servers": [
            {
              "tag": "external",
              "address": "tls://8.8.8.8",
              "detour": "select"
            },
            {
              "tag": "internal",
              "address": "223.5.5.5",
              "detour": "direct"
            },
            {
              "tag": "dns_block",
              "address": "rcode://refused"
            }
          ],
          "rules": [
            {
              "outbound": [
                "any"
              ],
              "action": "route",
              "server": "internal"
            }
          ]
        },
        "route": {
          "rules": [
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
            }
          ],
          "final": "select",
          "auto_detect_interface": true
        },
        "outbounds": [
          {
            "tag": "direct",
            "type": "direct"
          },
          {
            "tag": "block",
            "type": "block"
          }
        ]
      }
    }
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|该接口会返回所有的配置信息, 包括机场、规则集以及模板|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|object|true|none||none|
|»» providers|[object]|true|none||none|
|»»» name|string|true|none||none|
|»»» path|string|true|none||none|
|»»» remote|boolean|true|none||none|
|»»» detour|string|true|none||none|
|»» rulesets|[object]|true|none||none|
|»»» tag|string|true|none||none|
|»»» type|string|true|none||none|
|»»» path|string|true|none||none|
|»»» format|string|true|none||none|
|»»» china|boolean|true|none||none|
|»»» name_server|string|false|none||none|
|»»» label|string|true|none||none|
|»»» download_detour|string|false|none||none|
|»»» update_interval|string|false|none||none|
|»» templates|object|true|none||none|
|»»» default|object|true|none||none|
|»»»» log|object|false|none||none|
|»»»»» output|string|true|none||none|
|»»»»» disabled|boolean|true|none||none|
|»»»»» timestamp|boolean|true|none||none|
|»»»» experimental|object|false|none||none|
|»»»»» cache_file|object|true|none||none|
|»»»»»» enabled|boolean|true|none||none|
|»»»»» clash_api|object|true|none||none|
|»»»»»» external_controller|string|true|none||none|
|»»»»»» external_ui|string|true|none||none|
|»»»»»» external_ui_download_url|string|true|none||none|
|»»»»»» external_ui_download_detour|string|true|none||none|
|»»»»»» Secret|string|true|none||none|
|»»»» inbounds|[object]|true|none||none|
|»»»»» type|string|false|none||none|
|»»»»» tag|string|false|none||none|
|»»»»» interface_name|string|false|none||none|
|»»»»» address|[string]|false|none||none|
|»»»»» mtu|integer|false|none||none|
|»»»»» auto_route|boolean|false|none||none|
|»»»»» strict_route|boolean|false|none||none|
|»»»»» stack|string|false|none||none|
|»»»» dns|object|true|none||none|
|»»»»» final|string|true|none||none|
|»»»»» strategy|string|true|none||none|
|»»»»» reverse_mapping|boolean|true|none||none|
|»»»»» servers|[object]|true|none||none|
|»»»»»» tag|string|true|none||none|
|»»»»»» address|string|true|none||none|
|»»»»»» detour|string|true|none||none|
|»»»»» rules|[object]|true|none||none|
|»»»»»» outbound|[string]|false|none||none|
|»»»»»» action|string|false|none||none|
|»»»»»» server|string|false|none||none|
|»»»» route|object|true|none||none|
|»»»»» rules|[object]|true|none||none|
|»»»»»» user|[string]|false|none||none|
|»»»»»» action|string|true|none||none|
|»»»»»» outbound|string|true|none||none|
|»»»»»» port|[integer]|false|none||none|
|»»»»»» protocol|[string]|true|none||none|
|»»»»»» ip_is_private|boolean|false|none||none|
|»»»»» final|string|true|none||none|
|»»»»» auto_detect_interface|boolean|true|none||none|
|»»»» outbounds|[object]|true|none||none|
|»»»»» tag|string|false|none||none|
|»»»»» type|string|false|none||none|
|»»»» custom_outbounds|[object]|false|none||none|
|»»»»» method|string|false|none||none|
|»»»»» password|string|false|none||none|
|»»»»» server|string|false|none||none|
|»»»»» server_port|integer|false|none||none|
|»»»»» tag|string|false|none||none|
|»»»»» type|string|false|none||none|

## GET 获取默认模板

GET /api/configuration/recover

该接口返回默认的模板配置

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |浏览器存储的JWT密钥|

> 返回示例

> 该接口会返回所有的配置信息, 包括机场、规则集以及模板

```json
{
  "message": {
    "log": {
      "disabled": true
    },
    "experimental": {
      "cache_file": {
        "enabled": true
      },
      "clash_api": {
        "external_controller": "0.0.0.0:9090",
        "external_ui": "ui",
        "external_ui_download_url": "https://github.com/MetaCubeX/Yacd-meta/archive/gh-pages.zip",
        "external_ui_download_detour": "select",
        "Secret": "123456"
      }
    },
    "inbounds": [
      {
        "type": "tun",
        "tag": "tun-in",
        "interface_name": "tun0",
        "address": [
          "172.18.0.1/30",
          "fdfe:dcba:9876::1/126"
        ],
        "mtu": 9000,
        "auto_route": true,
        "strict_route": true,
        "stack": "mixed"
      }
    ],
    "dns": {
      "final": "external",
      "strategy": "prefer_ipv4",
      "reverse_mapping": true,
      "servers": [
        {
          "tag": "external",
          "address": "tls://8.8.8.8",
          "detour": "select"
        },
        {
          "tag": "internal",
          "address": "223.5.5.5",
          "detour": "direct"
        },
        {
          "tag": "dns_block",
          "address": "rcode://refused"
        }
      ],
      "rules": [
        {
          "outbound": [
            "any"
          ],
          "action": "route",
          "server": "internal"
        }
      ]
    },
    "route": {
      "rules": [
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
        }
      ],
      "final": "select",
      "auto_detect_interface": true
    },
    "outbounds": [
      {
        "tag": "direct",
        "type": "direct"
      },
      {
        "tag": "block",
        "type": "block"
      }
    ]
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|该接口会返回所有的配置信息, 包括机场、规则集以及模板|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|object|true|none||none|
|»» log|object|true|none||none|
|»»» disabled|boolean|true|none||none|
|»» experimental|object|true|none||none|
|»»» cache_file|object|true|none||none|
|»»»» enabled|boolean|true|none||none|
|»»» clash_api|object|true|none||none|
|»»»» external_controller|string|true|none||none|
|»»»» external_ui|string|true|none||none|
|»»»» external_ui_download_url|string|true|none||none|
|»»»» external_ui_download_detour|string|true|none||none|
|»»»» Secret|string|true|none||none|
|»» inbounds|[object]|true|none||none|
|»»» type|string|false|none||none|
|»»» tag|string|false|none||none|
|»»» interface_name|string|false|none||none|
|»»» address|[string]|false|none||none|
|»»» mtu|integer|false|none||none|
|»»» auto_route|boolean|false|none||none|
|»»» strict_route|boolean|false|none||none|
|»»» stack|string|false|none||none|
|»» dns|object|true|none||none|
|»»» final|string|true|none||none|
|»»» strategy|string|true|none||none|
|»»» reverse_mapping|boolean|true|none||none|
|»»» servers|[object]|true|none||none|
|»»»» tag|string|true|none||none|
|»»»» address|string|true|none||none|
|»»»» detour|string|true|none||none|
|»»» rules|[object]|true|none||none|
|»»»» outbound|[string]|false|none||none|
|»»»» action|string|false|none||none|
|»»»» server|string|false|none||none|
|»» route|object|true|none||none|
|»»» rules|[object]|true|none||none|
|»»»» port|[integer]|false|none||none|
|»»»» action|string|true|none||none|
|»»»» protocol|[string]|true|none||none|
|»»»» ip_is_private|boolean|false|none||none|
|»»»» outbound|string|false|none||none|
|»»» final|string|true|none||none|
|»»» auto_detect_interface|boolean|true|none||none|
|»» outbounds|[object]|true|none||none|
|»»» tag|string|true|none||none|
|»»» type|string|true|none||none|

## POST 添加配置

POST /api/configuration/add

该接口接受json数据,其中rulesets键会影响是否全局更新配置文件,如果没有新添加的规则集则不应该存在这个字段

> Body 请求参数

```json
"{\r\n    \"providers\": [\r\n        {\r\n            \"name\": \"夜煞云5\",\r\n            \"path\": \"https://45.137.180.216/api/v1/client/nyth?token=aa751df806caf70136e33c9310531fe6\",\r\n            \"detour\": \"select\",\r\n            \"remote\": true\r\n        }\r\n    ],\r\n    // \"rulesets\": [\r\n    //     {\r\n    //         \"type\": \"remote\",\r\n    //         \"path\": \"https://github.com/MetaCubeX/meta-rules-dat/raw/sing/geo/geosite/cn.srs1\",\r\n    //         \"format\": \"binary\",\r\n    //         \"tag\": \"china-site1\",\r\n    //         \"download_detour\": \"select\",\r\n    //         \"update_interval\": \"1d\",\r\n    //         \"label\": \"china\",\r\n    //         \"china\": true,\r\n    //         \"name_server\": \"internal\"\r\n    //     },\r\n    //     {\r\n    //         \"type\": \"remote\",\r\n    //         \"path\": \"https://github.com/MetaCubeX/meta-rules-dat/raw/sing/geo/geoip/cn.srs1\",\r\n    //         \"format\": \"binary\",\r\n    //         \"tag\": \"china-ip1\",\r\n    //         \"download_detour\": \"select\",\r\n    //         \"update_interval\": \"1d\",\r\n    //         \"label\": \"china\",\r\n    //         \"china\": true\r\n    //     }\r\n    // ]\r\n}"
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |浏览器存储的JWT密钥|
|body|body|object| 否 |none|
|» providers|body|[object]| 是 |none|
|»» name|body|string| 否 |none|
|»» path|body|string| 否 |none|
|»» detour|body|string| 否 |none|
|»» remote|body|boolean| 否 |none|
|» rulesets|body|[object]| 否 |如没有不应该存在该字段|
|»» type|body|string| 是 |none|
|»» path|body|string| 是 |none|
|»» format|body|string| 是 |none|
|»» tag|body|string| 是 |none|
|»» download_detour|body|string| 是 |none|
|»» update_interval|body|string| 是 |none|
|»» label|body|string| 是 |none|
|»» china|body|boolean| 是 |none|
|»» name_server|body|string| 否 |none|

> 返回示例

```json
{
  "message": "success"
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
|» message|string|true|none||none|

## POST 添加配置文件

POST /api/configuration/files

该接口接受上传的yaml文件,如果没有文件上传则不应使用该接口

> Body 请求参数

```yaml
files:
  - file://C:\Users\19822\Downloads\20250119.yaml
  - file://C:\Users\19822\Downloads\20250119-2.yaml

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |浏览器存储的JWT密钥|
|body|body|object| 否 |none|
|» files|body|string(binary)| 否 |none|

> 返回示例

```json
{
  "message": "success"
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
|» message|string|true|none||none|

## DELETE 删除配置

DELETE /api/configuration/items

该接口的表单字段应保证元素存在, 没有相关的元素则不应设置该字段的表单

> Body 请求参数

```yaml
providers:
  - "20250119"
  - 20250119-2
rulesets: ""
templates: ""

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |浏览器存储的JWT密钥|
|body|body|object| 否 |none|
|» providers|body|[string]| 否 |none|
|» rulesets|body|[string]| 否 |none|
|» templates|body|[string]| 否 |none|

> 返回示例

> 该接口会返回所有的配置信息, 包括机场、规则集以及模板

```json
{
  "message": [
    "打开'github1'文件失败"
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|该接口会返回所有的配置信息, 包括机场、规则集以及模板|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|[string]|true|none||none|

## POST 设置模板

POST /api/configuration/template

该接口通过路径参数设置模板的名称,发送json序列设置template的内容

> Body 请求参数

```json
{
  "log": {
    "output": "singbox.log",
    "disabled": false,
    "timestamp": true
  },
  "experimental": {
    "cache_file": {
      "enabled": true
    },
    "clash_api": {
      "external_controller": "0.0.0.0:9090",
      "external_ui": "ui",
      "external_ui_download_url": "https://github.com/MetaCubeX/Yacd-meta/archive/gh-pages.zip",
      "external_ui_download_detour": "select",
      "Secret": "123456"
    }
  },
  "inbounds": [
    {
      "type": "tun",
      "tag": "tun-in",
      "interface_name": "tun0",
      "address": [
        "172.18.0.1/30",
        "fdfe:dcba:9876::1/126"
      ],
      "mtu": 9000,
      "auto_route": true,
      "strict_route": true,
      "stack": "mixed"
    }
  ],
  "dns": {
    "final": "external",
    "strategy": "prefer_ipv4",
    "reverse_mapping": true,
    "servers": [
      {
        "tag": "external",
        "address": "tls://8.8.8.8",
        "detour": "select"
      },
      {
        "tag": "internal",
        "address": "192.168.1.2:5353",
        "detour": "direct"
      },
      {
        "tag": "dns_block",
        "address": "rcode://refused"
      }
    ],
    "rules": [
      {
        "outbound": [
          "any"
        ],
        "action": "route",
        "server": "internal"
      }
    ]
  },
  "route": {
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
      }
    ],
    "final": "select",
    "auto_detect_interface": true
  },
  "outbounds": [
    {
      "tag": "direct",
      "type": "direct"
    }
  ],
  "custom_outbounds": [
    {
      "method": "chacha20-ietf-poly1305",
      "password": "123456dsaf",
      "server": "127.0.0.2",
      "server_port": 443,
      "tag": "香港",
      "type": "shadowsocks"
    }
  ]
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|name|query|string| 否 |名称|
|Authorization|header|string| 是 |浏览器存储的JWT密钥|
|body|body|object| 否 |none|

> 返回示例

```json
{
  "message": "success"
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
|» message|string|true|none||none|

# 文件

## GET 获取文件列表

GET /api/files/fetch

该接口返回目前生成的配置文件, 其内容包含模板,机场名称以及token

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |认证用的token, 否则谁都能获取所有文件列表|

> 返回示例

```json
{
  "message": {
    "1": [
      {
        "label": "20250119",
        "path": "api/file?path=4abf1b7994dcaaee427c557d85985fad.json&template=1&token=f7ba36460a4b059e6f52297497456949"
      },
      {
        "label": "夜煞云",
        "path": "api/file?path=78d40fb467cdc48a90486b2a7e74f7d6.json&template=1&token=f7ba36460a4b059e6f52297497456949"
      },
      {
        "label": "20250119-2",
        "path": "api/file?path=8557823b7ce5d63bfaa60359e93d2452.json&template=1&token=f7ba36460a4b059e6f52297497456949"
      }
    ],
    "default": [
      {
        "label": "20250119",
        "path": "api/file?path=4abf1b7994dcaaee427c557d85985fad.json&template=default&token=f7ba36460a4b059e6f52297497456949"
      },
      {
        "label": "夜煞云",
        "path": "api/file?path=78d40fb467cdc48a90486b2a7e74f7d6.json&template=default&token=f7ba36460a4b059e6f52297497456949"
      },
      {
        "label": "20250119-2",
        "path": "api/file?path=8557823b7ce5d63bfaa60359e93d2452.json&template=default&token=f7ba36460a4b059e6f52297497456949"
      }
    ],
    "ios": [
      {
        "label": "20250119",
        "path": "api/file?path=4abf1b7994dcaaee427c557d85985fad.json&template=ios&token=f7ba36460a4b059e6f52297497456949"
      },
      {
        "label": "夜煞云",
        "path": "api/file?path=78d40fb467cdc48a90486b2a7e74f7d6.json&template=ios&token=f7ba36460a4b059e6f52297497456949"
      },
      {
        "label": "20250119-2",
        "path": "api/file?path=8557823b7ce5d63bfaa60359e93d2452.json&template=ios&token=f7ba36460a4b059e6f52297497456949"
      }
    ]
  }
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
|» message|object|true|none||none|
|»» 1|[object]|true|none||none|
|»»» label|string|true|none||none|
|»»» path|string|true|none||none|
|»» default|[object]|true|none||none|
|»»» label|string|true|none||none|
|»»» path|string|true|none||none|
|»» ios|[object]|true|none||none|
|»»» label|string|true|none||none|
|»»» path|string|true|none||none|

## GET 获取文件列表 Copy

GET /api/file

该接口返回目前生成的配置文件, 其内容包含模板,机场名称以及token

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|template|query|string| 否 |none|
|token|query|string| 否 |none|
|path|query|string| 否 |none|
|Authorization|header|string| 否 |none|

> 返回示例

```json
{
  "message": [
    {
      "label": "夜煞云",
      "path": "api/files/78d40fb467cdc48a90486b2a7e74f7d6.json?label=%E5%A4%9C%E7%85%9E%E4%BA%91&template=1&token=f7ba36460a4b059e6f52297497456949"
    },
    {
      "label": "夜煞云",
      "path": "api/files/78d40fb467cdc48a90486b2a7e74f7d6.json?label=%E5%A4%9C%E7%85%9E%E4%BA%91&template=default&token=f7ba36460a4b059e6f52297497456949"
    },
    {
      "label": "夜煞云",
      "path": "api/files/78d40fb467cdc48a90486b2a7e74f7d6.json?label=%E5%A4%9C%E7%85%9E%E4%BA%91&template=ios&token=f7ba36460a4b059e6f52297497456949"
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
|»» label|string|true|none||none|
|»» path|string|true|none||none|

# 执行

## GET 开启

GET /api/exec/boot

启动singbox

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |浏览器存储的JWT密钥|

> 返回示例

```json
{
  "message": "启动成功"
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
|» message|string|true|none||none|

## GET 刷新

GET /api/exec/refresh

启动singbox

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |浏览器存储的JWT密钥|

> 返回示例

```json
{
  "message": "启动成功"
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
|» message|string|true|none||none|

## GET 停止

GET /api/exec/stop

关闭singbox

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |浏览器存储的JWT密钥|

> 返回示例

```json
{
  "message": "success"
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
|» message|string|true|none||none|

## GET 重启

GET /api/exec/restart

重启singbox

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |浏览器存储的JWT密钥|

> 返回示例

```json
{
  "message": "success"
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
|» message|string|true|none||none|

## GET 重载

GET /api/exec/reload

重新载入配置文件

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |浏览器存储的JWT密钥|

> 返回示例

```json
{
  "message": "success"
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
|» message|string|true|none||none|

## GET 查看

GET /api/exec/status

获取singbox运行状态

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |浏览器存储的JWT密钥|

> 返回示例

```json
{
  "message": true
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
|» message|boolean|true|none||none|

# 应用

## POST 设置模板

POST /api/application/set/template

> Body 请求参数

```yaml
value: default

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |浏览器存储的JWT密钥|
|body|body|object| 否 |none|
|» value|body|string| 否 |none|

> 返回示例

```json
{
  "message": "success"
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
|» message|string|true|none||none|

## POST 设置机场

POST /api/application/set/provider

> Body 请求参数

```yaml
value: "20250119"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |浏览器存储的JWT密钥|
|body|body|object| 否 |none|
|» value|body|string| 否 |none|

> 返回示例

```json
{
  "message": "success"
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
|» message|string|true|none||none|

## POST 设置间隔

POST /api/application/interval

该接口发送cron表达式设置更新间隔

> Body 请求参数

```yaml
interval: "*/2 * * * *"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |浏览器存储的JWT密钥|
|body|body|object| 否 |none|
|» interval|body|string| 否 |none|

> 返回示例

```json
{
  "message": "success"
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
|» message|string|true|none||none|

## GET 获取singbox监听端口

GET /api/application/fetch

该接口返回singbox的监听地址和密钥

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |浏览器存储的JWT密钥|

> 返回示例

```json
{
  "message": {
    "current_provider": "夜煞云",
    "current_template": "1",
    "listen": "http://192.168.1.2",
    "secret": "123456"
  }
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
|» message|object|true|none||none|
|»» current_provider|string|true|none||none|
|»» current_template|string|true|none||none|
|»» listen|string|true|none||none|
|»» secret|string|true|none||none|

# 迁移

## GET 导出

GET /api/migrate/export

该接口将返回当前配置的yaml文件

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |浏览器存储的JWT密钥|

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

## POST 导入配置

POST /api/migrate/import

> Body 请求参数

```yaml
file: file://C:\Users\19822\Desktop\conf.yaml

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 否 |none|
|body|body|object| 否 |none|
|» file|body|string(binary)| 否 |none|

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

