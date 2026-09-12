package model

// OperationLog 操作日志
type OperationLog struct {
	BaseModel
	UserID        uint64 `gorm:"index;default:0" json:"user_id"`
	Username      string `gorm:"size:100;default:''" json:"username"`
	Action        string `gorm:"size:100;index;not null" json:"action"`    // 操作类型: create/update/delete/publish/unpublish/approve/reject/batch_delete/import/export等
	Module        string `gorm:"size:100;default:''" json:"module"`        // 模块名称: article/video/comment/menu/system等
	Description   string `gorm:"size:500;default:''" json:"description"`   // 操作描述
	IPAddress     string `gorm:"size:45;default:''" json:"ip_address"`     // IP地址（支持IPv6）
	UserAgent     string `gorm:"size:500;default:''" json:"user_agent"`    // 浏览器UA
	RequestMethod string `gorm:"size:10;default:''" json:"request_method"` // GET/POST/PUT/DELETE等
	RequestPath   string `gorm:"size:255;default:''" json:"request_path"`  // 请求路径
	RequestData   string `gorm:"type:text" json:"request_data"`            // 请求数据(JSON)
	ResponseBody  string `gorm:"type:text" json:"response_body"`           // 响应数据(JSON)
	DurationMs    int    `gorm:"default:0" json:"duration_ms"`             // 耗时(毫秒)
	Status        int8   `gorm:"default:1" json:"status"`                  // 1成功 0失败
	ErrorMessage  string `gorm:"size:500;default:''" json:"error_message"` // 错误信息（失败时）
}

func (OperationLog) TableName() string {
	return "operation_logs"
}

// OperationLogListRequest 操作日志列表查询请求
type OperationLogListRequest struct {
	Page      int    `form:"page" binding:"min=1"`
	PageSize  int    `form:"page_size" binding:"max=200"`
	Module    string `form:"module"`
	Action    string `form:"action"`
	Username  string `form:"username"`
	IPAddress string `form:"ip_address"`
	Status    *int8  `form:"status"`
	StartDate string `form:"start_date"` // 筛选起始日期
	EndDate   string `form:"end_date"`   // 筛选结束日期
}

// OperationLogStats 操作日志统计
type OperationLogStats struct {
	TotalCount   int64        `json:"total_count"`   // 总记录数
	TodayCount   int64        `json:"today_count"`   // 今日记录数
	SuccessCount int64        `json:"success_count"` // 成功次数
	FailCount    int64        `json:"fail_count"`    // 失败次数
	ModuleStats  []ModuleStat `json:"module_stats"`  // 按模块统计
	ActionStats  []ActionStat `json:"action_stats"`  // 按操作类型统计
}

// ModuleStat 模块统计
type ModuleStat struct {
	Module string `json:"module"`
	Count  int64  `json:"count"`
}

// ActionStat 操作统计
type ActionStat struct {
	Action string `json:"action"`
	Count  int64  `json:"count"`
}
