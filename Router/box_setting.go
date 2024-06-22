package router

import (
	"net/http"
	"path/filepath"
	controller "sifu-box/Controller"
	middleware "sifu-box/Middleware"
	utils "sifu-box/Utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// add_items 在给定的路由组中添加处理新增项和文件上传的路由
// group: 一个gin.RouterGroup实例,用于组织和注册路由
func add_items(group *gin.RouterGroup) {
    // 创建一个子路由组,专门处理与"添加"相关的路由
    add_router := group.Group("/add")

    // 注册一个处理添加项的POST请求路由
    add_router.POST("/item", func(ctx *gin.Context) {
        // 解析请求中的JSON数据到config变量
        var config utils.Box_config
        if err := ctx.ShouldBindJSON(&config); err != nil {
            // 日志记录JSON绑定失败,并返回错误响应
            utils.Logger_caller("Marshal json failed!", err, 1)
            ctx.JSON(http.StatusBadRequest, gin.H{"error": "Add items failed."})
            return
        }
        // 调用控制器方法添加项,处理业务逻辑
        if err := controller.Add_items(config); err != nil {
            // 日志记录添加项失败,并返回错误响应
            utils.Logger_caller("Add items failed!", err, 1)
            ctx.JSON(http.StatusBadRequest, gin.H{"error": "Add items failed."})
            return
        }
        // 如果添加成功,返回成功的响应
        ctx.JSON(http.StatusOK, gin.H{"result": "success"})
    })

    // 注册一个处理添加文件的POST请求路由
    add_router.POST("/files", func(ctx *gin.Context) {
        // 解析上传的多部分表单
        form, err := ctx.MultipartForm()
        if err != nil {
            // 日志记录获取多部分表单失败,并返回错误响应
            utils.Logger_caller("get json files failed!", err, 1)
            ctx.JSON(http.StatusBadRequest, gin.H{"error": "Add files failed."})
            return
        }
        // 获取上传的文件列表
        files := form.File["files"]
        // 获取项目目录路径
        project_dir, err := utils.Get_value("project-dir")
        if err != nil {
            // 日志记录获取项目目录失败,并返回内部服务器错误响应
            utils.Logger_caller("get project dir failed!", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "get project dir failed."})
            return
        }
        // 初始化配置结构体和URL列表
        config := utils.Box_config{}
        urls := make([]utils.Box_url, len(files))
        // 检查temp目录是否存在,不存在则创建
        if err := utils.Dir_Create(filepath.Join(project_dir.(string),"temp"),0755); err != nil{
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "creat temp dir failed."})
            return
        }
        // 遍历上传的文件,处理并保存每个文件
        for i, file := range files {
            // 解析文件名,用于生成标签
            name_slice := strings.Split(file.Filename, ".")
            var label string
            if len(name_slice) <= 2 {
                label = name_slice[0]
            } else {
                label = strings.Join(name_slice[0:len(name_slice)-2], "")
            }
            // 构建文件保存路径,并初始化URL结构体
            urls[i] = utils.Box_url{Path: filepath.Join(project_dir.(string), "temp", file.Filename), Proxy: false, Label: label, Remote: false}
            // 保存上传的文件到指定路径
            if err := ctx.SaveUploadedFile(file, filepath.Join(project_dir.(string), "temp", file.Filename)); err != nil {
                // 如果保存文件失败,返回内部服务器错误响应
                ctx.JSON(http.StatusInternalServerError, gin.H{"error": "save file failed."})
                return
            }
        }
        // 将处理后的URL列表赋值给配置结构体
        config.Url = urls
        // 调用控制器方法添加配置,处理业务逻辑
        if err := controller.Add_items(config); err != nil {
            // 日志记录添加失败,并返回错误响应
            utils.Logger_caller("Add items failed!", err, 1)
            ctx.JSON(http.StatusBadRequest, gin.H{"error": "Add items failed."})
            return
        }
        // 如果添加成功,返回成功的响应
        ctx.JSON(http.StatusOK, gin.H{"result": "success"})
    })
}
// fetch_items 设置与获取物品相关的路由
// 使用 gin.RouterGroup 作为参数,允许在已有的路由组中嵌套创建新的路由组,专门处理 /fetch 路径下的请求
func fetch_items(group *gin.RouterGroup) {
    // 在 /fetch 路径下创建一个新的路由组
    fetch_router := group.Group("/fetch")
    
    // 定义 GET 请求的路由,用于获取物品信息
    fetch_router.GET("/items", func(ctx *gin.Context) {
        // 调用 controller 中的 Fetch_items 函数尝试获取物品信息
        config, err := controller.Fetch_items()
        // 如果获取失败,记录错误日志并返回内部服务器错误的响应
        if err != nil {
            utils.Logger_caller("fetch items failed!", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "fetch items failed."})
            return
        }
        // 如果获取成功,返回物品信息
        ctx.JSON(http.StatusOK, config)
    })
}
func delete_items(group *gin.RouterGroup) {
    delete_router := group.Group("/delete")
    delete_router.DELETE("/items", func(ctx *gin.Context) {
        type delete_config struct {
            Urls []int `json:"urls"`
            Rulesets []int `json:"rulesets"`
        }
        var items delete_config
        if err := ctx.BindJSON(&items);err!=nil{
            utils.Logger_caller("marshal json failed!", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "marshal json failed."})
            return
        }
        if err := controller.Delete_items(map[string][]int{"urls":items.Urls,"rulesets":items.Rulesets}); err != nil{
            utils.Logger_caller("delete items failed!", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "delete items failed."})
            return
        }
        ctx.JSON(http.StatusOK, gin.H{"result": "success"})
    })
}
func Setting_box(group *gin.RouterGroup) {
	setting_router := group.Group("/config")
	setting_router.Use(middleware.Token_auth())
	add_items(setting_router)
	fetch_items(setting_router)
    delete_items(setting_router)
}