package router

import (
	"net/http"
	controller "sifu-box/Controller"

	middleware "sifu-box/Middleware"
	utils "sifu-box/Utils"

	"github.com/gin-gonic/gin"
)

// add_server 添加服务器路由处理函数
// 该函数使用gin框架的RouterGroup来组织路由，专门处理与添加服务器相关的HTTP请求。
// 参数:
//   group *gin.RouterGroup: 路由组，用于定义一组具有相同前缀的路由。
func add_server(group *gin.RouterGroup) {
    // 创建一个子路由组,专门处理添加服务器的POST请求。
    add_router := group.Group("/add")
    
    // 定义路由：接收POST请求，路径为/add/server
    add_router.POST("/server",func(ctx *gin.Context) {
        // 定义一个结构体变量,用于存储从JSON请求体中解析出的服务器信息
        var content utils.Server
        
        // 尝试将请求体中的JSON数据绑定到content变量上
        if err := ctx.BindJSON(&content); err != nil {
            // 如果绑定失败,记录错误日志并返回内部服务器错误
            utils.Logger_caller("Marshal json failed!", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "marshal failed!"})
            return
        }
		// 判断IP是否指向本机
        is_localhost,err := controller.Is_localhost(content.Url) 
		if err != nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		content.Localhost = is_localhost
        // 尝试将content变量中的服务器信息插入到数据库中
        if err := utils.Db.Create(&content).Error; err != nil {
            // 如果插入失败,记录错误日志并返回内部服务器错误
            utils.Logger_caller("Write to the database failed!", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "write to the database failed!"})
            return
        }
        
        // 如果插入成功,返回状态码200和成功信息
        ctx.JSON(http.StatusOK, gin.H{"message": true})
    })
}

// remove_server 从服务器组中移除指定的服务器
// 使用DELETE请求方法,通过URL参数指定要移除的服务器URL
// 参数:
//   group *gin.RouterGroup - gin路由器组，用于定义路由
func remove_server(group *gin.RouterGroup) {
    // 创建一个子路由组,专门处理移除操作。
    remove_router := group.Group("/remove")
    
    // 定义删除服务器的路由路径
    remove_router.DELETE("/server", func(ctx *gin.Context) {
        // 从POST表单数据中获取待删除服务器的URL
        url := ctx.PostForm("url")
        
        // 使用GORM从数据库中删除URL对应的服务器记录
        // 如果删除操作失败,记录错误并返回内部服务器错误响应
        if err := utils.Db.Where("url = ?", url).Delete(&utils.Server{}).Error; err != nil {
            utils.Logger_caller("Delete from the database failed!", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "delete data from database failed!"})
            return
        }
        
        // 如果删除成功,返回成功的响应
        ctx.JSON(http.StatusOK, gin.H{"message": true})
    })
}

// Setting_server 配置服务器路由组
// 该函数旨在通过gin框架设置与服务器相关的路由规则
// 该子路由组上应用Token认证的中间件,以确保只有经过授权的请求才能访问服务器设置相关的信息
// 调用add_server函数,将具体的路由处理函数注册到这个子路由组上
//
// 参数:
//   group *gin.RouterGroup - 父路由组，用于创建子路由组并继承其配置
func Setting_server(group *gin.RouterGroup){
    // 创建一个名为"setting"的子路由组，用于处理所有与设置相关的请求
    setting_router := group.Group("/setting")
    
    // 在"setting"子路由组上应用Token认证中间件，确保所有请求都需要通过认证
    setting_router.Use(middleware.Token_auth())
    
    // 注册处理服务器添加相关请求的路由处理函数
    add_server(setting_router)
	remove_server(setting_router)
}