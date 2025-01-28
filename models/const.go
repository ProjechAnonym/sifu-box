package models

import "time"

const SIFUBOXSETTINGKEY = "setting:sifubox"
const DEFAULTTEMPLATEKEY = "template:default"
const SINGBOXSETTINGKEY = "setting:singbox"
const CURRENTPROVIDER = "singbox:provider"
const CURRENTTEMPLATE = "singbox:template"
const DEFAULTTEMPLATEPATH = "default.template.yaml"
const STATICDIR = "static"
const TEMPLATEDIR = "template"
const CLASHCONFIGFILE = "clash"
const SIFUBOXSETTINGFILE = "setting.config.yaml"
const TEMPDIR = "temp"
const SINGBOXCONFIGFILEDIR = "config"
const SINGBOXBACKUPCONFIGFILE = "config.json.bak"
const BACKUPDIR = "backup"
const BOOTCOMMAND = "boot_command"
const RELOADCOMMAND = "reload_command"
const STOPCOMMAND = "stop_command"
const RESTARTCOMMAND = "restart_command"
const CHECKCOMMAND = "check_command"

const EXPIRETIME = time.Hour * 24 