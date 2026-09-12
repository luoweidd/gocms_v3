package webhook

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// EventType Webhook事件类型
type EventType string

const (
	EventArticleCreated   EventType = "article.created"
	EventArticleUpdated   EventType = "article.updated"
	EventArticleDeleted   EventType = "article.deleted"
	EventArticlePublished EventType = "article.published"
	EventUserRegistered   EventType = "user.registered"
	EventUserUpdated      EventType = "user.updated"
	EventVideoCreated     EventType = "video.created"
	EventVideoUpdated     EventType = "video.updated"
	EventCommentAdded     EventType = "comment.added"
	EventCommentApproved  EventType = "comment.approved"
)

// WebhookEvent Webhook事件结构
type WebhookEvent struct {
	ID         string                 `json:"id"`
	Type       EventType              `json:"type"`
	Timestamp  time.Time              `json:"timestamp"`
	Data       map[string]interface{} `json:"data"`
	RetryCount int                    `json:"-"`
	Status     WebhookStatus          `json:"-"`
}

// WebhookStatus Webhook状态
type WebhookStatus string

const (
	StatusPending  WebhookStatus = "pending"
	StatusSuccess  WebhookStatus = "success"
	StatusFailed   WebhookStatus = "failed"
	StatusRetrying WebhookStatus = "retrying"
)

// WebhookSubscription Webhook订阅配置
type WebhookSubscription struct {
	ID            string      `json:"id"`
	URL           string      `json:"url"`
	Events        []EventType `json:"events"`
	Secret        string      `json:"secret"` // 签名密钥
	Active        bool        `json:"active"`
	CreatedAt     time.Time   `json:"created_at"`
	LastTriggered time.Time   `json:"last_triggered"`
}

// WebhookManager Webhook管理器
type WebhookManager struct {
	subscriptions sync.Map // subscriptionID -> *WebhookSubscription
	queue         chan *WebhookEvent
	handlers      map[EventType][]func(WebhookEvent)
	mu            sync.RWMutex
	running       bool
	stopCh        chan struct{}
}

// NewWebhookManager 创建Webhook管理器
func NewWebhookManager() *WebhookManager {
	return &WebhookManager{
		queue:    make(chan *WebhookEvent, 1000),
		handlers: make(map[EventType][]func(WebhookEvent)),
		stopCh:   make(chan struct{}),
	}
}

// GlobalWebhookManager 全局Webhook管理器实例
var GlobalWebhookManager *WebhookManager

// GetWebhookManager 获取全局Webhook管理器单例
func GetWebhookManager() *WebhookManager {
	if GlobalWebhookManager == nil {
		GlobalWebhookManager = NewWebhookManager()
		go GlobalWebhookManager.run()
	}
	return GlobalWebhookManager
}

// Subscribe 订阅事件
func (wm *WebhookManager) Subscribe(subscription *WebhookSubscription) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if subscription.Events == nil {
		subscription.Events = []EventType{}
	}
	if !subscription.Active {
		subscription.Active = true
	}
	subscription.CreatedAt = time.Now()

	wm.subscriptions.Store(subscription.ID, subscription)
	log.Printf("Webhook订阅已创建: %s 事件=%v", subscription.ID, subscription.Events)
	return nil
}

// Unsubscribe 取消订阅
func (wm *WebhookManager) Unsubscribe(id string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wm.subscriptions.Delete(id)
	log.Printf("Webhook订阅已取消: %s", id)
}

// GetSubscriptions 获取所有有效订阅
func (wm *WebhookManager) GetSubscriptions(eventType EventType) []*WebhookSubscription {
	var subs []*WebhookSubscription
	wm.subscriptions.Range(func(key, value interface{}) bool {
		sub := value.(*WebhookSubscription)
		if sub.Active {
			for _, e := range sub.Events {
				if e == eventType {
					subs = append(subs, sub)
					break
				}
			}
		}
		return true
	})
	return subs
}

// Emit 触发事件
func (wm *WebhookManager) Emit(eventType EventType, data map[string]interface{}) error {
	event := &WebhookEvent{
		ID:         fmt.Sprintf("wh_%d", time.Now().UnixNano()),
		Type:       eventType,
		Timestamp:  time.Now(),
		Data:       data,
		RetryCount: 0,
		Status:     StatusPending,
	}

	select {
	case wm.queue <- event:
		log.Printf("Webhook事件已触发: %s ID=%s", eventType, event.ID)
	default:
		log.Printf("Webhook队列已满，丢弃事件: %s", eventType)
	}

	return nil
}

