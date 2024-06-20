package controller

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	utils "sifu-box/Utils"
	"sort"
	"strings"
)

// remove_item 移除数组中指定索引的元素
// 参数 array 是一个接口类型,它可以是任何类型的数组或切片
// 参数 index 是一个整数切片,包含了需要移除的元素的索引
// 返回值是一个接口类型,它代表了移除指定索引后的新数组或切片
func Remove_item(array interface{}, index []int) interface{} {
    // 使用反射获取 array 的反射值,以便操作其元素
    slice := reflect.ValueOf(array)
    // 创建一个新的切片,类型与原数组相同,初始长度为 0,容量为原数组长度减去需要移除的索引数量
    new_array := reflect.MakeSlice(slice.Type(), 0, slice.Len()-len(index))
    // 使用 map 记录需要移除的索引,以便高效查找
    index_map := make(map[int]bool)
    for _, value := range(index) {
        index_map[value] = true
    }
    // 遍历原数组,如果当前索引不在移除列表中,则将该元素添加到新数组中
    for i := 0; i < slice.Len(); i++ {
        if !index_map[i] {
            new_array = reflect.Append(new_array, slice.Index(i))
        }
    }
    // 将新数组转换为接口类型并返回
    return new_array.Interface()
}
// Check_array 检查给定的删除索引数组是否有效
// delete_index 是待检查的删除索引数组
// length 是目标数组的长度
// 返回值表示删除索引数组是否有效
func Check_array(delete_index []int, length int) bool {
    if len(delete_index) == 0 {
        return true
    }
    // 使用map统计每个索引出现的次数
    count := make(map[int]int)
    for _, value := range(delete_index) {
        count[value]++
        // 如果某个索引出现超过一次,则数组无效
        if count[value] > 1 {
            return false
        }
    }
    // 对删除索引数组进行排序
    sort.Ints(delete_index)
    // 检查排序后的最后一个索引是否超过了目标数组的长度,或者删除索引的数量是否超过了目标数组的长度
    // 如果超过,则数组无效
    if delete_index[len(delete_index)-1] >= length || len(delete_index) > length {
        return false
    }
    // 如果以上检查都通过,则数组有效
    return true
}
// Delete_config 根据删除索引数组,从原始数组中删除配置项,并清理相关的配置文件
// origin_array: 原始的配置项数组
// delete_index: 需要删除的配置项的索引数组
// project_dir: 项目目录路径
// 返回值: 删除操作可能产生的错误
func Delete_config(origin_array []utils.Box_url, delete_index []int,project_dir string) error{
    // 创建一个新的数组,用于存储需要删除的配置项
    new_array := make([]utils.Box_url,len(delete_index))
    // 根据delete_index从origin_array中复制需要删除的配置项到new_array
    for i,value := range(delete_index) {
        new_array[i] = origin_array[value]
    }
    // 打开项目目录中的template子目录
    template_dir, err := os.Open(filepath.Join(project_dir,"template"))
    if err != nil {
        // 记录打开template目录失败的日志
        utils.Logger_caller("failed to open template directory", err,1)
    }
    defer template_dir.Close()
    // 读取template目录中的所有文件
    entries, err := template_dir.ReadDir(-1) // -1 表示读取所有条目
    if err != nil {
        // 记录读取template目录失败的日志
        utils.Logger_caller("failed to read template directory", err,1)
        return err
    }
    // 遍历new_array中的每个配置项,进行删除操作
    for _,value := range(new_array) {
        // 如果配置项的Remote标志为false,则删除对应的临时配置文件
        if !value.Remote{
            if err := utils.File_delete(value.Path);err != nil {
                // 记录删除静态配置文件失败的日志
                utils.Logger_caller("failed to delete config file",err,1)
                return err
            }
        }
        // 对配置项的Label进行MD5加密
        label,err := utils.Encryption_md5(value.Label)
        if err!=nil {
            // 记录加密失败的日志
            utils.Logger_caller("md5 encryption failed",err,1)
            return err
        }
        // 遍历template目录中的每个文件,删除与当前配置项相关的静态配置文件
        for _,entry := range(entries) {
            template := strings.Split(entry.Name(), ".")[0]
            err = utils.File_delete(filepath.Join(project_dir,"static",template,fmt.Sprintf("%s.json",label)))
            if err != nil {
                // 记录删除静态配置文件失败的日志
                utils.Logger_caller("failed to delete config file",err,1)
                return err
            }
        }
    }
    // 删除操作完成,返回nil表示无错误发生
    return nil
}