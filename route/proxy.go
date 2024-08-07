package route

import (
	"net/http"
	"path/filepath"
	"sifu-box/controller"
	"sifu-box/middleware"
	"sifu-box/models"
	"sifu-box/utils"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// SettingProxy 配置代理相关的API路由
// 参数：
// - group: *gin.RouterGroup,路由组,用于组织相关的HTTP请求处理函数
// - lock: *sync.Mutex,用于同步访问共享资源的互斥锁
func SettingProxy(group *gin.RouterGroup, lock *sync.Mutex) {
    // 创建代理路由的子组
    route := group.Group("/proxy")
    // 使用Token认证中间件
    route.Use(middleware.TokenAuth())

    // 获取代理配置的处理函数
    route.GET("fetch", func(ctx *gin.Context) {
        // 尝试从controller获取代理配置
        config, err := controller.FetchItems()
        if err != nil {
            // 如果获取失败,返回500错误和错误信息
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "获取代理配置失败"})
            return
        }
        // 如果获取成功,返回200状态码和配置数据
        ctx.JSON(http.StatusOK, config)
    })

    // 添加代理配置的处理函数
    route.POST("add", func(ctx *gin.Context) {
        // 解析请求体中的JSON数据到proxy结构体
        var proxy models.Proxy
        if err := ctx.ShouldBindJSON(&proxy); err != nil {
            // 如果解析JSON失败,记录错误并返回400错误和错误信息
            utils.LoggerCaller("序列化json失败", err, 1)
            ctx.JSON(http.StatusBadRequest, gin.H{"message": "序列化json失败"})
            return
        }

        // 尝试添加解析得到的代理配置
        if errs := controller.AddItems(proxy, lock); len(errs) != 0 {
            // 如果添加失败,收集错误信息并返回500错误和错误信息
            var errors []string
            for _, addErr := range errs {
                errors = append(errors, addErr.Error())
            }
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": errors})
            return
        }

        // 如果添加成功,返回200状态码和成功信息
        ctx.JSON(http.StatusOK, gin.H{"message": true})
    })

    // 删除代理配置的处理函数
    route.DELETE("delete", func(ctx *gin.Context) {
        // 解析请求体中的JSON数据到deleteMap
        deleteMap := make(map[string][]int)
        if err := ctx.ShouldBindJSON(&deleteMap); err != nil {
            // 如果解析JSON失败,记录错误并返回400错误和错误信息
            utils.LoggerCaller("序列化json失败", err, 1)
            ctx.JSON(http.StatusBadRequest, gin.H{"message": "序列化json失败"})
            return
        }

        // 尝试删除指定的代理配置
        if errs := controller.DeleteProxy(deleteMap, lock); len(errs) != 0 {
            // 如果删除失败,收集错误信息并返回500错误和错误信息
            var errors []string
            for _, addErr := range errs {
                errors = append(errors, addErr.Error())
            }
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": errors})
            return
        }

        // 如果删除成功,返回200状态码和成功信息
        ctx.JSON(http.StatusOK, gin.H{"message": true})
    })

    // 上传文件以添加代理配置的处理函数
    route.POST("files", func(ctx *gin.Context) {
        // 解析multipart/form-data格式的表单
        form, err := ctx.MultipartForm()
        if err != nil {
            // 如果解析表单失败,记录错误并返回400错误和错误信息
            utils.LoggerCaller("解析表单失败", err, 1)
            ctx.JSON(http.StatusBadRequest, gin.H{"message": "解析表单失败"})
            return
        }

        // 获取文件列表
        files := form.File["files"]

        // 获取项目目录
        projectDir, err := utils.GetValue("project-dir")
        if err != nil {
            // 如果获取项目目录失败,记录错误并返回500错误和错误信息
            utils.LoggerCaller("获取工作目录失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "获取工作目录失败"})
            return
        }

        // 准备Providers列表
        providers := make([]models.Provider, len(files))

        // 确保项目目录下的temp文件夹存在
        if err := utils.DirCreate(filepath.Join(projectDir.(string), "temp")); err != nil {
            utils.LoggerCaller("创建temp文件夹失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "创建temp文件夹失败"})
            return
        }

        // 处理每个上传的文件
        for i, file := range files {
            // 根据文件名生成Provider的Name和Path
            nameSlice := strings.Split(file.Filename, ".")
            var label string
            if len(nameSlice) <= 2 {
                label = nameSlice[0]
            } else {
                label = strings.Join(nameSlice[0:len(nameSlice)-2], "")
            }

            providers[i] = models.Provider{
                Path: filepath.Join(projectDir.(string), "temp", file.Filename),
                Proxy: false,
                Name: label,
                Remote: false,
            }

            // 保存上传的文件到项目目录的temp文件夹
            if err := ctx.SaveUploadedFile(file, filepath.Join(projectDir.(string), "temp", file.Filename)); err != nil {
                ctx.JSON(http.StatusInternalServerError, gin.H{"message": "保存文件失败"})
                return
            }
        }

        // 尝试添加通过文件解析得到的代理配置
        if err := controller.AddItems(models.Proxy{Providers: providers, Rulesets: []models.Ruleset{}}, lock); err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "添加代理配置失败"})
            return
        }

        // 如果添加成功,返回200状态码和成功信息
        ctx.JSON(http.StatusOK, gin.H{"message": true})
    })
}