// EmitWithSubscribers 触发事件并通知指定订阅者
func (wm *WebhookManager) EmitWithSubscribers(eventType EventType, data map[string]interface{}, subscriberIDs []string) error {
	event := &WebhookEvent{
		ID:         fmt.Sprintf("wh_%d", time.Now().UnixNano()),
		Type:       eventType,
		Timestamp:  time.Now(),
		Data:       data,
		RetryCount: 0,
		Status:     StatusPending,
	}

	// 查找指定订阅者
	for _, id := range subscriberIDs {
		val, ok := wm.subscriptions.Load(id)
		if !ok {
			continue
		}
		sub := val.(*WebhookSubscription)

		// 生成签名
		signature := wm.generateSignature(event, sub.Secret)

		// 发送Webhook
		go wm.sendWebhook(sub, event, signature)
	}

	return nil
}

// sendWebhook 发送Webhook通知
func (wm *WebhookManager) sendWebhook(sub *WebhookSubscription, event *WebhookEvent, signature string) {
	payload := map[string]interface{}{
		"id":        event.ID,
		"type":      string(event.Type),
		"timestamp": event.Timestamp.Format(time.RFC3339),
		"data":      event.Data,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Webhook序列化失败: %v", err)
		event.Status = StatusFailed
		return
	}

	req, err := http.NewRequest("POST", sub.URL, bytes.NewReader(payloadBytes))
	if err != nil {
		log.Printf("Webhook创建请求失败: %v", err)
		event.Status = StatusFailed
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Signature", signature)
	req.Header.Set("X-Webhook-Event", string(event.Type))
	req.Header.Set("X-Webhook-ID", event.ID)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Webhook发送失败: %v", err)
		event.Status = StatusFailed
		event.RetryCount++
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		event.Status = StatusSuccess
		sub.LastTriggered = time.Now()
		log.Printf("Webhook发送成功: %s URL=%s", event.Type, sub.URL)
	} else {
		event.Status = StatusFailed
		event.RetryCount++
		log.Printf("Webhook响应异常: %d URL=%s", resp.StatusCode, sub.URL)
	}
}

// generateSignature 生成Webhook签名
func (wm *WebhookManager) generateSignature(event *WebhookEvent, secret string) string {
	if secret == "" {
		return ""
	}

	signString := fmt.Sprintf("%s.%d.%s", event.ID, event.Timestamp.Unix(), string(event.Type))
	data := []byte(signString + secret)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("sha256_%x", hash[:16])
}

// run 运行Webhook处理循环
func (wm *WebhookManager) run() {
	for {
		select {
		case event := <-wm.queue:
			wm.processEvent(event)
		case <-wm.stopCh:
			return
		}
	}
}

// processEvent 处理单个事件
func (wm *WebhookManager) processEvent(event *WebhookEvent) {
	subs := wm.GetSubscriptions(event.Type)

	for _, sub := range subs {
		signature := wm.generateSignature(event, sub.Secret)
		wm.sendWebhook(sub, event, signature)

		// 触发事件处理器
		if handlers, ok := wm.handlers[event.Type]; ok {
			for _, handler := range handlers {
				go handler(*event)
			}
		}
	}
}

// RegisterHandler 注册事件处理器
func (wm *WebhookManager) RegisterHandler(eventType EventType, handler func(WebhookEvent)) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wm.handlers[eventType] = append(wm.handlers[eventType], handler)
}

// Stop 停止Webhook管理器
func (wm *WebhookManager) Stop() {
	close(wm.stopCh)
	log.Println("Webhook管理器已停止")
}

// CreateSubscription 创建新订阅
func CreateSubscription(id, url string, events []EventType, secret string) *WebhookSubscription {
	return &WebhookSubscription{
		ID:        id,
		URL:       url,
		Events:    events,
		Secret:    secret,
		Active:    true,
		CreatedAt: time.Now(),
	}
}

// DefaultEventPayload 默认事件数据工厂
type DefaultEventPayload struct {
	Event     EventType              `json:"event"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// NewEventPayload 创建默认事件载荷
func NewEventPayload(eventType EventType, data map[string]interface{}) *DefaultEventPayload {
	return &DefaultEventPayload{
		Event:     eventType,
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      data,
	}
}
