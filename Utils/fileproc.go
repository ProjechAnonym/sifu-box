package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)
func Encryption_md5(str string) (string,error) {
	h := md5.New()
	_,err := h.Write([]byte(str))
	if err != nil {
		return "",err
	}
	return hex.EncodeToString(h.Sum(nil)),nil
}

func File_delete(dst string) error{
	// 检查目标文件是否存在,若存在则删除
	info, err := os.Stat(dst)
	if err != nil {
        if os.IsNotExist(err) {
            // 文件不存在,不需要删除,直接返回
            return nil
        } else {
            // 其他错误,例如权限问题,返回错误
            Logger_caller("Failed to check file status!", err, 1)
            return err
        }
    }
    // 如果是目录,返回错误
    if info.IsDir() {
        err := fmt.Errorf("cannot remove '%s': it is a directory", dst)
        Logger_caller("Delete file failed!", err, 1)
        return err
    }

    // 尝试删除文件
    if err := os.Remove(dst); err != nil {
        Logger_caller("Delete file failed!", err, 1)
        return err
    }
	return nil
}
// File_write 写入内容到指定文件
// content: 需要写入文件的内容,以字节切片形式提供
// dst: 目标文件的路径
// perm: 文件的权限设置,以文件模式切片形式提供
// 返回值: 如果操作失败,返回错误信息；成功则返回nil
func File_write(content []byte, dst string,perm []fs.FileMode) error {
    // 检查目标文件目录是否存在,若不存在则创建
    if _, err := os.Stat(filepath.Dir(dst)); err != nil {
        if os.IsNotExist(err) {
            os.MkdirAll(filepath.Dir(dst), perm[0])
        }
    }

    // 检查目标文件是否存在,若存在则删除,为新建文件做准备
    if err := File_delete(dst); err != nil{
		Logger_caller("Delete file failed!", err,1)
		return err
	}


    // 打开(若不存在则创建)文件,准备进行写操作
    file, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR, perm[1])
    defer func() {
        // 确保文件在函数返回前关闭,避免资源泄露
        if err := file.Close(); err != nil {
            Logger_caller("File can not close!", err,1)
        }
    }()
    if err != nil {
        Logger_caller("Create file failed", err,1)
        return err
    }

    // 将内容写入文件
    _, err = file.WriteString(string(content))
    if err != nil {
        Logger_caller("Write config failed!", err,1)
        return err
    }

    // 操作成功,返回nil
    return nil
}
// File_copy 实现了源文件到目标文件的复制功能
// src: 源文件路径
// dst: 目标文件路径
// perm: 目标文件的权限模式
// 返回值: 如果复制过程中发生错误,则返回错误；否则返回nil
func File_copy(src, dst string,perm fs.FileMode) error{
    // 打开源文件
    src_file, err := os.Open(src)
    if err != nil {
        // 记录打开源文件失败的日志并返回错误
        Logger_caller("Open file failed!", err,1)
        return err
    }
    // 确保在函数返回前关闭源文件
    defer src_file.Close()
    
    // 创建或打开目标文件,以读写模式,并设置指定的权限
    target_file, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC, perm)
    if err != nil {
        // 记录创建目标文件失败的日志并返回错误
        Logger_caller("Create file failed!", err,1)
        return err
    }
    // 确保在函数返回前关闭目标文件
    defer target_file.Close()
    
    // 使用io.Copy函数将源文件的内容复制到目标文件
    if _,err = io.Copy(target_file,src_file);err != nil{
        // 记录文件复制失败的日志并返回错误
        Logger_caller("Copy file failed!", err,1)
        return err
    }
    
    // 如果复制成功,返回nil
    return nil
}

func Dir_Create(src string,perm fs.FileMode) error{
    // 检查临时目录是否存在,不存在则创建临时目录
    if _,err := os.Stat(src); err != nil {
        if os.IsNotExist(err){
            if err := os.MkdirAll(src,0755); err != nil{
                Logger_caller(fmt.Sprintf("create %s failed!",src),err,1)
                return err
            }
        }else{
            Logger_caller(fmt.Sprintf("check %s has error!",src),err,1)
            return err
        }
    }
    return nil
}