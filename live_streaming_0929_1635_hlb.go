// 代码生成时间: 2025-09-29 16:35:44
package main

import (
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/gorilla/websocket"
)

// LiveStreamingService 结构体封装了直播服务相关的配置
type LiveStreamingService struct {
    upgrader websocket.Upgrader
}

// NewLiveStreamingService 创建一个新的直播服务实例
func NewLiveStreamingService() *LiveStreamingService {
    return &LiveStreamingService{
        upgrader: websocket.Upgrader{
            // 允许所有Origin连接
            CheckOrigin: func(r *http.Request) bool {
                return true
            },
        },
    }
}

// HandleStream 处理直播流的WebSocket连接
func (service *LiveStreamingService) HandleStream(w http.ResponseWriter, r *http.Request) {
    conn, err := service.upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("Failed to set websocket upgrade: ", err)
        return
    }
    defer conn.Close()

    // 循环读取来自客户端的消息
    for {
        mt, message, err := conn.ReadMessage()
        if err != nil {
            fmt.Println("error: ", err)
            break
        }
        fmt.Printf("recv: %s", message)

        // 将消息广播给所有连接的客户端
        // 这里仅作为示例，实际应用中可能需要实现更复杂的逻辑
        if err := conn.WriteMessage(mt, message); err != nil {
            fmt.Println("error: ", err)
            return
        }
    }
}

// main 函数设置路由并启动服务器
func main() {
    service := NewLiveStreamingService()
    r := mux.NewRouter()
    r.HandleFunc("/live", service.HandleStream)

    fmt.Println("Starting live streaming service...

    http.ListenAndServe(":8080", r)
}
