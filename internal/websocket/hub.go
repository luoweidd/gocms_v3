package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// MessageType WebSocket消息类型
type MessageType string

const (
	MessageTypeNotification MessageType = "notification"
	MessageTypeChat         MessageType = "chat"
	MessageTypePing         MessageType = "ping"
	MessageTypePong         MessageType = "pong"
	MessageTypeConnect      MessageType = "connect"
)

// WSMessage WebSocket消息结构
type WSMessage struct {
	Type      MessageType            `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Payload   json.RawMessage        `json:"payload,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// WSHub WebSocket连接中心管理器
type WSHub struct {
	connections sync.Map // clientID -> *WSClient
	broadcastCh chan WSMessage
	mu          sync.RWMutex
	running     bool
	stopCh      chan struct{}
}

// WSClient WebSocket客户端封装
type WSClient struct {
	ID       string          `json:"id"`
	UserID   uint            `json:"user_id"`
	Conn     *websocket.Conn `json:"-"`
	SendMu   sync.Mutex      `json:"-"`
	LastSeen time.Time       `json:"last_seen"`
	Roles    []string        `json:"roles"`
	TenantID uint            `json:"tenant_id,omitempty"`
}

// WSHub全局实例
var hubInstance *WSHub

// GetWSHub 获取WebSocket Hub单例
func GetWSHub() *WSHub {
	if hubInstance == nil {
		hubInstance = newWSHub(1000)
		go hubInstance.run()
	}
	return hubInstance
}

// newWSHub 创建新的WSHub实例
func newWSHub(broadcastSize int) *WSHub {
	return &WSHub{
		connections: sync.Map{},
		broadcastCh: make(chan WSMessage, broadcastSize),
		stopCh:      make(chan struct{}),
	}
}

// Start 启动WebSocket Hub
func (h *WSHub) Start() error {
	h.mu.Lock()
	if h.running {
		return nil
	}
	h.running = true
	h.mu.Unlock()

	log.Println("WebSocket Hub started")
	return nil
}

// Stop 停止WebSocket Hub
func (h *WSHub) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()

	close(h.stopCh)
	h.running = false

	// 关闭所有连接
	h.connections.Range(func(key, value interface{}) bool {
		client := value.(*WSClient)
		client.Conn.Close()
		h.connections.Delete(key)
		return true
	})

	log.Println("WebSocket Hub stopped")
}

// AddClient 添加新的WebSocket客户端
func (h *WSHub) AddClient(client *WSClient) string {
	client.LastSeen = time.Now()
	h.connections.Store(client.ID, client)

	// 发送连接确认消息
	connectMsg := WSMessage{
		Type:      MessageTypeConnect,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message":       "connected",
			"client_id":     client.ID,
			"total_clients": h.ClientCount(),
		},
	}

	h.broadcastTo(client.UserID, connectMsg)

	log.Printf("WebSocket client %s connected. Total: %d", client.ID, h.ClientCount())
	return client.ID
}

// RemoveClient 移除WebSocket客户端
func (h *WSHub) RemoveClient(clientID string) {
	h.connections.Delete(clientID)
	log.Printf("WebSocket client %s disconnected", clientID)
}

// SendToUser 向指定用户发送消息
func (h *WSHub) SendToUser(userID uint, msg WSMessage) error {
	msg.Type = MessageTypeNotification
	msg.Timestamp = time.Now()

	var sent bool
	h.connections.Range(func(key, value interface{}) bool {
		client := value.(*WSClient)
		if client.UserID == userID {
			if err := h.sendToClient(client, msg); err != nil {
				log.Printf("Failed to send message to user %d: %v", userID, err)
			} else {
				sent = true
			}
		}
		return true
	})

	if sent {
		log.Printf("Message sent to user %d", userID)
	}
	return nil
}

// Broadcast 向所有用户广播消息
func (h *WSHub) Broadcast(msg WSMessage) error {
	msg.Type = MessageTypeNotification
	msg.Timestamp = time.Now()

	h.broadcastCh <- msg
	return nil
}

// SendToTenant 向指定租户的所有用户发送消息
func (h *WSHub) SendToTenant(tenantID uint, msg WSMessage) error {
	msg.Timestamp = time.Now()
	count := 0

	h.connections.Range(func(key, value interface{}) bool {
		client := value.(*WSClient)
		if client.TenantID == tenantID {
			if err := h.sendToClient(client, msg); err != nil {
				log.Printf("Failed to send message to tenant %d: %v", tenantID, err)
			} else {
				count++
			}
		}
		return true
	})

	log.Printf("Message sent to %d clients in tenant %d", count, tenantID)
	return nil
}

