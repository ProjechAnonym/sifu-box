package controller

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	utils "sifu-box/Utils"
	"strings"
)

// Fetch_links 获取模板文件中的链接信息,并根据配置文件中的代理和服务器设置进行处理
// 返回一个映射,其中键是模板名称,值是包含链接标签和路径的数组如果发生错误,则返回错误
func Fetch_links() (map[string][]map[string]string,error){
    // 获取项目目录路径
	project_dir,err := utils.Get_value("project-dir")
    if err != nil{
        // 记录获取项目目录失败的日志
        utils.Logger_caller("get project dir failed",err,1)
        return nil,err
    }

    // 打开静态文件目录
	// 打开目录
	static_dir, err := os.Open(filepath.Join(project_dir.(string),"static"))
	if err != nil {
		// 记录打开静态文件目录失败的日志
		utils.Logger_caller("failed to open static directory", err,1)
		return nil,err
	}
	defer static_dir.Close()

    // 读取目录中的所有文件和子目录
	// 读取目录条目
	dirs, err := static_dir.ReadDir(-1) // -1 表示读取所有条目
	if err != nil {
		// 记录读取目录失败的日志
		utils.Logger_caller("failed to read template directory", err,1)
		return nil,err
	}

    // 获取代理配置
	proxy_config,err := utils.Get_value("Proxy")
	if err != nil {
		// 记录获取代理配置失败的日志
		utils.Logger_caller("failed to get proxy config", err,1)
		return nil,err
	}

    // 获取服务器配置
	server_config,err := utils.Get_value("Server")
	if err != nil {
		// 记录获取服务器配置失败的日志
		utils.Logger_caller("failed to get server config", err,1)
		return nil,err
	}

    // 对服务器配置中的令牌进行MD5加密
	md5_token,err := utils.Encryption_md5(server_config.(utils.Server_config).Token)
	if err !=  nil{
		// 记录加密令牌失败的日志
		utils.Logger_caller("failed to get md5 token", err,1)
		return nil,err
	}

    // 初始化存储链接信息的映射
	template_links := make(map[string][]map[string]string)
	for _,dir := range(dirs){
        // 打开模板文件目录
		template_file_dir,err := os.Open(filepath.Join(project_dir.(string),"static",dir.Name()))
		if err != nil{
			// 记录打开模板文件目录失败的日志
			utils.Logger_caller("failed to open template file directory", err,1)
		}
		defer template_file_dir.Close()

        // 读取模板文件目录中的所有文件和子目录
		template_file_list,err := template_file_dir.ReadDir(-1)
		if err != nil{
			// 记录读取模板文件目录失败的日志
			utils.Logger_caller("failed to read template directory file list", err,1)
		}

        // 初始化存储当前模板链接的数组
		var links []map[string]string
		for _,file := range template_file_list{
            // 跳过子目录
			if file.IsDir(){
				// 记录模板目录包含子目录的日志
				utils.Logger_caller("template directory contains subdirectory", fmt.Errorf("%s is a subdirectory",file.Name()),1)
				continue
			}
            // 遍历代理配置中的链接,匹配文件名
			for _,link := range(proxy_config.(utils.Box_config).Url){
				md5_link,err := utils.Encryption_md5(link.Label)
				if err != nil{
					// 记录加密链接标签失败的日志
					utils.Logger_caller("failed to encrypt link", err,1)
				}
                // 如果文件名的MD5与链接标签的MD5匹配,则处理该链接
				if md5_link == strings.Split(file.Name(), ".")[0]{
					// 构建链接的完整路径
					path,_ := url.JoinPath("api","files",file.Name())
					params := url.Values{}
					params.Add("token",md5_token)
					params.Add("template",dir.Name())
					params.Add("label",link.Label)
					path += "?" + params.Encode()
                    // 将处理后的链接添加到数组中
					links = append(links, map[string]string{"label":link.Label,"path":path})
					break
				}
			}
            // 将当前模板的链接数组添加到映射中
			template_links[dir.Name()] = links
		}
	}
	
	return template_links,nil
}

func Verify_link(token string) (error) {
	// 获取配置文件
	server_config,err := utils.Get_value("Server")
	if err != nil {
		utils.Logger_caller("get server config failed",err,1)
		return err
	}
	md5_token,err := utils.Encryption_md5(server_config.(utils.Server_config).Token)
	if err != nil {
		utils.Logger_caller("encryption md5 failed",err,1)
		return err
	} 
	if token == md5_token {
		return nil
	} else {
		return errors.New("token error")
	}	

}