// 代码生成时间: 2025-09-11 12:27:52
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
)

// BatchRenameTool 结构体，包含重命名相关的信息
type BatchRenameTool struct {
    // 旧文件路径
    OldPath string
    // 新文件路径
    NewPath string
    // 文件前缀
    Prefix string
    // 起始编号
    StartNumber int
    // 文件类型
    FileType string
}

// NewBatchRenameTool 创建BatchRenameTool实例
func NewBatchRenameTool(oldPath, newpath, prefix string, startNumber int, fileType string) *BatchRenameTool {
    return &BatchRenameTool{
        OldPath:   oldPath,
        NewPath:   newpath,
        Prefix:    prefix,
        StartNumber: startNumber,
        FileType: fileType,
    }
}

// RenameFiles 批量重命名文件函数
func (b *BatchRenameTool) RenameFiles() error {
    // 获取旧路径下所有文件
    files, err := os.ReadDir(b.OldPath)
    if err != nil {
        return err
    }

    // 按文件名排序
    fileSlice := make([]*os.DirEntry, len(files))
    for i, file := range files {
        fileSlice[i] = &files[i]
    }
    sort.Slice(fileSlice, func(i, j int) bool {
        return fileSlice[i].Name() < fileSlice[j].Name()
    })

    // 初始化编号
    count := b.StartNumber

    for _, file := range fileSlice {
        if file.Type().IsRegular() && strings.HasSuffix(file.Name(), b.FileType) {
            // 构造新文件名
            newFileName := fmt.Sprintf("%s%d%s", b.Prefix, count, b.FileType)
            oldFile := filepath.Join(b.OldPath, file.Name())
            newFile := filepath.Join(b.NewPath, newFileName)

            // 重命名文件
            if err := os.Rename(oldFile, newFile); err != nil {
                return err
            }
            count++
        }
    }

    return nil
}

func main() {
    // 创建重命名工具实例
    brt := NewBatchRenameTool("./old", "./new", "prefix_", 1, ".txt")

    // 执行重命名操作
    if err := brt.RenameFiles(); err != nil {
        log.Fatalf("Error renaming files: %v", err)
    } else {
        fmt.Println("Files renamed successfully.")
    }
}
