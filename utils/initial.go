package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sifu-box/models"
	"strings"

	"github.com/glebarez/sqlite"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var DiskDb *gorm.DB
var MemoryDb *gorm.DB
func GetDatabase() error{
	project_dir, err := GetValue("project-dir"); 
	if err != nil {
		return err
	}
	if DiskDb, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s/sifu-box.db", project_dir)), &gorm.Config{}); err != nil{
		return err
	}
	if MemoryDb,err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{}); err != nil{
		return err
	}
	DiskDb.AutoMigrate(&models.Host{})
	MemoryDb.AutoMigrate(&models.Provider{},&models.Ruleset{})
	return nil
}

// LoadConfig 加载配置文件并根据配置类型执行相应的操作
// 参数dst表示配置文件的路径,参数class表示配置的类型（如"mode"或"proxy"）
// 返回error类型,表示加载过程中可能出现的错误
func LoadConfig(dst string, class string) error {
    // 获取项目目录值
    projectDir, err := GetValue("project-dir"); 
    if err != nil {
        // 如果获取失败,返回错误
        return err
    }
    // 设置配置文件路径
    viper.SetConfigFile(filepath.Join(projectDir.(string), dst))
    // 读取配置文件,如果失败返回错误
    if err = viper.ReadInConfig(); err != nil {
        return err
    }
    // 根据配置类型执行相应操作
    switch class {
    case "mode":
        // 如果是"mode"类型,加载服务器配置
        var server models.Server
        if err = viper.Unmarshal(&server); err != nil {
            return err
        }
        // 设置配置值
        SetValue(server, class)
    case "proxy":
        // 如果是"proxy"类型,加载代理配置
        var proxy models.Proxy
        if err = viper.Unmarshal(&proxy); err != nil {
            return err
        }
        // 如果代理配置中有提供者,则创建它们
        if len(proxy.Providers) != 0 {
            if err = MemoryDb.Create(&proxy.Providers).Error; err != nil {
                return err
            }
        }
        // 如果代理配置中有规则集,则创建它们
        if len(proxy.Rulesets) != 0 {
            if err = MemoryDb.Create(&proxy.Rulesets).Error; err != nil {
                return err
            }
        }
    default:
        // 如果配置类型不正确,返回错误
        return fmt.Errorf("类型'%s'不正确", class)
    }
    // 如果一切正常,返回nil
    return nil
}
// LoadTemplate 加载项目中的所有模板文件
// 该函数首先获取项目目录路径,然后打开模板目录,
// 读取目录中的所有文件,并尝试将每个文件解析为模板对象
// 最后,将这些模板对象存储在一个映射中,并将该映射保存到配置中
// 如果在任何步骤中发生错误,函数将返回该错误
func LoadTemplate() error{
    // 获取项目目录路径
    projectDir, err := GetValue("project-dir"); 
    if err != nil {
        return err
    }
    // 打开项目模板目录
    templateDir, err := os.Open(filepath.Join(projectDir.(string),"template"))
    if err != nil {
        return err
    }
    defer templateDir.Close()
    // 读取模板目录中的所有文件
    files, err := templateDir.ReadDir(-1) // -1 表示读取所有条目
    if err != nil {
        return err
    }
    // 初始化模板映射,用于存储文件名和模板对象的映射
    templateMap := make(map[string]models.Template)
    for _, file := range files{
        var template models.Template
        // 获取文件名（不带扩展名）
        fileName := strings.Split(file.Name(), ".")[0]
        // 设置viper配置文件路径
        viper.SetConfigFile(filepath.Join(projectDir.(string),"template",file.Name()))
        // 读取配置文件
        if err = viper.ReadInConfig();err != nil {
            return err
        }
        // 将配置解析为模板对象
        if err = viper.Unmarshal(&template); err != nil {
            return err
        }
        // 将模板对象添加到映射中
        templateMap[fileName] = template
    }
    // 将模板映射保存到配置中
    if err := SetValue(templateMap,"templates"); err != nil {
        return err
    }
    return nil
}
func GetProjectDir() string {
	// base_dir := filepath.Dir(os.Args[0])
	base_dir := "E:/Myproject/sifubox"
	// base_dir := "/root/sifu-clash"
	return base_dir
}
