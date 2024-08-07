package controller

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sifu-box/models"
	"sifu-box/utils"
	"strings"
)

// FetchLinks 从静态目录中检索所有链接,并按分类整理
// 它返回一个映射,其中每个键是链接的分类名称,每个值是该分类下所有链接的列表
// 每个链接都是一个映射,包含链接的标签和路径
// 如果在获取工作目录、打开静态文件目录、读取分类目录、获取代理配置、读取服务配置或加密过程中出现错误,函数将返回错误
func FetchLinks() (map[string][]map[string]string, error) {
	// 获取项目目录
	project_dir, err := utils.GetValue("project-dir")
	if err != nil {
		// 记录错误并返回
		utils.LoggerCaller("获取工作目录失败", err, 1)
		return nil, fmt.Errorf("获取工作目录失败")
	}

	// 打开项目目录下的static目录
	staticDir, err := os.Open(filepath.Join(project_dir.(string), "static"))
	if err != nil {
		// 记录错误并返回
		utils.LoggerCaller("打开静态文件目录失败", err, 1)
		return nil, fmt.Errorf("打开静态文件目录失败")
	}
	defer staticDir.Close()

	// 读取static目录下的所有文件和目录
	dirs, err := staticDir.ReadDir(-1)
	if err != nil {
		// 记录错误并返回
		utils.LoggerCaller("无法读取分类目录", err, 1)
		return nil, fmt.Errorf("无法读取分类目录")
	}

	// 获取所有代理配置
	var providers []models.Provider
	if err := utils.MemoryDb.Find(&providers).Error; err != nil {
		// 记录错误并返回
		utils.LoggerCaller("无法获得代理配置", err, 1)
		return nil, fmt.Errorf("无法获得代理配置")
	}

	// 获取服务配置
	serverConfig, err := utils.GetValue("mode")
	if err != nil {
		// 记录错误并返回
		utils.LoggerCaller("无法读取服务配置", err, 1)
		return nil, fmt.Errorf("无法读取服务配置")
	}

	// 加密服务配置中的密钥
	md5Token, err := utils.EncryptionMd5(serverConfig.(models.Server).Key)
	if err != nil {
		// 记录错误并返回
		utils.LoggerCaller("md5加密失败", err, 1)
		return nil, fmt.Errorf("文件托管密钥加密失败")
	}

	// 初始化模板链接映射
	templateLinks := make(map[string][]map[string]string)
	for _, dir := range dirs {
		// 打开每个分类目录
		templateFileDir, err := os.Open(filepath.Join(project_dir.(string), "static", dir.Name()))
		if err != nil {
			// 记录错误并继续处理其他目录
			utils.LoggerCaller(fmt.Sprintf("无法打开'%s'目录", dir.Name()), err, 1)
			continue
		}
		defer templateFileDir.Close()

		// 读取分类目录下的所有文件
		templateFileList, err := templateFileDir.ReadDir(-1)
		if err != nil {
			// 记录错误并继续处理其他目录
			utils.LoggerCaller(fmt.Sprintf("无法读取'%s'目录文件", dir.Name()), err, 1)
			continue
		}

		// 初始化链接列表
		var links []map[string]string
		for _, file := range templateFileList {
			// 检查文件是否为目录
			if file.IsDir() {
				// 记录错误并继续处理其他文件
				utils.LoggerCaller("分类文件夹下存在子文件夹", fmt.Errorf("'%s'是个子文件夹", file.Name()), 1)
				continue
			}

			// 匹配文件名与代理配置
			for _, provider := range providers {
				md5Link, err := utils.EncryptionMd5(provider.Name)
				if err != nil {
					// 记录错误并继续处理其他代理配置
					utils.LoggerCaller("无法加密", err, 1)
					continue
				}

				// 如果文件名匹配代理配置
				if md5Link == strings.Split(file.Name(), ".")[0] {
					// 构建链接路径和参数
					path, _ := url.JoinPath("api", "files", file.Name())
					params := url.Values{}
					params.Add("token", md5Token)
					params.Add("template", dir.Name())
					params.Add("label", provider.Name)
					path += "?" + params.Encode()

					// 将链接添加到列表
					links = append(links, map[string]string{"label": provider.Name, "path": path})
					break
				}
			}
		}
		// 将链接列表添加到模板链接映射
		templateLinks[dir.Name()] = links
	}

	// 返回模板链接映射
	return templateLinks, nil
}

// VerifyLink 验证链接有效性
// 参数 token: 传入的token字符串,用于验证链接是否有效
// 返回值: 如果链接有效,返回nil；否则返回错误信息
func VerifyLink(token string) error {
	// 获取服务器配置信息
	serverConfig, err := utils.GetValue("mode")
	if err != nil {
		// 记录获取运行配置失败的错误日志
		utils.LoggerCaller("获取运行配置失败", err, 1)
		// 返回获取运行配置失败的错误信息
		return fmt.Errorf("获取运行配置失败")
	}
	// 加密服务器配置中的密钥
	md5Token, err := utils.EncryptionMd5(serverConfig.(models.Server).Key)
	if err != nil {
		// 记录加密预置密钥失败的错误日志
		utils.LoggerCaller("加密预置密钥失败", err, 1)
		// 返回加密预置密钥失败的错误信息
		return fmt.Errorf("加密预置密钥失败")
	}
	// 比较传入的token与加密后的密钥是否一致
	if token == md5Token {
		// 如果一致,返回nil,表示链接有效
		return nil
	} else {
		// 如果不一致,返回密钥错误的错误信息
		return errors.New("密钥错误")
	}

}