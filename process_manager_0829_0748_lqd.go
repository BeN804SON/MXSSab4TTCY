// 代码生成时间: 2025-08-29 07:48:26
// process_manager.go
// 此文件定义一个进程管理器，可以启动和停止进程。

package main

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "syscall"
)

// ProcessManager 结构体包含一个进程
type ProcessManager struct {
    process *os.Process
}

// NewProcessManager 创建一个新的进程管理器实例
func NewProcessManager() *ProcessManager {
    return &ProcessManager{}
}

// StartProcess 启动一个新的进程
func (pm *ProcessManager) StartProcess(command string, args ...string) error {
    // 创建命令
    cmd := exec.Command(command, args...)
    
    // 启动进程
    if err := cmd.Start(); err != nil {
        return fmt.Errorf("启动进程失败: %w", err)
    }
    
    // 保存进程引用
    pm.process = cmd.Process
    
    return nil
}

// StopProcess 停止进程
func (pm *ProcessManager) StopProcess() error {
    if pm.process == nil {
        return fmt.Errorf("没有要停止的进程")
    }
    
    // 发送信号停止进程
    if err := pm.process.Signal(syscall.SIGTERM); err != nil {
        return fmt.Errorf("停止进程失败: %w", err)
    }
    
    // 等待进程退出
    _, err := pm.process.Wait()
    if err != nil {
        return fmt.Errorf("等待进程退出失败: %w", err)
    }
    
    // 清理
    pm.process = nil
    
    return nil
}

func main() {
    pm := NewProcessManager()
    
    // 启动一个示例进程
    if err := pm.StartProcess("echo", "Hello, World!"); err != nil {
        log.Fatal(err)
    }

    // 停止进程
    if err := pm.StopProcess(); err != nil {
        log.Fatal(err)
    }
}
