package utils

import (
	"io"
	"os"
	"path/filepath"
)

// FileWrite 将给定的字节切片内容写入指定路径的文件中
// 如果文件不存在,将创建新文件;如果文件存在,将覆盖原有内容
// 参数:
//   content: 待写入文件的字节切片
//   dst: 文件路径
// 返回值:
//   error: 写入操作中可能遇到的错误
func FileWrite(content []byte, dst string) error {
	
	// 检查文件所在目录是否存在,不存在则创建
	if _, err := os.Stat(filepath.Dir(dst)); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(dst),0755); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// 打开或创建文件,使用读写方式,并在文件存在时清空文件内容
	file, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	// 确保文件在函数返回前被关闭
	defer func() {
		if err := file.Close(); err != nil {
			LoggerCaller("文件无法关闭", err, 1)
		}
	}()
	if err != nil {
		return err
	}

	// 将字节切片转换为字符串并写入文件
	_, err = file.WriteString(string(content))
	if err != nil {
		return err
	}
	
	// 文件写入成功,返回nil表示没有发生错误
	return nil
}

// DirCreate 检查给定路径的文件或文件夹是否存在,如果不存在,则创建该路径
// src: 需要检查和创建的文件或文件夹路径
func DirCreate(src string) error {
    // 检查路径是否存在
    if _, err := os.Stat(src); err != nil {
        // 如果路径不存在
        if os.IsNotExist(err) {
            // 尝试创建路径,使用0755权限确保对所有用户可读可执行
            if err := os.MkdirAll(src, 0755); err != nil {
                return err
            }
        } else {
            // 如果存在其他错误,直接返回错误
            return err
        }
    }
    // 路径已存在或已成功创建,返回nil
    return nil
}

// FileDelete 删除指定路径的文件或目录
// 参数 dst 是待删除的文件或目录的路径
// 返回值为错误类型,表示删除操作的结果
func FileDelete(dst string) error{
    // 检查目标文件或目录是否存在
    _, err := os.Stat(dst)
    if err != nil {
        // 如果不存在,则直接返回,不进行删除操作
        if os.IsNotExist(err) {
            return nil
        } else {
            // 如果存在其他错误,则返回错误信息
            return err
        }
    }

    // 删除文件或目录,如果失败则返回错误信息
    if err := os.RemoveAll(dst); err != nil {
        return err
    }
    // 删除成功,返回nil
    return nil
}

// FileCopy 文件复制函数,将源文件复制到目标文件
// 参数 src 为源文件路径,dst 为目标文件路径
// 返回 error,表示在复制过程中可能发生的错误
func FileCopy(src, dst string) error {
    // 打开源文件以读取内容
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    // 确保在函数返回前关闭文件
    defer srcFile.Close()
    
    // 打开或创建目标文件,如果文件已存在,则会进行截断
    targetFile, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
    if err != nil {
        return err
    }
    // 确保在函数返回前关闭文件
    defer targetFile.Close()
    
    // 使用 io.Copy 进行文件内容的复制
    if _, err = io.Copy(targetFile, srcFile); err != nil {
        return err
    }
    
    // 文件复制成功,返回 nil 表示没有发生错误
    return nil
}

// FileRead 读取文件内容
// 参数 src 是待读取文件的路径
// 返回值是文件内容的字节切片和可能的错误
func FileRead(src string) ([]byte, error) {
    // 打开文件,srcFile 为文件的抽象表示
    srcFile, err := os.Open(src)
    // 如果打开文件时发生错误,直接返回错误
    if err != nil {
        return nil, err
    }
    
    // 确保在函数返回前关闭文件
    defer srcFile.Close()
    
    // 读取文件全部内容到字节切片中
    content, err := io.ReadAll(srcFile)
    // 如果读取过程中发生错误,返回错误
    if err != nil {
        return nil, err
    }
    
    // 返回文件内容和无错误
    return content, nil
}