// sendToClient 向单个客户端发送消息
func (h *WSHub) sendToClient(client *WSClient, msg WSMessage) error {
	client.SendMu.Lock()
	defer client.SendMu.Unlock()

	msg.Timestamp = time.Now()
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = client.Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}

	client.LastSeen = time.Now()
	return nil
}

// broadcastTo 向指定用户组广播消息（内部方法）
func (h *WSHub) broadcastTo(userID uint, msg WSMessage) {
	msg.Type = MessageTypeNotification
	msg.Timestamp = time.Now()

	h.connections.Range(func(key, value interface{}) bool {
		client := value.(*WSClient)
		if userID == 0 || client.UserID == userID {
			go func(c *WSClient, m WSMessage) {
				if err := h.sendToClient(c, m); err != nil {
					log.Printf("Failed to broadcast: %v", err)
				}
			}(client, msg)
		}
		return true
	})
}

// ClientCount 获取当前连接数
func (h *WSHub) ClientCount() int {
	count := 0
	h.connections.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// GetConnectedUsers 获取已连接的用户列表
func (h *WSHub) GetConnectedUsers() []uint {
	var userIDs []uint
	seen := make(map[uint]bool)

	h.connections.Range(func(key, value interface{}) bool {
		client := value.(*WSClient)
		if !seen[client.UserID] {
			userIDs = append(userIDs, client.UserID)
			seen[client.UserID] = true
		}
		return true
	})

	return userIDs
}

// run 运行消息广播循环
func (h *WSHub) run() {
	for {
		select {
		case msg := <-h.broadcastCh:
			h.Broadcast(msg)
		case <-h.stopCh:
			return
		}
	}
}

// CreateMessage 创建标准WebSocket消息
func CreateMessage(msgType MessageType, data map[string]interface{}) WSMessage {
	return WSMessage{
		Type:      msgType,
		Timestamp: time.Now(),
		Data:      data,
	}
}

// CreateNotificationMessage 创建通知消息
func CreateNotificationMessage(userID uint, title string, message string, level string) WSMessage {
	return CreateMessage(MessageTypeNotification, map[string]interface{}{
		"user_id": userID,
		"title":   title,
		"message": message,
		"level":   level, // info, warning, error
		"read":    false,
	})
}

// CreateChatMessage 创建聊天消息
func CreateChatMessage(fromID uint, toID uint, content string) WSMessage {
	return CreateMessage(MessageTypeChat, map[string]interface{}{
		"from":      fromID,
		"to":        toID,
		"content":   content,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// UpgradeToWebSocket Gin WebSocket升级中间件
func UpgradeToWebSocket() gin.HandlerFunc {
	var upgrder = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return func(c *gin.Context) {
		conn, err := upgrder.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v", err)
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("websocket upgrade failed"))
			return
		}

		clientID := c.GetString("client_id")
		userID := c.GetUint("user_id")
		tenantID := c.GetUint("tenant_id")
		roles, _ := c.Get("roles")

		var rolesSlice []string
		if roles != nil {
			rolesSlice = roles.([]string)
		}

		client := &WSClient{
			ID:       clientID,
			UserID:   userID,
			Conn:     conn,
			LastSeen: time.Now(),
			Roles:    rolesSlice,
			TenantID: tenantID,
		}

		hub := GetWSHub()
		hub.AddClient(client)

		// 启动读取循环
		go func(c *WSClient) {
			defer func() {
				c.Conn.Close()
				hub.RemoveClient(c.ID)
			}()

			for {
				_, message, err := c.Conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Printf("WebSocket read error: %v", err)
					}
					break
				}

				// 处理接收到的消息
				var msg WSMessage
				if err := json.Unmarshal(message, &msg); err != nil {
					log.Printf("Failed to unmarshal message: %v", err)
					continue
				}

				switch msg.Type {
				case MessageTypePing:
					pongMsg := CreateMessage(MessageTypePong, map[string]interface{}{
						"timestamp": time.Now().Unix(),
					})
					c.SendMu.Lock()
					data, _ := json.Marshal(pongMsg)
					c.Conn.WriteMessage(websocket.TextMessage, data)
					c.SendMu.Unlock()

				case MessageTypeChat:
					// 转发聊天消息给目标用户
					if toID, ok := msg.Data["to"].(float64); ok {
						hub.SendToUser(uint(toID), msg)
					}
				}
			}
		}(client)

		c.Next()
	}
}
