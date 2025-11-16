package model

import "time"

const EXPIRE_TIME = 3 * 24 * time.Hour      // 登录token有效期
const LINK_VALID_TIME = 24 * 30 * time.Hour // 托管文件的token有效期
const RULE_SET = "rule_set"
